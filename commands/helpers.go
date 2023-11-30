package commands

import (
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	mariadb "pfc2/mariaDB"
	"strings"
	"time"
)

func validateRegistrationData(c []discordgo.MessageComponent) (bool, error, mariadb.ModalSubmitData) {
	modalSubmitData := mariadb.ModalSubmitData{
		Region:     c[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value,
		PlayStyle:  c[1].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value,
		PlayerType: c[2].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value,
		DOB:        c[3].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value,
		InGameName: c[4].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value,
	}
	// first, validate age, so we don't have to validate more if we don't have to
	err := modalSubmitData.ValidateDateOfBirth()
	if err != nil {
		return false, err, modalSubmitData
	}
	err = modalSubmitData.ValidatePlayerType()
	if err != nil {
		return false, err, modalSubmitData
	}

	err = modalSubmitData.ValidatePlayStyle()
	if err != nil {
		return false, err, modalSubmitData
	}

	err = modalSubmitData.ValidateRegion()
	if err != nil {
		return false, err, modalSubmitData
	}
	err = modalSubmitData.ValidateInGameName()
	if err != nil {
		return false, err, modalSubmitData
	}
	return true, nil, modalSubmitData
}

// Status returns a random string of the bot's "status" when a user runs the /status command.
func Status(s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: statusResponse(),
		},
	},
	)
	if err != nil {
		log.Print(err)
	}
}

//statusResponse returns a response string for the Status() func - at random.
func statusResponse() string {
	rand.Seed(time.Now().UnixNano())
	randomIndex := rand.Intn(len(statusOptions))
	return statusOptions[randomIndex]
}

func writeToAliasFile(discordId string, inGameName string) error {
	filePath := os.Getenv("ALIASES_FILEPATH")
	// Read the JSON data from the file
	inGameName = "q-" + inGameName
	jsonFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	// Unmarshal the JSON data into a map
	var data map[string]interface{}
	if err := json.Unmarshal(jsonFile, &data); err != nil {
		return err
	}

	// Check if the "players" object exists
	players, ok := data["players"].(map[string]interface{})
	if !ok {
		players = make(map[string]interface{})
		data["players"] = players
	}

	// Check if the discordId already exists in "players"
	existingName, exists := players[discordId].(string)
	if exists {
		// If inGameName is different, update it
		if existingName != inGameName {
			players[discordId] = inGameName
		}
	} else {
		// Add the new player
		players[discordId] = inGameName
	}

	// Marshal the updated data back to JSON with indentation
	updatedJSON, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	// Write the updated JSON back to the file
	if err := ioutil.WriteFile(filePath, updatedJSON, 0644); err != nil {
		return err
	}

	return nil
}

// UserHasRole takes in a roleId and return a bool to represent whether user has that roleId/role
func UserHasRole(s *discordgo.Session, i *discordgo.InteractionCreate, targetRoleID string) bool {
	// Fetch the member
	member, err := s.GuildMember(os.Getenv("GUILD_ID"), i.Member.User.ID)
	if err != nil {
		return false
	}

	// Iterate through the roles of the member and check against target role
	for _, roleID := range member.Roles {
		if roleID == targetRoleID {
			return true
		}
	}
	// Member does not have target role
	return false
}

func generateNodeRepresentation(teams []mariadb.PendingTeam) string {
	var result strings.Builder

	// Build header
	result.WriteString("```\n")

	// Build data
	for _, team := range teams {
		result.WriteString(fmt.Sprintf("\n%s  --%s\n____________\n\n", team.Team, team.Region))
		for i, playerName := range team.PlayerNames {
			result.WriteString(fmt.Sprintf("   * %s (%s)  %s\n", playerName, team.InGameNames[i], team.PlayerDOB[i]))
		}
	}

	result.WriteString("```\n")

	return result.String()
}
