package commands

import (
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"pfc2/components"
	"pfc2/interactions"
	mariadb "pfc2/mariaDB"
	"time"
)

// HandleCommands handles all slash commands and InteractionModalSubmit interactions like button presses inside a modal.
func HandleCommands(s *discordgo.Session, i *discordgo.InteractionCreate, db mariadb.DBHandler) {
	switch i.Type {
	case discordgo.InteractionModalSubmit:

		modalSubmitData := i.ModalSubmitData()
		{
			switch modalSubmitData.CustomID {
			case "createRegistration":
				{
					isValid, err, submitData := validateRegistrationData(modalSubmitData.Components)
					if !isValid {
						interactions.RegistrationErrorResponse(s, i, err)
						break
					} else {
						err := interactions.SubmitRegistration(i, db, submitData)
						if err != nil {
							//there was an error calling the DB
							log.Print(err)
							interactions.RegistrationErrorResponse(s, i, err)
							break
						}
						interactions.RegistrationSuccessResponse(s, i)
						// give user PlayerType Role - Enlisted/Draft - PickupsOnly
						err = writeToAliasFile(i.Member.User.ID, submitData.InGameName)
						if err != nil {
							log.Print(err)
						}
						break
					}
				}
			case "updateRegistration":
				{
					isValid, err, submitData := validateRegistrationData(modalSubmitData.Components)
					if !isValid {
						interactions.RegistrationErrorResponse(s, i, err)
						break
					} else {
						submitData.DiscordId = i.Member.User.ID
						err := interactions.SubmitUpdatedRegistration(db, submitData)
						if err != nil {
							//there was an error calling the DB
							log.Print(err)
							interactions.RegistrationErrorResponse(s, i, err)
							break
						}
						interactions.UpdatedRegistrationSuccessResponse(s, i)

						err = writeToAliasFile(i.Member.User.ID, submitData.InGameName)
						if err != nil {
							log.Print(err)
						}

						break
					}
				}

			}
		}
	case discordgo.InteractionApplicationCommand:
		{
			interaction := i.ApplicationCommandData()
			switch interaction.Name {
			case "status":
				{
					Status(s, i)
					break
				}
			case "addbutton":
				{
					components.InitializeButtons(s, i, interaction.Options)
					break
				}
			case "recoveraccount":
				{
					fmt.Println("You did it!")
				}
			}

		}
	}
}

func validateRegistrationData(c []discordgo.MessageComponent) (bool, error, mariadb.ModalSubmitData) {
	modalSubmitData := mariadb.ModalSubmitData{
		Region:     c[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value,
		PlayStyle:  c[1].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value,
		PlayerType: c[2].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value,
		DOB:        c[3].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value,
		InGameName: c[4].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value,
	}
	// first, validate age, so we don't have to validate more if we don't have to
	err := modalSubmitData.ValidateDateOfBirth()
	if err != nil {
		return false, err, modalSubmitData
	}

	err = modalSubmitData.ValidatePlayerType()
	if err != nil {
		return false, err, modalSubmitData
	}

	err = modalSubmitData.ValidatePlayStyle()
	if err != nil {
		return false, err, modalSubmitData
	}

	err = modalSubmitData.ValidateRegion()
	if err != nil {
		return false, err, modalSubmitData
	}

	err = modalSubmitData.ValidateInGameName()
	if err != nil {
		return false, err, modalSubmitData
	}
	return true, nil, modalSubmitData
}

// Status returns a random string of the bot's "status" when a user runs the /status command.
func Status(s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: statusResponse(),
		},
	},
	)
	if err != nil {
		//do nothing, I guess?
	}
}

//statusResponse returns a response string for the Status() func - at random.
func statusResponse() string {
	rand.Seed(time.Now().UnixNano())
	randomIndex := rand.Intn(len(statusOptions))
	return statusOptions[randomIndex]
}

func writeToAliasFile(discordId string, inGameName string) error {
	filePath := os.Getenv("ALIASES_FILEPATH")
	// Read the JSON data from the file
	inGameName = "q-" + inGameName
	jsonFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	// Unmarshal the JSON data into a map
	var data map[string]interface{}
	if err := json.Unmarshal(jsonFile, &data); err != nil {
		return err
	}

	// Check if the "players" object exists
	players, ok := data["players"].(map[string]interface{})
	if !ok {
		players = make(map[string]interface{})
		data["players"] = players
	}

	// Check if the discordId already exists in "players"
	existingName, exists := players[discordId].(string)
	if exists {
		// If inGameName is different, update it
		if existingName != inGameName {
			players[discordId] = inGameName
		}
	} else {
		// Add the new player
		players[discordId] = inGameName
	}

	// Marshal the updated data back to JSON with indentation
	updatedJSON, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	// Write the updated JSON back to the file
	if err := ioutil.WriteFile(filePath, updatedJSON, 0644); err != nil {
		return err
	}

	return nil
}
