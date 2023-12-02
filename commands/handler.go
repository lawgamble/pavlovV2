package commands

import (
	"database/sql"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
	"pfc2/components"
	"pfc2/interactions"
	mariadb "pfc2/mariaDB"
	"strconv"
	"sync"
)

func HandleModalSubmit(s *discordgo.Session, i *discordgo.InteractionCreate, db mariadb.DBHandler) {
	modalSubmitData := i.ModalSubmitData()
	registerRoleId := os.Getenv("REGISTERED_ROLE_ID")
	enlistedRoleId := os.Getenv("ENLISTED_ROLE_ID")
	isValid, err, submitData := validateRegistrationData(modalSubmitData.Components)
	if !isValid {
		interactions.RegistrationErrorResponse(s, i, err)
		return
	}

	switch modalSubmitData.CustomID {
	case "createRegistration":
		{
			err := interactions.SubmitRegistration(i, db, submitData)
			if err != nil {
				//there was an error calling the DB
				log.Print(err)
				interactions.RegistrationErrorResponse(s, i, err)
				break
			}
			interactions.RegistrationSuccessResponse(s, i)
			interactions.GiveRoleToUser(s, i, registerRoleId)

			if submitData.PlayerType == "Draftable" {
				interactions.GiveRoleToUser(s, i, enlistedRoleId)
			}
			err = writeToAliasFile(i.Member.User.ID, submitData.InGameName)
			if err != nil {
				log.Print(err)
			}
			break
		}

	case "updateRegistration":
		{
			submitData.DiscordId = i.Member.User.ID
			err := interactions.SubmitUpdatedRegistration(db, submitData)
			if err != nil {
				//there was an error calling the DB
				log.Print(err)
				interactions.RegistrationErrorResponse(s, i, err)
				break
			}
			interactions.UpdatedRegistrationSuccessResponse(s, i)

			if submitData.PlayerType == "Draftable" {
				interactions.GiveRoleToUser(s, i, enlistedRoleId)
			} else {
				interactions.RemoveRoleFromUser(s, i, enlistedRoleId)
			}
			err = writeToAliasFile(i.Member.User.ID, submitData.InGameName)
			if err != nil {
				log.Print(err)
			}
			break
		}
	}
}

// HandleApplicationCommands handles all slash commands
func HandleApplicationCommands(s *discordgo.Session, i *discordgo.InteractionCreate, db mariadb.DBHandler) {
	interaction := i.ApplicationCommandData()
	switch interaction.Name {
	case "teamregister":
		{
			interactions.RegisterTeam(s, i, interaction, db)
		}
	case "status":
		{
			Status(s, i)
			break
		}
	case "addbutton":
		{
			components.InitializeButtons(s, i, interaction.Options)
			break
		}
	case "repeat":
		{
			HandleRepeatCommand(s, i, interaction)
		}
	case "recoveraccount":
		{
			fmt.Println("You did it!")
		}
	case "listapprovals":
		{
			HandleListApprovalsCommand(s, i, db)
		}
	case "approveteam":
		{
			ApproveTeam(s, i, interaction, db)
		}
	}
}

func ApproveTeam(s *discordgo.Session, i *discordgo.InteractionCreate, interaction discordgo.ApplicationCommandInteractionData, db mariadb.DBHandler) {
	teamName := interaction.Options[0].Value.(string)
	returnedTeam, _ := db.DB.ReadTeamByTeamName(teamName)

	if returnedTeam.TeamName == "" {
		// send error that team does not exist
		ApproveTeamMessageResponse(s, i, teamDoesNotExist)
		return
	}
	// if team exists: // Validate the team is not already in an active state
	if returnedTeam.TeamStatus != "Pending" {
		// bail here as the team is already active or denied, etc.
		ApproveTeamMessageResponse(s, i, notPending)
		return
	}
	playersOnTeam, _ := db.DB.ReadAllPlayersOnTempTeam(teamName)
	if len(playersOnTeam) != 5 {
		// return error here - should be 5
		ApproveTeamMessageResponse(s, i, not5)
		return
	}
	// check that each player does not have a team in the USERS table (TEAM should be empty)
	listOfPlayersOnAnotherTeam := playersOnOtherTeams(playersOnTeam, db)
	if len(listOfPlayersOnAnotherTeam) > 0 {
		// we have to bail here and the team can not be registered
		ApproveTeamMessageResponse(s, i, listPlayersOnOtherTeams+fmt.Sprintf("%s", listOfPlayersOnAnotherTeam))
		return

	}

	// Create DiscordRole "Case sensitive" team name
	mentionable := true
	roleParameters := discordgo.RoleParams{
		Name:        teamName,
		Color:       nil,
		Hoist:       nil,
		Permissions: nil,
		Mentionable: &mentionable,
	}
	newTeamRole, err := s.GuildRoleCreate(i.GuildID, &roleParameters)
	if err != nil {
		// need to say there was an error, but can continue
		ApproveTeamMessageResponse(s, i, errorTeamRole+teamName)
	}

	newTeamRoleId := newTeamRole.ID
	leagueMemberRoleId := os.Getenv("LEAGUE_MEMBER_ROLE_ID")
	teamCaptainRoleId := os.Getenv("TEAM_CAPTAIN_ROLE_ID")

	//if none fail, we are approved
	// loop through all tempPlayers
	for _, tempPlayer := range playersOnTeam {
		discordIdString := strconv.FormatInt(tempPlayer.DiscordId, 10)
		// add team name to each player on the USER table
		err := db.DB.UpdatePlayerTeamName(teamName, tempPlayer.DiscordId)
		if err != nil {
			//probably just continue here and send a res to mod
			ApproveTeamMessageResponse(s, i, errorPlayerTeamName+discordIdString)
		}
		// give player team role
		err = s.GuildMemberRoleAdd(i.GuildID, discordIdString, newTeamRoleId)
		if err != nil {
			// continue, but let mod know
			ApproveTeamMessageResponse(s, i, errorTeamRoleId+discordIdString)
		}
		// give player "League Member" Role
		err = s.GuildMemberRoleAdd(i.GuildID, discordIdString, leagueMemberRoleId)
		if err != nil {
			// continue, but let mod know
			ApproveTeamMessageResponse(s, i, errorGuildRole+discordIdString)
		}
		// remove all players from temp table, including duplicates
		deleteErr := db.DB.DeletePlayerFromTempTable(discordIdString)
		if deleteErr != nil {
			// let mod know they need to do it manually
			ApproveTeamMessageResponse(s, i, tempTableDeleteError+discordIdString)
		}
	}

	// change team status on TEAM table to "Active"
	team, updateErr := db.DB.UpdateTeamStatus(teamName, "Active")
	if updateErr != nil {
		//tell mod to do it manually
		ApproveTeamMessageResponse(s, i, teamStatusErr)
	}
	// give team captain role to captain on TEAM table
	teamCaptainString := strconv.FormatInt(team.TeamCaptain, 10)
	err = s.GuildMemberRoleAdd(i.GuildID, teamCaptainString, teamCaptainRoleId)
	if err != nil {
		// continue, but let mod know user needs role
		ApproveTeamMessageResponse(s, i, teamCaptainRoleError+teamCaptainString)
	}
	// send message to a channel? In game names blah blah.
	ApproveTeamMessageResponse(s, i, successfullyApproved+teamName)
}

func ApproveTeamMessageResponse(s *discordgo.Session, i *discordgo.InteractionCreate, message string) {
	resErr := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: message,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	if resErr != nil {
		log.Print(resErr)
		return
	}
}

func playersOnOtherTeams(team []mariadb.TempTeamMember, db mariadb.DBHandler) []string {
	var registeredList []string
	var wg sync.WaitGroup
	resultCh := make(chan sql.Row, len(team))

	for _, player := range team {
		wg.Add(1)
		go func(player mariadb.TempTeamMember) {
			defer wg.Done()

			// Perform the database query for each player
			onOtherTeam := db.DB.ReadPlayerOnRegisteredTeam(player.DiscordId)

			// Send the result to the channel
			resultCh <- onOtherTeam
		}(player)
	}
	go func() {
		wg.Wait()
		close(resultCh)
	}()

	// Check results from the channel
	for result := range resultCh {
		var player mariadb.Player
		var team sql.NullString // Declare a sql.NullString for the Team column

		err := result.Scan(
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
			continue
		} else {
			registeredList = append(registeredList, player.InGameName)
		}
	}
	return registeredList
}

func HandleListApprovalsCommand(s *discordgo.Session, i *discordgo.InteractionCreate, db mariadb.DBHandler) {
	if !UserHasRole(s, i, os.Getenv("LEAGUE_MANAGER_ROLE_ID")) {
		err := s.InteractionRespond(i.Interaction, interactions.NotPermittedInteractionResponse)
		if err != nil {
			log.Print(err)
		}
		return
	}
	pendingTeams, err := db.DB.ReadAllPendingTeamsAndPlayers()
	if err != nil {
		interactions.DefaultErrorResponse(s, i, err)
	}
	// build func to respond with a table of all teams.
	tempTableString := generateNodeRepresentation(pendingTeams)
	interactions.SendTempTeamsTable(s, i, tempTableString)
}

func HandleRepeatCommand(s *discordgo.Session, i *discordgo.InteractionCreate, interaction discordgo.ApplicationCommandInteractionData) {
	if !UserHasRole(s, i, os.Getenv("LEAGUE_MANAGER_ROLE_ID")) {
		err := s.InteractionRespond(i.Interaction, interactions.NotPermittedInteractionResponse)
		if err != nil {
			log.Print(err)
		}
		return
	}
	// If we get here, user has permission to run command
	chanId := interaction.Options[0].Value.(string)
	message := interaction.Options[1].Value.(string)

	err := s.InteractionRespond(i.Interaction, interactions.RepeatCommandInteractionResponse)
	if err != nil {
		log.Print(err)
	}
	_, err = s.ChannelMessageSend(chanId, message)
	if err != nil {
		log.Print(err)
	}
}
