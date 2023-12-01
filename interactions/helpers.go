package interactions

import (
	"fmt"
	mariadb "pfc2/mariaDB"
	"regexp"
	"strconv"
	"strings"
)

func checkIfActiveTeamAlreadyExists(teamName string, db mariadb.DBHandler) error {
	returnedTeam, _ := db.DB.ReadTeamByTeamName(teamName) // a team can be returned here if needed. For now, err will determine if row was not found
	if returnedTeam.TeamName != "" {
		err := fmt.Errorf(" `%v`\n has already been chosen - try another team name", teamName)
		return err
	}
	return nil
}

func whoIsNotRegistered(players []string) string {
	msg := "Users are not registered: "
	for _, player := range players {
		msg += fmt.Sprintf("<@%s> ", player)
	}
	return strings.TrimSpace(msg)
}

// sendTeamRegistration adds players to a temp table in the DB, along with adding pending team to team table.
func sendTeamRegistration(players []mariadb.Player, teamName, teamRegion string, db mariadb.DBHandler, teamCaptain string) error {
	for _, player := range players {
		// avoiding concurrency here (not sure how DB will handle concurrent writes)
		err := db.DB.CreateTempRoster(strconv.FormatInt(player.DiscordId, 10), teamName)
		if err != nil {
			return err
		}
	}
	err := db.DB.CreateTempTeam(teamName, teamRegion, teamCaptain)
	if err != nil {
		return err
	}
	return nil
}

// hasDuplicates takes in a slice of string and checks to see if there are exact string matches within any index.
func hasDuplicates(slice []string) (bool, string) {
	seen := make(map[string]bool)

	for i, value := range slice {
		// If the value is already in the map, it's a duplicate
		if seen[value] {
			return true, slice[i]
		}

		// Mark the value as seen
		seen[value] = true
	}

	// No duplicates found
	return false, ""
}

func containsOnlyNumbers(input string) bool {
	// Use a regular expression to check if the string contains only numbers
	match, _ := regexp.MatchString("^[0-9]+$", input)
	return match
}
