package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

var (
	BotToken string
	config   *configStruct
)

type configStruct struct {
	BotToken string `json:"BOT_TOKEN"`
	MongoURI string `json:"MONGO_URI"`
}

func ReadConfig() {

	BotToken = os.Getenv("BOT_TOKEN")

	if BotToken == "" {

		file, err := ioutil.ReadFile("./config.json")

		if err != nil {
			log.Fatal(err.Error())
			os.Exit(1)
		}

		err = json.Unmarshal(file, &config)

		if err != nil {
			log.Fatal(err.Error())
			os.Exit(1)
		}

		BotToken = config.BotToken

	}

	log.Printf("Env variables loaded successfully!")
}
