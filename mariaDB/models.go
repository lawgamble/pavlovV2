package mariadb

import (
	"fmt"
	"strings"
	"time"
)

type Player struct {
	DiscordId        int64   `json:"DiscordId"`
	DiscordName      string  `json:"DiscordName"`
	DOB              []uint8 `json:"DOB"`
	PlayStyle        string  `json:"PlayStyle"`
	Region           string  `json:"Region"`
	InGameName       string  `json:"InGameName"`
	Team             string  `json:"Team"`
	PlayerType       string  `json:"PlayerType"`
	RegistrationDate []uint8 `json:"RegistrationDate"`
}

type ModalSubmitData struct {
	DiscordId   string
	DiscordName string
	Region      string
	PlayStyle   string
	PlayerType  string
	DOB         string
	InGameName  string
}

type TempTeamMember struct {
	DiscordId int64  `json:"DiscordId"`
	Team      string `json:"Team"`
}

type TempTeam []TempTeamMember

type Team struct {
	TeamId      string `json:"TeamId"`
	TeamName    string `json:"TeamName"`
	TeamStatus  string `json:"TeamStatus"`
	TeamRegion  string `json:"TeamRegion"`
	TeamCaptain int64  `json:"TeamCaptain"`
}

type PendingTeams []PendingTeam

type PendingTeam struct {
	Team          string   `json:"Team"`
	Region        string   `json:"Region"`
	PlayerNames   []string `json:"PlayerNames"`
	PlayStyles    []string `json:"PlayStyles"`
	PlayerRegions []string `json:"PlayerRegions"`
	InGameNames   []string `json:"InGameNames"`
	PlayerTypes   []string `json:"PlayerTypes"`
	PlayerDOB     []string `json:"PlayerDOB"`
}

func (p *PendingTeam) ConvertToSliceOfString(u []uint8) []string {
	strValue := string(u)
	return strings.Split(strValue, ",")
}

func (m *ModalSubmitData) ValidateRegion() error {
	switch m.Region {
	case "1":
		{
			m.Region = "USA"
			break
		}
	case "2":
		{
			m.Region = "USA ONLY"
			break
		}
	case "3":
		{
			m.Region = "EU"
			break
		}
	case "4":
		{
			m.Region = "EU ONLY"
			break
		}
	default:
		return fmt.Errorf("invalid region, try again")
	}
	return nil
}

func (m *ModalSubmitData) ValidatePlayStyle() error {
	switch m.PlayStyle {
	case "1":
		{
			m.PlayStyle = "Flex"
			break
		}
	case "2":
		{
			m.PlayStyle = "Rush/Entry"
			break
		}
	case "3":
		{
			m.PlayStyle = "Lurker"
			break
		}
	case "4":
		{
			m.PlayStyle = "Mid Player"
			break
		}
	case "5":
		{
			m.PlayStyle = "IGL"
			break
		}
	default:
		err := fmt.Errorf("%v is not a valid play style selection", m.PlayStyle)
		return err
	}
	return nil
}

func (m *ModalSubmitData) ValidatePlayerType() error {
	switch m.PlayerType {
	case "1":
		{
			m.PlayerType = "Draftable"
			break
		}
	case "2":
		{
			m.PlayerType = "Team Member"
			break
		}
	case "3":
		{
			m.PlayerType = "Pickups Only"
			break
		}
	default:
		err := fmt.Errorf("%v is not a valid player type selection", m.PlayerType)
		return err
	}
	return nil
}

func (m *ModalSubmitData) ValidateDateOfBirth() error {
	dob, err := time.Parse("01022006", m.DOB)
	if err != nil {
		errFormat := fmt.Errorf("for DOB, you must use the format: MMDDYYYY")
		return errFormat
	}
	currentDate := time.Now()
	age := currentDate.Year() - dob.Year()
	if age < 13 {
		return fmt.Errorf("you're not old enough to register! You must be at least 13")
	}

	m.DOB, err = reformatDOB(m.DOB)
	if err != nil {
		return err
	}
	return nil
}

func reformatDOB(birth string) (string, error) {
	// Parse the input DOB in MMDDYYYY format
	t, err := time.Parse("01022006", birth)
	if err != nil {
		return "", err
	}

	// Format the date in YYYY-MM-DD format
	formattedDOB := t.Format("2006-01-02")

	return formattedDOB, nil
}

func (m *ModalSubmitData) ValidateInGameName() error {
	if m.InGameName == "" {
		return fmt.Errorf("invalid in-game-name")
	}
	return nil
}

func (p Player) IsNotEmpty() bool {
	if p.DiscordId != 0 {
		return true
	}
	return false
}
