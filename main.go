package main

import (
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"log"
	"os"
	"os/signal"
	"pfc2/commands"
	"pfc2/interactions"
	mariadb "pfc2/mariaDB"
	"syscall"
)

var dbHandler mariadb.DBHandler

func main() {
	// make sure to change file name below to "./variables.env"
	_ = godotenv.Load("./vars.env")

	bot, err := discordgo.New(os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Panic(err)
		return
	}

	dbConfig := mariadb.BuildConfig()
	db, err := mariadb.NewMariaDB(dbConfig)

	dbHandler.DB = db

	if err != nil {
		log.Panic(err)
	}

	//Register all slashCommands
	_, err = bot.ApplicationCommandBulkOverwrite(os.Getenv("BOT_ID"), os.Getenv("GUILD_ID"), commands.SlashCommands)
	if err != nil {
		log.Panic(err)
	}

	bot.AddHandler(readyHandler)
	bot.AddHandler(commandHandler)
	err = bot.Open()
	if err != nil {
		log.Println("Error opening connection:", err)
		return
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-signalChan

	bot.Close()

}

func readyHandler(s *discordgo.Session, event *discordgo.Ready) {
	log.Println("Bot is now running. Press Ctrl+C to exit.")
}

func commandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionMessageComponent:
		{
			interactions.HandleMessageComponent(s, i, dbHandler)
			break
		}
	case discordgo.InteractionModalSubmit:
		{
			commands.HandleModalSubmit(s, i, dbHandler)
			break
		}
	case discordgo.InteractionApplicationCommand:
		{
			commands.HandleApplicationCommands(s, i)
			break
		}
	}
}
