package config

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectToDB(context context.Context) (*mongo.Client, error) {
	mongoURI := os.Getenv("MONGO_URI")

	if mongoURI == "" {

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
		mongoURI = config.MongoURI
	}

	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(context, clientOptions)

	if err != nil {
		log.Fatal(err.Error())
	}

	if err = client.Ping(context, nil); err != nil {
		log.Fatal(err.Error())
	}

	return client, nil
}
