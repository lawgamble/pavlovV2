package mariadb

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type DBI interface {
	CreatePlayer(m ModalSubmitData, time time.Time) error
	ReadAllPlayersOnTempTeam(teamName string) ([]TempTeamMember, error)
	ReadPlayerOnRegisteredTeam(playerId int64) sql.Row
	CreateTempRoster(discordId, teamName string) error
	ReadTeamByTeamName(teamName string) (Team, error)
	ReadAllPendingTeamsAndPlayers() (PendingTeams, error)
	CreateTempTeam(teamName, teamRegion, teamCaptain string) error
	ReadUsersByDiscordId(discordId string) (Player, error)
	Update(m ModalSubmitData) error
	UpdateTeamStatus(teamName, status string) (Team, error)
	UpdatePlayerTeamName(teamName string, discordId int64) error
	DeletePlayerFromTempTable(discordId string) error
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

func (db MariaDB) ReadPlayerOnRegisteredTeam(playerId int64) sql.Row {
	query := fmt.Sprintf("SELECT * FROM SND_USERS su WHERE DiscordId = %d AND TeamName IS NOT NULL AND TeamName != ''", playerId)
	row := db.DB.QueryRow(query)
	return *row
}

func (db MariaDB) CreatePlayer(m ModalSubmitData, time time.Time) error {
	query := "INSERT INTO SND_USERS (DiscordId, DiscordName, DOB ,PlayStyle, Region, InGameName, PlayerType, RegistrationDate) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"
	_, err := db.DB.Exec(query, m.DiscordId, m.DiscordName, m.DOB, m.PlayStyle, m.Region, m.InGameName, m.PlayerType, time)
	if err != nil {
		log.Println(err)
	}
	return err
}

//ReadUsersByDiscordId calls mariaDb with the user's DiscordId as the Primary Key
func (db MariaDB) ReadUsersByDiscordId(discordId string) (Player, error) {
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

func (db MariaDB) UpdateTeamStatus(teamName, status string) (Team, error) {
	var team Team
	query := `
		UPDATE SND_TEAMS
		SET TeamStatus = ?
		WHERE TeamName = ?`
	row := db.DB.QueryRow(query, status, teamName)

	err := row.Scan(
		&team.TeamId,
		&team.TeamName,
		&team.TeamStatus,
		&team.TeamRegion,
		&team.TeamCaptain,
	)

	return team, err
}

func (db MariaDB) UpdatePlayerTeamName(teamName string, discordId int64) error {
	query := `
        UPDATE SND_USERS
        SET TeamName = ?
        WHERE DiscordId = ?`

	_, err := db.DB.Exec(query, teamName, discordId)
	if err != nil {
		return err
	}

	return nil
}

func (db MariaDB) DeletePlayerFromTempTable(discordId string) error {
	query := fmt.Sprintf("DELETE FROM SND_TEMP_ROSTERS WHERE DiscordId = %s", discordId)
	_, err := db.DB.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func (db MariaDB) CreateTempRoster(discordId, teamName string) error {
	query := "INSERT INTO SND_TEMP_ROSTERS (DiscordId, Team) VALUES (?, ?)"
	_, err := db.DB.Query(query, discordId, teamName)
	return err
}

func (db MariaDB) CreateTempTeam(teamName, teamRegion, teamCaptain string) error {
	query := "INSERT INTO SND_TEAMS (TeamName, TeamRegion, TeamStatus, TeamCaptain) VALUES (?, ?, ?, ?)"
	_, err := db.DB.Query(query, teamName, teamRegion, "Pending", teamCaptain)
	return err
}

func (db MariaDB) ReadTeamByTeamName(teamName string) (Team, error) {
	var team Team

	query := "SELECT * FROM SND_TEAMS WHERE LOWER(TeamName) = LOWER(?) LIMIT 1"
	row := db.DB.QueryRow(query, teamName)

	err := row.Scan(
		&team.TeamId,
		&team.TeamName,
		&team.TeamStatus,
		&team.TeamRegion,
		&team.TeamCaptain,
	)

	return team, err
}

func (db MariaDB) ReadAllPlayersOnTempTeam(teamName string) ([]TempTeamMember, error) {
	var tempMembers []TempTeamMember

	query := "SELECT * FROM SND_TEMP_ROSTERS WHERE Team = (?)"
	rows, err := db.DB.Query(query, teamName)
	if err != nil {
		return []TempTeamMember{}, err
	}

	for rows.Next() {
		var tempPlayer TempTeamMember
		var discordId, team string

		err := rows.Scan(&discordId, &team)
		if err != nil {
			return []TempTeamMember{}, err
		}
		discordIdInt, _ := strconv.ParseInt(discordId, 10, 64)
		tempPlayer.DiscordId = discordIdInt
		tempPlayer.Team = team

		tempMembers = append(tempMembers, tempPlayer)
	}
	return tempMembers, err
}

func (db MariaDB) ReadAllPendingTeamsAndPlayers() (PendingTeams, error) {
	var allTeams PendingTeams
	rows, err := db.DB.Query(tempTeamQuery)
	if err != nil {
		return allTeams, err
	}
	defer rows.Close()

	for rows.Next() {
		var t PendingTeam
		var playerNames, playStyle, regions, inGameName, playerTypes, dob []uint8
		err := rows.Scan(
			&t.Team,
			&t.Region,
			&playerNames,
			&playStyle,
			&regions,
			&inGameName,
			&playerTypes,
			&dob,
		)
		if err != nil {
			return allTeams, err
		}

		age, _ := convertDOBListToAgeList(dob)

		t.PlayerNames = t.ConvertToSliceOfString(playerNames)
		t.PlayStyles = t.ConvertToSliceOfString(playStyle)
		t.PlayerRegions = t.ConvertToSliceOfString(regions)
		t.InGameNames = t.ConvertToSliceOfString(inGameName)
		t.PlayerTypes = t.ConvertToSliceOfString(playerTypes)
		t.PlayerDOB = age
		allTeams = append(allTeams, t)
	}
	return allTeams, nil
}

func convertDOBListToAgeList(dobList []uint8) ([]string, error) {
	// Convert []uint8 to string
	dobStr := string(dobList)

	// Split the string by commas
	dobStrings := strings.Split(dobStr, ",")

	// Initialize the age list
	var ageList []string

	// Loop through each DOB string and convert to age
	for _, dob := range dobStrings {
		trimmedDOB := strings.TrimSpace(dob)
		age, err := convertPlayerDOBtoAge([]uint8(trimmedDOB))
		if err != nil {
			return nil, err
		}
		ageList = append(ageList, age)
	}

	return ageList, nil
}

func convertPlayerDOBtoAge(dob []uint8) (string, error) {
	// Convert []uint8 to string
	dobStr := string(dob)

	// Parse the date of birth
	parsedDOB, err := time.Parse("2006-01-02", dobStr)
	if err != nil {
		return "", err
	}

	// Calculate age
	age := calculateAge(parsedDOB)

	// Return age as a string
	return fmt.Sprintf("%d", age), nil
}

func calculateAge(dob time.Time) int {
	// Get current time
	currentTime := time.Now()

	// Calculate age
	age := currentTime.Year() - dob.Year()

	// Adjust age if birthday hasn't occurred yet this year
	if currentTime.YearDay() < dob.YearDay() {
		age--
	}

	return age
}
