package interactions

import (
	"database/sql"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
	mariadb "pfc2/mariaDB"
	"strings"
	"sync"
	"time"
)

// HandleMessageComponent handles button presses - will switch between customIds.
func HandleMessageComponent(s *discordgo.Session, i *discordgo.InteractionCreate, db mariadb.DBHandler) {
	customId := i.MessageComponentData().CustomID
	switch customId {
	case "register":
		{
			Register(s, i, db)
			// make i nil here?
		}
	}
}

func SubmitRegistration(i *discordgo.InteractionCreate, db mariadb.DBHandler, m mariadb.ModalSubmitData) error {
	m.DiscordId = i.Member.User.ID
	m.DiscordName = i.Member.User.Username
	registrationDate := time.Now()
	// call db with unique query, if no error, send Success Message
	err := db.DB.CreatePlayer(m, registrationDate)
	return err
}

func SubmitUpdatedRegistration(db mariadb.DBHandler, m mariadb.ModalSubmitData) error {
	return db.DB.Update(m)
}

func Register(s *discordgo.Session, i *discordgo.InteractionCreate, db mariadb.DBHandler) {
	//first, check DB if user is registered
	discordId := i.Member.User.ID
	player, err := db.DB.ReadUsersByDiscordId(discordId)

	if err == sql.ErrNoRows {
		// player is not registered and needs to
		OpenRegistrationModal(s, i)
	} else {
		//player is registered and wants to update their information.
		updatedPlayer := SetFieldValuesToIntegers(player)
		OpenFilledRegistrationModal(s, i, updatedPlayer)
	}
}

func SetFieldValuesToIntegers(p mariadb.Player) mariadb.Player {
	if p.PlayerType != "" {
		switch p.PlayerType {
		case "Draftable":
			p.PlayerType = "1"
			break
		case "Team Member":
			p.PlayerType = "2"
			break
		case "Pickups Only":
			p.PlayerType = "3"
			break
		default:
			p.PlayerType = "X"
			break
		}
	}
	if p.PlayStyle != "" {
		switch p.PlayStyle {
		case "Flex":
			p.PlayStyle = "1"
			break
		case "Rush/Entry":
			p.PlayStyle = "2"
			break
		case "Lurker":
			p.PlayStyle = "3"
			break
		case "Mid Player":
			p.PlayStyle = "4"
			break
		case "IGL":
			p.PlayStyle = "5"
			break
		default:
			p.PlayStyle = "X"
			break
		}
	}
	if p.Region != "" {
		switch p.Region {
		case "USA":
			p.Region = "1"
			break
		case "USA ONLY":
			p.Region = "2"
			break
		case "EU":
			p.Region = "3"
			break
		case "EU ONLY":
			p.Region = "4"
			break
		default:
			p.Region = "X"
			break
		}
	}
	return p
}

func RegisterTeam(s *discordgo.Session, i *discordgo.InteractionCreate, interaction discordgo.ApplicationCommandInteractionData, db mariadb.DBHandler) {
	var playerIds []string

	player1 := i.Member.User.ID
	playerIds = append(playerIds, player1)

	teamName := interaction.Options[0].Value.(string)
	teamRegion := interaction.Options[1].Value.(string)

	err := checkIfActiveTeamAlreadyExists(teamName, db)
	if err != nil {
		// team already exists with registered name (case-insensitive)
		RegistrationErrorResponse(s, i, err)
		return
	}

	for j := 2; j < len(interaction.Options); j++ {
		playerId := interaction.Options[j].Value.(string)
		if playerId[1] == '@' {
			playerId = strings.ReplaceAll(playerId, "<@", "")
			playerId = strings.ReplaceAll(playerId, ">", "")
		}
		if containsOnlyNumbers(playerId) {
			playerIds = append(playerIds, playerId)
		} else {
			RegistrationErrorResponse(s, i, fmt.Errorf("you can only add valid discord users to your team. Ex: @PLAYER"))
			return
		}
	}

	// check to see if there's duplicate strings
	hasDuplicates, dupedUser := hasDuplicates(playerIds)
	if hasDuplicates {
		if dupedUser == player1 {
			dupErr := fmt.Errorf("no need to add yourself,<@%v>", dupedUser)
			RegistrationErrorResponse(s, i, dupErr)
			return
		}
		dupErr := fmt.Errorf("<@%v> can't be on the team twice", dupedUser)
		RegistrationErrorResponse(s, i, dupErr)
		return
	}
	isAllRegistered, players, unregisteredPlayers := ValidateAllTeamMembers(playerIds, db)
	if isAllRegistered {
		err := sendTeamRegistration(players, teamName, teamRegion, db, player1)
		if err != nil {
			// sendTeamRegistration error
			log.Print(err)
			RegistrationErrorResponse(s, i, err)
			return
		}
		// all players were registered and no error from db submission
		TeamRegistrationSuccessResponse(s, i, teamName)
		// TODO send message to a channel and tag a league manager (TeamRequests Channel)
		msg := "<@&878021874431434803> - " + teamName + " just registered!"
		_, _ = s.ChannelMessageSend(os.Getenv("TEAM_REQUESTS_CHAN_ID"), msg)
	} else {
		unregisteredErrorMsg := whoIsNotRegistered(unregisteredPlayers)
		RegistrationErrorResponse(s, i, fmt.Errorf("%v", unregisteredErrorMsg))
	}
}

func ValidateAllTeamMembers(ids []string, db mariadb.DBHandler) (bool, []mariadb.Player, []string) {
	var players []mariadb.Player
	var unregisteredPlayers []string
	var wg sync.WaitGroup
	var mu sync.Mutex // protect the concurrent write to players slice
	var errChan = make(chan error, len(ids))
	done := make(chan struct{})

	for _, id := range ids {
		wg.Add(1)

		go func(id string) {
			defer wg.Done() // Ensure that Done is called even if the condition below is true

			select {
			case <-done:
				// If done channel is closed, exit goroutine
				return
			default:
				player, err := db.DB.ReadUsersByDiscordId(id)
				if err != nil {
					unregisteredPlayers = append(unregisteredPlayers, id)
					// Sends error to the errChan
					errChan <- err
					return
				}

				if player.IsNotEmpty() {
					mu.Lock()
					players = append(players, player)
					mu.Unlock()
				}
			}
		}(id)
	}

	wg.Wait()
	close(done)

	select {
	case err := <-errChan:
		log.Printf(err.Error())
		return false, players, unregisteredPlayers
	default:
		return true, players, unregisteredPlayers
	}
}
