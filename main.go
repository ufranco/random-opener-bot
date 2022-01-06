package main

import (
	"context"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/ufranco/random-opener-bot/bot"
	"github.com/ufranco/random-opener-bot/config"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	context := context.Background()

	config.ReadConfig()

	mongoClient, err := config.ConnectToDB(context)

	if err != nil {
		log.Fatal(err.Error())
		os.Exit(1)
	}

	repository := bot.NewRepo((mongoClient.Database("openerBot").Collection("openers")))

	session, err := discordgo.New("Bot " + config.BotToken)

	if err != nil {
		log.Fatal(err.Error())
		os.Exit(1)
	}

	bot.Start(context, session, repository)

	sessionChannel := make(chan os.Signal, 1)
	signal.Notify(sessionChannel, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sessionChannel

	session.Close()
	defer mongoClient.Disconnect(context)

}
