package interactions

import (
	"github.com/bwmarrin/discordgo"
	"log"
)

func GiveRoleToUser(s *discordgo.Session, i *discordgo.InteractionCreate, roleId string) {
	err := s.GuildMemberRoleAdd(i.GuildID, i.Member.User.ID, roleId)
	if err != nil {
		log.Printf("User %v was not given the role with roleId of %v: %v", i.Member.User.Username, roleId, err)
	}
}

func RemoveRoleFromUser(s *discordgo.Session, i *discordgo.InteractionCreate, roleId string) {
	err := s.GuildMemberRoleRemove(i.GuildID, i.Member.User.ID, roleId)
	if err != nil {
		log.Printf("RoleId of %v for user %v was not removed: %v", roleId, i.Member.User.Username, err)
	}
}
