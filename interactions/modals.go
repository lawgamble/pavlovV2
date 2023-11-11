package interactions

import (
	"github.com/bwmarrin/discordgo"
	"log"
	mariadb "pfc2/mariaDB"
	"strings"
)

// OpenRegistrationModal sends a pop-up modal window, allowing a user to input text fields.
func OpenRegistrationModal(s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID: "createRegistration",
			Title:    "Register",
			Components: []discordgo.MessageComponent{
				RegionActionRow,
				PlayStyleActionRow,
				PlayerTypeActionRow,
				DOBActionRow,
				InGameNameActionRow,
			},
		},
	})
	if err != nil {
		log.Print(err)
	}
}

// OpenFilledRegistrationModal sends a pop-up modal window, allowing a user to update text fields. This modal is only viewable if the current user is registered.
func OpenFilledRegistrationModal(s *discordgo.Session, i *discordgo.InteractionCreate, p mariadb.Player) {
	region := RegionActionRow
	region.Components[0].(*discordgo.TextInput).Value = p.Region
	playStyle := PlayStyleActionRow
	playStyle.Components[0].(*discordgo.TextInput).Value = p.PlayStyle
	playerType := PlayerTypeActionRow
	playerType.Components[0].(*discordgo.TextInput).Value = p.PlayerType
	inGameName := InGameNameActionRow
	inGameName.Components[0].(*discordgo.TextInput).Value = p.InGameName

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID: "updateRegistration",
			Title:    "Update Registration",
			Components: []discordgo.MessageComponent{
				region,
				playStyle,
				playerType,
				formatDOBForModal(p),
				inGameName,
			},
		},
	})
	if err != nil {
		log.Print(err)
	}
}

func formatDOBForModal(p mariadb.Player) discordgo.MessageComponent {
	dob := DOBActionRow
	d := p.DOB
	stringDob := string(d)
	dateComponents := strings.Split(stringDob, "-")

	year := dateComponents[0]
	month := dateComponents[1]
	day := dateComponents[2]

	// Concatenate and format the components as "MMDDYYYY"
	dob.Components[0].(*discordgo.TextInput).Value = month + day + year
	return dob
}

// RegistrationErrorResponse responds with the validation error upon registration. It is only visible to the user.
func RegistrationErrorResponse(s *discordgo.Session, i *discordgo.InteractionCreate, err error) {
	resErr := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Registration Failed: " + err.Error(),
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	if resErr != nil {
		log.Print(resErr)
		return
	}
}

// RegistrationSuccessResponse returns a success response only visible to the user.
func RegistrationSuccessResponse(s *discordgo.Session, i *discordgo.InteractionCreate) {
	resErr := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Success! You've been registered to PCL!",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	if resErr != nil {
		log.Print(resErr)
		return
	}
}

func UpdatedRegistrationSuccessResponse(s *discordgo.Session, i *discordgo.InteractionCreate) {
	resErr := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Success! Your info has been updated!",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	if resErr != nil {
		log.Print(resErr)
		return
	}
}
