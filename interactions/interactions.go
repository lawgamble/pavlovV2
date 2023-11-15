package interactions

import (
	"database/sql"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
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
	player, err := db.DB.ReadByDiscordId(discordId)

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
	//collect all necessary data
	var playerIds []string

	teamName := interaction.Options[0].Value.(string)
	player1 := i.Member.User.ID

	playerIds = append(playerIds, player1)

	for i := 1; i < len(interaction.Options); i++ {
		playerId := interaction.Options[i].Value.(string)
		playerId = strings.ReplaceAll(playerId, "<@", "")
		playerId = strings.ReplaceAll(playerId, ">", "")
		playerIds = append(playerIds, playerId)
	}
	isValid, players, err := ValidateAllTeamMembers(playerIds, db, player1)
	if isValid {
		err := sendTeamRegistration(players, teamName, db)
		if err != nil {
			log.Print(err)
			RegistrationErrorResponse(s, i, err)
		}
	} else {
		RegistrationErrorResponse(s, i, fmt.Errorf("%v", err))
		return
	}
	TeamRegistrationSuccessResponse(s, i, teamName)
}

func ValidateAllTeamMembers(ids []string, db mariadb.DBHandler, player1 string) (bool, []mariadb.Player, error) {
	var players []mariadb.Player
	var wg sync.WaitGroup
	var mu sync.Mutex // protect the concurrent write to players slice
	var errChan = make(chan error, len(ids))
	done := make(chan struct{})

	for _, id := range ids {
		wg.Add(1)

		go func(id string) {
			defer wg.Done() // Ensure that Done is called even if the condition below is true

			if id == player1 {
				err := fmt.Errorf("no need to register yourself")
				errChan <- err
				return
			}

			select {
			case <-done:
				// If done channel is closed, exit goroutine
				return
			default:
				player, err := db.DB.ReadByDiscordId(id)
				if err != nil {
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
		log.Println("Error:", err)
		return false, players, err
	default:
		return true, players, nil
	}
}

func sendTeamRegistration(players []mariadb.Player, name string, db mariadb.DBHandler) error {
	return nil
}
