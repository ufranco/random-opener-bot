package bot

import (
	"context"
	"log"
	"math/rand"
	"os"
	"regexp"

	"github.com/bwmarrin/discordgo"
)

var (
	openerRepository   Repository
	applicationContext context.Context
	helpPattern        string
)

const (
	openerPattern string = "^\\!(randomOpener|ro)$"
	colorMaxValue int    = 16777215
)

func Start(
	inheritedContext context.Context,
	session *discordgo.Session,
	repository Repository,
) {

	applicationContext = inheritedContext
	user, err := session.User("@me")

	if err != nil {
		log.Fatal(err.Error())
		os.Exit(1)
	}

	helpPattern = "^(<@!" + user.ID + ">)((\\shelp){0,1})$"

	openerRepository = repository

	session.AddHandler(messageHandler)

	err = session.Open()

	if err != nil {
		log.Fatal(err.Error())
		os.Exit(1)
	}

	log.Printf("RandomOpenerBot running!\nPress Ctrl-C to exit")

}

func messageHandler(session *discordgo.Session, message *discordgo.MessageCreate) {
	if message.Author.ID == session.State.User.ID || message.GuildID == "" {
		return
	}

	matchesHelpPattern, err := regexp.MatchString(helpPattern, message.Content)

	if err != nil {
		log.Fatal(err.Error())
		return
	}

	log.Print(message.Content, "  ", helpPattern)

	if matchesHelpPattern {
		session.ChannelMessageSend(message.ChannelID, "I got ya fam! Type ***!randomOpener*** or just ***!ro***")
	}

	matchesOpenerPattern, err := regexp.MatchString(openerPattern, message.Content)

	if err != nil {
		log.Fatal(err.Error())
		return
	}

	if matchesOpenerPattern {

		nickname := getUserNickname(message)
		log.Printf("%s: %s", nickname, message.Content)

		openerFound, err := openerRepository.GetRandomOpener()

		if err != nil {
			log.Fatal(err.Error())
			return
		}

		embedImage := discordgo.MessageEmbedImage{
			URL: openerFound.ImageURL,
		}

		embedMessage := discordgo.MessageEmbed{
			Title: openerFound.Name,
			Image: &embedImage,
			Type:  "image",
			Footer: &discordgo.MessageEmbedFooter{
				Text: openerFound.Description,
			},
			Color: rand.Intn(colorMaxValue),
		}

		session.ChannelMessageSendEmbed(message.ChannelID, &embedMessage)

	}
}

func getUserNickname(message *discordgo.MessageCreate) string {

	if message.Member.Nick == "" {
		return message.Author.Username
	}

	return message.Member.Nick
}
