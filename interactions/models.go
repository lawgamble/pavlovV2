package interactions

import (
	"github.com/bwmarrin/discordgo"
)

var RegionActionRow = discordgo.ActionsRow{
	Components: []discordgo.MessageComponent{
		discordgo.TextInput{
			CustomID:    "region",
			Label:       "Region",
			Style:       discordgo.TextInputShort,
			Placeholder: "1-USA, 2-USA ONLY, 3-EU, 4-EU ONLY",
			Required:    true,
			MaxLength:   1,
			MinLength:   1,
		},
	},
}

var FilledRegionActionRow = discordgo.ActionsRow{
	Components: []discordgo.MessageComponent{
		&discordgo.TextInput{
			CustomID:    "region",
			Label:       "Region",
			Style:       discordgo.TextInputShort,
			Placeholder: "1-USA, 2-USA ONLY, 3-EU, 4-EU ONLY",
			Required:    true,
			MaxLength:   1,
			MinLength:   1,
		},
	},
}

var PlayStyleActionRow = discordgo.ActionsRow{
	Components: []discordgo.MessageComponent{
		discordgo.TextInput{
			CustomID:    "playStyle",
			Label:       "Play-Style?",
			Style:       discordgo.TextInputShort,
			Placeholder: "1-Flex, 2-Rush/Entry, 3-Lurker, 4-Mid Player, 5-IGL",
			Required:    true,
			MaxLength:   1,
			MinLength:   1,
		},
	},
}

var FilledPlayStyleActionRow = discordgo.ActionsRow{
	Components: []discordgo.MessageComponent{
		&discordgo.TextInput{
			CustomID:    "playStyle",
			Label:       "Play-Style?",
			Style:       discordgo.TextInputShort,
			Placeholder: "1-Flex, 2-Rush/Entry, 3-Lurker, 4-Mid Player, 5-IGL",
			Required:    true,
			MaxLength:   1,
			MinLength:   1,
		},
	},
}

var ActivityActionRow = discordgo.ActionsRow{
	Components: []discordgo.MessageComponent{
		discordgo.TextInput{
			CustomID:    "activity",
			Label:       "Registration Type",
			Style:       discordgo.TextInputShort,
			Placeholder: "1-Enlisted/Draft, 2-Pickups Only, 3-Inactive",
			Required:    true,
			MaxLength:   1,
			MinLength:   1,
		},
	},
}

var FilledActivityActionRow = discordgo.ActionsRow{
	Components: []discordgo.MessageComponent{
		&discordgo.TextInput{
			CustomID:    "activity",
			Label:       "Registration Type",
			Style:       discordgo.TextInputShort,
			Placeholder: "1-Enlisted/Draft, 2- Pickups Only, 3-Inactive",
			Required:    true,
			MaxLength:   1,
			MinLength:   1,
		},
	},
}

var DOBActionRow = discordgo.ActionsRow{
	Components: []discordgo.MessageComponent{
		discordgo.TextInput{
			CustomID:    "dateOfBirth",
			Label:       "Date Of Birth",
			Style:       discordgo.TextInputShort,
			Placeholder: "Ex: MMDDYYYY; 09291988",
			Required:    true,
			MaxLength:   8,
		},
	},
}

var FilledDOBActionRow = discordgo.ActionsRow{
	Components: []discordgo.MessageComponent{
		&discordgo.TextInput{
			CustomID:    "dateOfBirth",
			Label:       "Date Of Birth",
			Style:       discordgo.TextInputShort,
			Placeholder: "Ex: MMDDYYYY; 09291988",
			Required:    true,
			MaxLength:   8,
		},
	},
}

var InGameNameActionRow = discordgo.ActionsRow{
	Components: []discordgo.MessageComponent{
		discordgo.TextInput{
			CustomID:    "inGameName",
			Label:       "In-Game Name",
			Style:       discordgo.TextInputShort,
			Placeholder: "Ex: WinnerWinner42069",
			Required:    true,
			MaxLength:   25,
			MinLength:   1,
		},
	},
}

var FilledInGameNameActionRow = discordgo.ActionsRow{
	Components: []discordgo.MessageComponent{
		&discordgo.TextInput{
			CustomID:  "inGameName",
			Label:     "In-Game Name",
			Style:     discordgo.TextInputShort,
			Required:  true,
			MaxLength: 25,
			MinLength: 1,
		},
	},
}
