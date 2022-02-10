package bot

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var (
	logicLayer         Service
	applicationContext context.Context
	helpPattern        string
)

const (
	getRandomOpenerPattern  string = `^\!(randomOpener|ro)$`
	getLeaderboardPattern   string = `^\!(randomOpener|ro) top$`
	setFavoriteOpenerPatter string = `^\!(randomOpener favorite|ro fav|rof)\s([a-zA-Z\s)]{3,20})$`

	colorMaxValue int = 16777215
)

func Start(
	session *discordgo.Session,
	logic Service,
) {
	logicLayer = logic
	user, err := session.User("@me")

	if err != nil {
		log.Fatal(err.Error())
		os.Exit(1)
	}

	helpPattern = "^(<@!" + user.ID + ">)((\\shelp){0,1})$"

	session.AddHandler(messageCreateHandler)
	//session.AddHandler(reactionHandler)

	err = session.Open()

	if err != nil {
		log.Fatal(err.Error())
		os.Exit(1)
	}

	log.Printf("RandomOpenerBot running!\nPress Ctrl-C to exit")

}

/* func reactionHandler(session *discordgo.Session, reaction *discordgo.MessageReactionAdd) {

	log.Printf("%s reacted with '%s'", reaction.UserID, reaction.Emoji.Name)

	if reaction.GuildID == "" || reaction.UserID == session.State.User.ID {
		return
	}

	log.Printf("Channel: %s \nMessage: %s", reaction.ChannelID, reaction.MessageID)

	//TODO: this does not work help
	reactedTo, err := session.State.Message(reaction.ChannelID, reaction.MessageID)

	if err != nil {
		log.Fatal(err.Error())
		return
	}

	if reactedTo.Author.ID != session.State.User.ID {
		return
	}

	logicLayer.ProcessReaction(reaction, reactedTo)
	session.ChannelMessageSend(reaction.ChannelID, "bancá que todavia no se procesar emociones :c")

}
*/

func messageCreateHandler(session *discordgo.Session, message *discordgo.MessageCreate) {
	if message.GuildID == "" || message.Author.ID == session.State.User.ID {
		return
	}

	log.Printf("%s said '%s'", getUserNickname(message), message.Content)

	matchesHelpPattern, err := regexp.MatchString(helpPattern, message.Content)

	if err != nil {
		log.Fatal(err.Error())
		return
	}

	matchesGetRandomOpenerPattern, err := regexp.MatchString(getRandomOpenerPattern, message.Content)

	if err != nil {
		log.Fatal(err.Error())
		return
	}

	matchesGetLeaderboardPattern, err := regexp.MatchString(getLeaderboardPattern, message.Content)

	if err != nil {
		log.Fatal(err.Error())
		return
	}

	matchesSetFavoriteOpenerPattern, err := regexp.MatchString(setFavoriteOpenerPatter, message.Content)

	if err != nil {
		log.Fatal(err.Error())
		return
	}

	if matchesHelpPattern {
		session.ChannelMessageSend(message.ChannelID, "I got ya fam! Type ***!randomOpener*** or just ***!ro***")
	}

	if matchesGetRandomOpenerPattern {

		opener, err := logicLayer.GetRandomOpener(message)

		if err != nil {
			session.ChannelMessageSend(message.ChannelID, "se rompió todo, bancame un cachito uwu")
			return
		}

		embedImage := discordgo.MessageEmbedImage{
			URL: opener.ImageURL,
		}

		embedMessage := discordgo.MessageEmbed{
			Title: opener.Name,
			Image: &embedImage,
			Type:  "image",
			Footer: &discordgo.MessageEmbedFooter{
				Text: opener.Description,
			},
			Color: rand.Intn(colorMaxValue),
		}

		session.ChannelMessageSendEmbed(message.ChannelID, &embedMessage)
		return
	}

	if matchesGetLeaderboardPattern {
		openers, err := logicLayer.GetOpenerLeaderboard()

		if err != nil {
			session.ChannelMessageSend(message.ChannelID, "se rompió todo, bancame un cachito uwu")
			return
		}

		var builder strings.Builder

		for index, opener := range openers {
			fmt.Fprintf(
				&builder,
				"%d - %s -> %d favorites\n",
				index+1,
				opener.Name,
				opener.Reactions,
			)

		}

		finalDescription := builder.String()

		embedImage := discordgo.MessageEmbedImage{
			URL: openers[0].ImageURL,
		}

		embedMessage := discordgo.MessageEmbed{
			Title:       "OPENER LEADERBOARD",
			Image:       &embedImage,
			Description: finalDescription,
			Type:        "image",
			Color:       rand.Intn(colorMaxValue),
		}

		session.ChannelMessageSendEmbed(message.ChannelID, &embedMessage)
		return
	}

	if matchesSetFavoriteOpenerPattern {

		newFavoriteOpener, err := logicLayer.SetFavoriteOpener(message)

		var str string

		if err != nil {
			str = fmt.Sprintf(err.Error())
		} else {

			str = fmt.Sprintf("Che bro ahora tu opener favorito es el ***%s***", newFavoriteOpener)
		}

		session.ChannelMessageSend(message.ChannelID, str)
		return
	}
}

func getUserNickname(message *discordgo.MessageCreate) string {
	if message.Member != nil {
		return message.Member.Nick
	}

	return message.Author.Username
}
