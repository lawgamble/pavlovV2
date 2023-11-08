package mariadb

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
	"time"
)

type DBI interface {
	CreatePlayer(m ModalSubmitData, time time.Time) error
	ReadByDiscordId(discordId string) (Player, error)
	Update(m ModalSubmitData) error
	Delete()
}

type DBHandler struct {
	DB DBI
}

type MariaDB struct {
	DB *sql.DB
}
type DBConfig struct {
	Name     string
	Host     string
	Port     string
	User     string
	Password string
}

// NewMariaDB initializes a new MariaDB connection. We should pass in the config struct that is created and initialized using env variables
func NewMariaDB(c DBConfig) (DBI, error) {
	connURL := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v", c.User, c.Password, c.Host, c.Port, c.Name)
	db, err := sql.Open("mysql", connURL)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &MariaDB{DB: db}, nil
}

func BuildConfig() DBConfig {
	return DBConfig{
		Name:     os.Getenv("DBName"),
		Host:     os.Getenv("DBHOST"),
		Port:     os.Getenv("DBPORT"),
		User:     os.Getenv("DBUSER"),
		Password: os.Getenv("DBPASS"),
	}
}

func (db MariaDB) CreatePlayer(m ModalSubmitData, time time.Time) error {
	query := "INSERT INTO SND_USERS (DiscordId, DiscordName, DOB ,PlayStyle, Region, InGameName, PlayerType, RegistrationDate) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"
	_, err := db.DB.Exec(query, m.DiscordId, m.DiscordName, m.DOB, m.PlayStyle, m.Region, m.InGameName, m.PlayerType, time)
	if err != nil {
		log.Println(err)
	}
	return err
}

//ReadByDiscordId calls mariaDb with the user's DiscordId as the Primary Key
func (db MariaDB) ReadByDiscordId(discordId string) (Player, error) {
	query := fmt.Sprintf("SELECT * FROM SND_USERS WHERE DiscordId = %v", discordId)
	row := db.DB.QueryRow(query)

	var player Player
	var team sql.NullString // Declare a sql.NullString for the Team column

	err := row.Scan(
		&player.DiscordId,
		&player.DiscordName,
		&player.DOB,
		&player.PlayStyle,
		&player.Region,
		&player.InGameName,
		&team, // Scan into the sql.NullString for the Team column
		&player.PlayerType,
		&player.RegistrationDate,
	)

	if err == sql.ErrNoRows {
		log.Println("Player not found: " + discordId)
		return player, err
	} else if err != nil {
		log.Println(err)
		return player, err
	}

	// Check if team.Valid is true to see if 'Team' is not NULL
	if team.Valid {
		player.Team = team.String
	} else {
		player.Team = "" // Set a default value for 'Team' when it's NULL
	}
	log.Println(discordId + " is currently registered")
	return player, nil
}

func (db MariaDB) Update(m ModalSubmitData) error {
	query := `
        UPDATE SND_USERS
        SET DOB = ?,
            PlayStyle = ?,
            Region = ?,
            InGameName = ?,
            PlayerType = ?
        WHERE DiscordId = ?`

	_, err := db.DB.Exec(query, m.DOB, m.PlayStyle, m.Region, m.InGameName, m.PlayerType, m.DiscordId)
	if err != nil {
		return err
	}

	log.Println("User data updated successfully.")
	return nil
}

func (db MariaDB) Delete() {}
