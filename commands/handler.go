package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
	"pfc2/components"
	"pfc2/interactions"
	mariadb "pfc2/mariaDB"
)

func HandleModalSubmit(s *discordgo.Session, i *discordgo.InteractionCreate, db mariadb.DBHandler) {
	modalSubmitData := i.ModalSubmitData()
	registerRoleId := os.Getenv("REGISTERED_ROLE_ID")
	enlistedRoleId := os.Getenv("ENLISTED_ROLE_ID")
	isValid, err, submitData := validateRegistrationData(modalSubmitData.Components)
	if !isValid {
		interactions.RegistrationErrorResponse(s, i, err)
		return
	}

	switch modalSubmitData.CustomID {
	case "createRegistration":
		{
			err := interactions.SubmitRegistration(i, db, submitData)
			if err != nil {
				//there was an error calling the DB
				log.Print(err)
				interactions.RegistrationErrorResponse(s, i, err)
				break
			}
			interactions.RegistrationSuccessResponse(s, i)
			interactions.GiveRoleToUser(s, i, registerRoleId)

			if submitData.PlayerType == "Draftable" {
				interactions.GiveRoleToUser(s, i, enlistedRoleId)
			}
			err = writeToAliasFile(i.Member.User.ID, submitData.InGameName)
			if err != nil {
				log.Print(err)
			}
			break
		}

	case "updateRegistration":
		{
			submitData.DiscordId = i.Member.User.ID
			err := interactions.SubmitUpdatedRegistration(db, submitData)
			if err != nil {
				//there was an error calling the DB
				log.Print(err)
				interactions.RegistrationErrorResponse(s, i, err)
				break
			}
			interactions.UpdatedRegistrationSuccessResponse(s, i)

			if submitData.PlayerType == "Draftable" {
				interactions.GiveRoleToUser(s, i, enlistedRoleId)
			} else {
				interactions.RemoveRoleFromUser(s, i, enlistedRoleId)
			}
			err = writeToAliasFile(i.Member.User.ID, submitData.InGameName)
			if err != nil {
				log.Print(err)
			}
			break
		}
	}
}

// HandleApplicationCommands handles all slash commands
func HandleApplicationCommands(s *discordgo.Session, i *discordgo.InteractionCreate) {
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
	case "repeat":
		{
			HandleRepeatCommand(s, i, interaction)
		}
	case "recoveraccount":
		{
			fmt.Println("You did it!")
		}
	}
}

func HandleRepeatCommand(s *discordgo.Session, i *discordgo.InteractionCreate, interaction discordgo.ApplicationCommandInteractionData) {
	if !UserHasRole(s, i, os.Getenv("MODERATOR_ROLE_ID")) {
		err := s.InteractionRespond(i.Interaction, interactions.NotPermittedInteractionResponse)
		if err != nil {
			log.Print(err)
		}
		return
	}
	// If we get here, user has permission to run command
	chanId := interaction.Options[0].Value.(string)
	message := interaction.Options[1].Value.(string)
	err := s.InteractionRespond(i.Interaction, interactions.RepeatCommandInteractionResponse)
	if err != nil {
		log.Print(err)
	}
	_, err = s.ChannelMessageSend(chanId, message)
	if err != nil {
		log.Print(err)
	}
}
