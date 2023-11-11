package interactions

import (
	"database/sql"
	"github.com/bwmarrin/discordgo"
	mariadb "pfc2/mariaDB"
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
