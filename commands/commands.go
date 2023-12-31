package commands

import (
	"github.com/bwmarrin/discordgo"
)

var SlashCommands = []*discordgo.ApplicationCommand{
	&statusSlashCommand,
	&addButtonsSlashCommand,
	&repeatSlashCommand,
	&teamRegistrationSlashCommand,
	&listApprovalsSlashCommand,
	&approveTeamSlashCommand,
	&denyTeamSlashCommand,
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

var teamRegistrationSlashCommand = discordgo.ApplicationCommand{
	Name:        "teamregister",
	Type:        1,
	Description: "register a team and the players on it; you are player1",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "team_name",
			Description: "your team's name:",
			Required:    true,
			MaxLength:   50,
		},
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "team_region",
			Description: "US or EU?",
			Required:    true,
			MaxLength:   2,
		},
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "player2",
			Description: "2nd player",
			Required:    true,
			MaxLength:   35,
		},
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "player3",
			Description: "3rd player",
			Required:    true,
			MaxLength:   35,
		},
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "player4",
			Description: "4th player",
			Required:    true,
			MaxLength:   35,
		},
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "player5",
			Description: "5th player",
			Required:    true,
			MaxLength:   35,
		},
	},
}

// /listapprovals, show "pending" status teams + Players + In-Game-Name (chart)?

var listApprovalsSlashCommand = discordgo.ApplicationCommand{
	Name:        "listapprovals",
	Type:        1,
	Description: "show pending status teams + Players + In-Game-Name",
}

var approveTeamSlashCommand = discordgo.ApplicationCommand{
	Name:        "approveteam",
	Type:        1,
	Description: "approve a registered team",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "team_name",
			Description: "name of team you're approving",
			Required:    true,
			MaxLength:   35,
		},
	},
}

var denyTeamSlashCommand = discordgo.ApplicationCommand{
	Name:        "denyteam",
	Type:        1,
	Description: "deny a registered team",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "team_name",
			Description: "name of team you're denying",
			Required:    true,
			MaxLength:   35,
		},
	},
}
