package main

import (
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"log"
	"os"
	"os/signal"
	"syscall"
	"velxBOT/commands/fun"
	"velxBOT/commands/help"
	"velxBOT/commands/moderation"
)

func main() {

	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}

	velx, err := discordgo.New("Bot " + os.Getenv("TOKEN"))
	if err != nil {
		log.Fatal("Can't create Discord session, ", err)
		return
	}

	//Handlers
	velx.AddHandler(fun.HandleFunCommands)
	velx.AddHandler(moderation.HandleModerationCommands)
	velx.AddHandler(help.HandleHelpCommand)

	velx.Identify.Intents = discordgo.IntentsAllWithoutPrivileged

	err = velx.Open()
	if err != nil {
		log.Fatal("Can't open bot, ", err)
		return
	}
	defer velx.Close()

	//Cleanly close down the Discord session
	log.Println("Bot is running!")
	waitForTerminationSignal()
}

func waitForTerminationSignal() {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}
