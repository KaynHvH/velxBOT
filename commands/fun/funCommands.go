package fun

import (
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func HandleFunCommands(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	args := strings.Split(m.Content, " ")
	if args[0] != os.Getenv("PREFIX") {
		return
	}

	switch args[1] {
	case "ping":
		s.ChannelMessageSend(m.ChannelID, "pong")
	case "dog":
		sendRandomDogImage(s, m)
	case "answer":
		fmt.Println("Answer command recognized. Args:", args)
		sendRandomAnswer(s, m, args)
	case "whois":
		fmt.Println("Whois command recognized. Args:", args)
		whoisCommand(s, m, args)
	case "avatar":
		fmt.Println("Avatar command recognized. Args:", args)
		userAvatar(s, m, args)
	case "dice":
		fmt.Println("Dice command recognized. Args:", args)
		rollDice(s, m, args)
	}
}

type ApiResponse struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

func userAvatar(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if len(args) < 3 {
		s.ChannelMessageSend(m.ChannelID, "Please mention a user to check.")
		return
	}

	mention := args[2]
	userID := strings.Trim(mention, "<@!>")
	user, err := s.User(userID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Error retrieving user information.")
		return
	}

	embed := discordgo.MessageEmbed{
		Title: user.Username + " avatar",
		Image: &discordgo.MessageEmbedImage{URL: user.AvatarURL("")},
		Author: &discordgo.MessageEmbedAuthor{
			Name:    "Requested by: " + m.Author.Username,
			IconURL: m.Author.AvatarURL(""),
		},
	}
	s.ChannelMessageSendEmbed(m.ChannelID, &embed)
}

func sendRandomDogImage(s *discordgo.Session, m *discordgo.MessageCreate) {
	url := "https://dog.ceo/api/breeds/image/random"
	response, err := http.Get(url)
	if err != nil {
		log.Println(err)
		return
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
		return
	}

	var apiResponse ApiResponse
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		log.Println(err)
		return
	}

	embed := discordgo.MessageEmbed{
		Title: "Random dog image",
		Author: &discordgo.MessageEmbedAuthor{
			Name:    "Requested by: " + m.Author.Username,
			IconURL: m.Author.AvatarURL(""),
		},
		Image: &discordgo.MessageEmbedImage{
			URL: apiResponse.Message,
		},
	}
	s.ChannelMessageSendEmbed(m.ChannelID, &embed)
}

func sendRandomAnswer(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	var answers = []string{
		"Yes", "Of course", "No", "Maybe", "I don't know", "Probably yes", "Probably not", "I don't think so",
	}
	randIndex := rand.Intn(len(answers))
	question := strings.Join(args[2:len(args)], " ")

	embed := discordgo.MessageEmbed{
		Title:       "User " + m.Author.Username + " asked: " + question,
		Description: "Answer: **" + answers[randIndex] + "**",
		Author: &discordgo.MessageEmbedAuthor{
			Name:    "Requested by: " + m.Author.Username,
			IconURL: m.Author.AvatarURL(""),
		},
	}
	s.ChannelMessageSendEmbed(m.ChannelID, &embed)
	return
}

func whoisCommand(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if len(args) < 3 {
		s.ChannelMessageSend(m.ChannelID, "Please mention a user to check.")
		return
	}

	mention := args[2]
	userID := strings.Trim(mention, "<@!>")
	user, err := s.User(userID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Error retrieving user information.")
		return
	}

	embed := discordgo.MessageEmbed{
		Title: "User Information",
		Author: &discordgo.MessageEmbedAuthor{
			Name:    "Requested by: " + m.Author.Username,
			IconURL: m.Author.AvatarURL(""),
		},
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: user.AvatarURL(""),
		},
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Username",
				Value:  user.Username,
				Inline: false,
			},
			{
				Name:   "ID",
				Value:  user.ID,
				Inline: false,
			},
			{
				Name:   "MFA",
				Value:  strconv.FormatBool(user.MFAEnabled),
				Inline: false,
			},
			{
				Name:   "Is bot",
				Value:  strconv.FormatBool(user.Bot),
				Inline: false,
			},
		},
	}

	s.ChannelMessageSendEmbed(m.ChannelID, &embed)
	return
}

func rollDice(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	rand.Seed(time.Now().UnixNano())
	randomNumber := rand.Intn(6)
	var dice string = "https://skalkuluj.pl/assets/images/dice" + strconv.Itoa(randomNumber) + ".webp"

	embed := discordgo.MessageEmbed{
		Title: "Rolling the dice...",
		Image: &discordgo.MessageEmbedImage{URL: dice},
		Author: &discordgo.MessageEmbedAuthor{
			Name:    "Requested by: " + m.Author.Username,
			IconURL: m.Author.AvatarURL(""),
		},
	}
	s.ChannelMessageSendEmbed(m.ChannelID, &embed)
	return
}
