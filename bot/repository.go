package bot

import (
	"errors"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var repoErr error = errors.New("Unable to handle Repo request")

type repo struct {
	collection *mongo.Collection
}

func NewRepo(collection *mongo.Collection) Repository {
	return &repo{collection}
}

func (repo *repo) GetRandomOpener() (Opener, error) {
	openers := make([]Opener, 1)

	pipeline := []bson.M{{"$match": bson.D{}}, {"$sample": bson.M{"size": 1}}}

	loadedCursor, err := repo.collection.Aggregate(
		applicationContext,
		pipeline,
	)

	if err != nil {
		log.Fatal(err.Error())
		return Opener{}, repoErr
	}

	if err = loadedCursor.All(applicationContext, &openers); err != nil {
		return Opener{}, err
	}

	return openers[0], nil
}

func (repo *repo) GetLeaderboard() ([]Opener, error) {
	return make([]Opener, 0), nil

}
