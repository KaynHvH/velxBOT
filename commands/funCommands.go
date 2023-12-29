package commands

import (
	"github.com/bwmarrin/discordgo"
	"math/rand"
	"strconv"
	"strings"
)

var dogs = []string{
	"https://zooart.com.pl/blog/wp-content/uploads/2020/09/alaskan-malamute-1.jpg",
	"https://zooart.com.pl/blog/wp-content/uploads/2020/09/alaskan-malamute-3.jpg",
	"https://i0.wp.com/piesrasowy.pl/wp-content/uploads/2019/06/4.jpg?fit=1200%2C707&ssl=1",
	"https://zielonyogrodek.pl/i/2020/03/16/81456-b56e-1800x0-sc1x77778_siberian-husky-fot-riley-sullivan-unsplash.jpg",
	"https://puppyintraining.com/wp-content/uploads/german-shepherd-growth-chart.jpg",
	"https://warsawdog.com/wp-content/uploads/2021/08/kangal-1-1024x683.jpg.webp",
	"https://paradepets.com/.image/t_share/MTkxMzY1Nzg4NjczMzIwNTQ2/cutest-dog-breeds-jpg.jpg",
	"https://www.princeton.edu/sites/default/files/styles/1x_full_2x_half_crop/public/images/2022/02/KOA_Nassau_2697x1517.jpg?itok=Bg2K7j7J",
	"https://i.pinimg.com/originals/72/83/6c/72836ce645389ce7fd6435ffe844dcab.jpg",
	"https://www.southernliving.com/thmb/Rz-dYEhwq_82C5_Y9GLH2ZlEoYw=/1500x0/filters:no_upscale():max_bytes(150000):strip_icc()/gettyimages-837898820-1-4deae142d4d0403dbb6cb542bfc56934.jpg",
}

var answers = []string{
	"Yes", "Of course", "No", "Maybe", "I don't know", "Probably yes", "Probably not", "I don't think so",
}

func HandleFunCommands(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	args := strings.Split(m.Content, " ")
	if args[0] != PREFIX {
		return
	}

	switch args[1] {
	case "ping":
		s.ChannelMessageSend(m.ChannelID, "sronk")
	case "dog":
		sendRandomDogImage(s, m)
	case "answer":
		sendRandomAnswer(s, m, args)
	case "whois":
		whoisCommand(s, m, args)
	}
}

func sendRandomDogImage(s *discordgo.Session, m *discordgo.MessageCreate) {
	randIndex := rand.Intn(len(dogs))

	embed := discordgo.MessageEmbed{
		Title: "Random dog image",
		Author: &discordgo.MessageEmbedAuthor{
			Name:    "Requested by: " + m.Author.Username,
			IconURL: m.Author.AvatarURL(""),
		},
		Image: &discordgo.MessageEmbedImage{
			URL: dogs[randIndex],
		},
	}
	s.ChannelMessageSendEmbed(m.ChannelID, &embed)
}

func sendRandomAnswer(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
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
}
