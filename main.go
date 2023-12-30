package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"velxBOT/commands"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func main() {
	if err := loadEnv(); err != nil {
		log.Fatalf("Error loading environment variables: %s", err)
	}

	velx, err := createDiscordSession()
	if err != nil {
		log.Fatalf("Error creating Discord session: %s", err)
		return
	}
	defer velx.Close()

	h := commands.NewHandler(velx)
	registerHandlers(velx, h)

	startBot(velx)
}

func loadEnv() error {
	return godotenv.Load(".env")
}

func createDiscordSession() (*discordgo.Session, error) {
	token := "Bot " + os.Getenv("TOKEN")
	return discordgo.New(token)
}

func registerHandlers(velx *discordgo.Session, h *commands.Handler) {
	velx.AddHandler(h.HandleHelpCommand)
	velx.AddHandler(h.HandleModerationCommands)
	velx.Identify.Intents = discordgo.IntentsAllWithoutPrivileged
}

func startBot(velx *discordgo.Session) {
	err := velx.Open()
	if err != nil {
		log.Fatalf("Error opening bot: %s", err)
		return
	}

	log.Println("Bot is running!")
	waitForTerminationSignal()
}

func waitForTerminationSignal() {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}
