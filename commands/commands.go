package commands

import (
	"github.com/bwmarrin/discordgo"
)

var SlashCommands = []*discordgo.ApplicationCommand{
	&statusSlashCommand,
	&addButtonsSlashCommand,
	&repeatSlashCommand,
	//&recoverSlashCommand,
}

var statusSlashCommand = discordgo.ApplicationCommand{
	Name:        "status",
	Type:        1,
	Description: "Check status of PFC Bot",
}

var recoverSlashCommand = discordgo.ApplicationCommand{
	Name:        "recoveraccount",
	Type:        1,
	Description: "this command helps you recover your PCL data in the event of a new discord ID",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        3,
			Name:        "discordname",
			Description: "What's the discord name you previously registered with?",
			Required:    true,
			MaxValue:    1,
			MaxLength:   32,
		},
	},
}

var repeatSlashCommand = discordgo.ApplicationCommand{
	Name:        "repeat",
	Type:        1,
	Description: "You know what to do",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionChannel,
			Name:        "channel",
			Description: "Where to",
			Required:    true,
		},
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "message",
			Description: "What to say",
			Required:    true,
			MaxLength:   2000,
		},
	},
}

var addButtonsSlashCommand = discordgo.ApplicationCommand{
	Name:        "addbutton",
	Type:        1,
	Description: "Choose what button to initialize",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        3,
			Name:        "name",
			Description: "Choose button to initialize",
			Required:    true,
			Choices: []*discordgo.ApplicationCommandOptionChoice{
				{
					Name:  "register",
					Value: "register",
				},
			},
			MaxValue:  1,
			MaxLength: 32,
		},
	},
}
