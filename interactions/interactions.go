package interactions

import (
	"database/sql"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	mariadb "pfc2/mariaDB"
	"regexp"
	"strconv"
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

	err := checkIfActiveTeamAlreadyExists(teamName, db)
	if err != nil {
		// team already exists with registered name (case-insensitive)
		RegistrationErrorResponse(s, i, err)
		return
	}

	for j := 1; j < len(interaction.Options); j++ {
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
		err := sendTeamRegistration(players, teamName, db)
		if err != nil {
			// sendTeamRegistration error
			log.Print(err)
			RegistrationErrorResponse(s, i, err)
			return
		}
		// all players were registered and no error from db submission
		TeamRegistrationSuccessResponse(s, i, teamName)
	} else {
		unregisteredErrorMsg := whoIsNotRegistered(unregisteredPlayers)
		RegistrationErrorResponse(s, i, fmt.Errorf("%v", unregisteredErrorMsg))
	}
}

func checkIfActiveTeamAlreadyExists(teamName string, db mariadb.DBHandler) error {
	returnedTeam, _ := db.DB.ReadTeamByTeamName(teamName) // a team can be returned here if needed. For now, err will determine if row was not found
	if returnedTeam.TeamName != "" {
		err := fmt.Errorf(" `%v`\n has already been chosen - try another team name", teamName)
		return err
	}
	return nil
}

func whoIsNotRegistered(players []string) string {
	msg := "Users are not registered: "
	for _, player := range players {
		msg += fmt.Sprintf("<@%s> ", player)
	}
	return strings.TrimSpace(msg)
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

// sendTeamRegistration adds players to a temp table in the DB, along with adding pending team to team table.
func sendTeamRegistration(players []mariadb.Player, teamName string, db mariadb.DBHandler) error {
	for _, player := range players {
		// avoiding concurrency here (not sure how DB will handle concurrent writes)
		err := db.DB.CreateTempRoster(strconv.FormatInt(player.DiscordId, 10), teamName)
		if err != nil {
			return err
		}
	}
	err := db.DB.CreateTempTeam(teamName)
	if err != nil {
		return err
	}
	return nil
}

// hasDuplicates takes in a slice of string and checks to see if there are exact string matches within any index.
func hasDuplicates(slice []string) (bool, string) {
	seen := make(map[string]bool)

	for i, value := range slice {
		// If the value is already in the map, it's a duplicate
		if seen[value] {
			return true, slice[i]
		}

		// Mark the value as seen
		seen[value] = true
	}

	// No duplicates found
	return false, ""
}

func containsOnlyNumbers(input string) bool {
	// Use a regular expression to check if the string contains only numbers
	match, _ := regexp.MatchString("^[0-9]+$", input)
	return match
}
