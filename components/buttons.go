package components

import (
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
)

func InitializeButtons(s *discordgo.Session, i *discordgo.InteractionCreate, o []*discordgo.ApplicationCommandInteractionDataOption) {
	if o[0] != nil {
		switch o[0].StringValue() {
		case "register":
			{
				err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: 4,
				})
				if err != nil {
					log.Print(err)
				}
				registerChannel := os.Getenv("REGISTERCHANID")
				message := discordgo.MessageSend{
					Content: "",
					Components: []discordgo.MessageComponent{
						RegistrationButtons,
					},
				}
				s.ChannelMessageSendComplex(registerChannel, &message)
			}
		}
	}
}

func deleteChannelMessages(s *discordgo.Session, chanId string, count int) error {
	msgs, err := s.ChannelMessages(chanId, count, "", "", "")

	for _, message := range msgs {
		err = s.ChannelMessageDelete(chanId, message.ID)
		if err != nil {
			log.Println(err)
		}
	}
	return err
}

func getChannels() map[string]string {
	c := make(map[string]string)
	c["register"] = os.Getenv("REGISTERCHANID")
	return c
}

var RegisterEnableButton = discordgo.Button{
	Label:    "Register",
	Style:    discordgo.SuccessButton,
	CustomID: "register",
}

//RegistrationButtons contains the buttons needed for the registration page.
var RegistrationButtons = discordgo.ActionsRow{
	Components: []discordgo.MessageComponent{
		RegisterEnableButton,
	},
}
