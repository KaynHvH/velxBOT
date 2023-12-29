package main

import (
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"log"
	"os"
	"os/signal"
	"syscall"
	"velxBOT/commands"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}

	velx, err := discordgo.New("Bot " + os.Getenv("TOKEN"))
	if err != nil {
		log.Fatal("Can't create Discord session, ", err)
		return
	}

	velx.AddHandler(commands.HandleFunCommands)
	velx.AddHandler(commands.HandleModerationCommands)
	velx.AddHandler(commands.HandleHelpCommand)
	velx.Identify.Intents = discordgo.IntentsAllWithoutPrivileged

	err = velx.Open()
	if err != nil {
		log.Fatal("Can't open bot, ", err)
		return
	}
	defer velx.Close()

	//Cleanly close down the Discord session
	log.Println("Bot is running!")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}
