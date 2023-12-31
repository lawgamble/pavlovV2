package interactions

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	mariadb "pfc2/mariaDB"
	"strings"
)

// OpenRegistrationModal sends a pop-up modal window, allowing a user to input text fields.
func OpenRegistrationModal(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var blankInteractionResponse = discordgo.InteractionResponse{
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
	}
	err := s.InteractionRespond(i.Interaction, &blankInteractionResponse)
	if err != nil {
		log.Print(err)
	}
}

// OpenFilledRegistrationModal sends a pop-up modal window, allowing a user to update text fields. This modal is only viewable if the current user is registered.
func OpenFilledRegistrationModal(s *discordgo.Session, i *discordgo.InteractionCreate, p mariadb.Player) {
	region := FilledRegionActionRow
	region.Components[0].(*discordgo.TextInput).Value = p.Region
	playStyle := FilledPlayStyleActionRow
	playStyle.Components[0].(*discordgo.TextInput).Value = p.PlayStyle
	playerType := FilledPlayerTypeActionRow
	playerType.Components[0].(*discordgo.TextInput).Value = p.PlayerType
	inGameName := FilledInGameNameActionRow
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
	dob := FilledDOBActionRow
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
			Content: "***Registration Failed***:\n" + err.Error(),
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

func TeamRegistrationSuccessResponse(s *discordgo.Session, i *discordgo.InteractionCreate, teamName string) {
	resErr := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("%v has been registered! Mods will need to approve your team.", teamName),
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	if resErr != nil {
		log.Print(resErr)
		return
	}
}

func DefaultErrorResponse(s *discordgo.Session, i *discordgo.InteractionCreate, err error) {
	resErr := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("%v", err),
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	if resErr != nil {
		log.Print(resErr)
		return
	}
}

func SendTempTeamsTable(s *discordgo.Session, i *discordgo.InteractionCreate, table string) {
	resErr := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: table,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	if resErr != nil {
		log.Print(resErr)
		return
	}
}
