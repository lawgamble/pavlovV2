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
	case "denyteam":
		{
			DenyTeam(s, i, interaction, db)
		}
	}
}

func DenyTeam(s *discordgo.Session, i *discordgo.InteractionCreate, interaction discordgo.ApplicationCommandInteractionData, db mariadb.DBHandler) {
	teamName := interaction.Options[0].Value.(string)
	returnedTeam, _ := db.DB.ReadTeamByTeamName(teamName)

	if returnedTeam.TeamName == "" {
		ApproveDenyResponse(s, i, teamDoesNotExist)
		return
	}
	//should only deny a team if they are "Pending"
	if returnedTeam.TeamStatus != "Pending" {
		ApproveDenyResponse(s, i, notPending)
		return
	}
	//read all players in order to remove their name off temp table associated with denied team
	playersOnTeam, _ := db.DB.ReadAllPlayersOnTempTeam(teamName)

	for _, tempPlayer := range playersOnTeam {
		deletePlayerQuery := fmt.Sprintf("DELETE FROM SND_TEMP_ROSTERS WHERE DiscordId = %s AND Team = %s", tempPlayer.DiscordId, teamName)
		err := db.DB.Delete(deletePlayerQuery)
		if err != nil {
			ApproveDenyResponse(s, i, err.Error())
		}
	}
	team, updateErr := db.DB.UpdateTeamStatus(teamName, "Denied")
	if updateErr != nil {
		ApproveDenyResponse(s, i, updateErr.Error())
	}
	//send msg to team captain that their team was denied.
	teamCaptainString := strconv.FormatInt(team.TeamCaptain, 10)
	_, _ = s.ChannelMessageSend(teamCaptainString, fmt.Sprintf(deniedTeamMessage, teamName))

	ApproveDenyResponse(s, i, teamName+" was denied. Fuck 'em")

}

func ApproveTeam(s *discordgo.Session, i *discordgo.InteractionCreate, interaction discordgo.ApplicationCommandInteractionData, db mariadb.DBHandler) {
	teamName := interaction.Options[0].Value.(string)
	returnedTeam, _ := db.DB.ReadTeamByTeamName(teamName)

	startProcessMsg := fmt.Sprintf("Starting 8 step approval process for %v...", teamName)
	log.Println(startProcessMsg)

	//step 1 - Does Team Exist?
	if returnedTeam.TeamName == "" {
		ApproveDenyResponse(s, i, teamDoesNotExist)
		return
	}
	log.Println("Completed step 1 0f 8 - Team exists in DB")
	//step 2 - Team Status - Must Be "Pending"
	if returnedTeam.TeamStatus != "Pending" {
		ApproveDenyResponse(s, i, notPending)
		return
	}
	log.Println("Completed step 2 0f 8 - Team is in 'Pending' state")
	//step 3 - Team Player Count - Must Be 5
	playersOnTeam, _ := db.DB.ReadAllPlayersOnTempTeam(teamName)
	if len(playersOnTeam) != 5 {
		ApproveDenyResponse(s, i, not5)
		return
	}
	log.Println("Completed step 3 0f 8 - Team has 5 players")
	//step 4 - Players on other Active Teams? - All players must not be on an existing "Active" Team
	listOfPlayersOnAnotherTeam := playersOnOtherTeams(playersOnTeam, db)
	if len(listOfPlayersOnAnotherTeam) > 0 {
		ApproveDenyResponse(s, i, listPlayersOnOtherTeams+fmt.Sprintf("%s", listOfPlayersOnAnotherTeam))
		return
	}
	log.Println("Completed step 4 0f 8 - No players exist on another team")

	//step 5 - Create Team Discord Role - Role name = name of Approved Team
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
		ApproveDenyResponse(s, i, errorTeamRole+teamName)
	}
	log.Println("Completed step 5 0f 8 - Team Role created")

	newTeamRoleId := newTeamRole.ID
	leagueMemberRoleId := os.Getenv("LEAGUE_MEMBER_ROLE_ID")
	teamCaptainRoleId := os.Getenv("TEAM_CAPTAIN_ROLE_ID")

	//step 6 - Update ALL DB tables and give Roles (For each player)
	for _, tempPlayer := range playersOnTeam {
		discordIdString := strconv.FormatInt(tempPlayer.DiscordId, 10)
		// add team name to each player on the USER table
		err := db.DB.UpdatePlayerTeamName(teamName, tempPlayer.DiscordId)
		if err != nil {
			ApproveDenyResponse(s, i, errorPlayerTeamName+discordIdString)
		}
		// give player team role
		err = s.GuildMemberRoleAdd(i.GuildID, discordIdString, newTeamRoleId)
		if err != nil {
			// continue, but let mod know
			ApproveDenyResponse(s, i, errorTeamRoleId+discordIdString)
		}
		// give player "League Member" Role
		err = s.GuildMemberRoleAdd(i.GuildID, discordIdString, leagueMemberRoleId)
		if err != nil {
			// continue, but let mod know
			ApproveDenyResponse(s, i, errorGuildRole+discordIdString)
		}
		// remove all players from temp table, including duplicates
		deleteErr := db.DB.DeletePlayerFromTempTable(discordIdString)
		if deleteErr != nil {
			ApproveDenyResponse(s, i, tempTableDeleteError+discordIdString)
		}
	}
	log.Println("Completed step 6 0f 8 - All players assigned to new Team and assigned roles")

	// change team status on TEAM table to "Active"
	//step 7 - Activate new Team
	team, updateErr := db.DB.UpdateTeamStatus(teamName, "Active")
	if updateErr != nil {
		//tell mod to do it manually
		ApproveDenyResponse(s, i, teamStatusErr)
	}
	log.Println("Completed step 7 0f 8 - Team is now 'Active'")
	// give team captain role to captain on TEAM table
	//step 8 - Give Team Captain the Team Captain Role
	teamCaptainString := strconv.FormatInt(team.TeamCaptain, 10)
	err = s.GuildMemberRoleAdd(i.GuildID, teamCaptainString, teamCaptainRoleId)
	if err != nil {
		// continue, but let mod know user needs role
		ApproveDenyResponse(s, i, teamCaptainRoleError+teamCaptainString)
	}
	log.Println("Completed step 8 0f 8 - Team Captain role assigned")
	// If we get here - SUCCESS!
	ApproveDenyResponse(s, i, successfullyApproved+teamName)
}

func ApproveDenyResponse(s *discordgo.Session, i *discordgo.InteractionCreate, message string) {
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
