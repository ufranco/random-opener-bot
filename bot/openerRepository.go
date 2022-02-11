package bot

import (
	"errors"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var errRepo error = errors.New("unable to handle repo request")

type openerRepo struct {
	collection *mongo.Collection
}

func NewOpenerRepository(collection *mongo.Collection) OpenerRepository {
	return &openerRepo{collection}
}

func (repo *openerRepo) FindById(id string) (Opener, error) {
	opener := Opener{}

	cursor := repo.collection.FindOne(applicationContext, bson.M{"_id": id})

	err := cursor.Decode(&opener)

	return opener, err
}

func (repo *openerRepo) GetRandomOpener() (Opener, error) {

	pipeline := []bson.M{{"$match": bson.D{}}, {"$sample": bson.M{"size": 1}}}

	loadedCursor, err := repo.collection.Aggregate(
		applicationContext,
		pipeline,
	)

	if err != nil {
		log.Fatal(err.Error())
		return Opener{}, errRepo
	}

	defer loadedCursor.Close(applicationContext)

	var openers []Opener

	if err = loadedCursor.All(applicationContext, &openers); err != nil {
		return Opener{}, err
	}

	return openers[0], nil
}

func (repo *openerRepo) GetLeaderboard() ([]Opener, error) {
	sort := options.Find()

	sort.SetSort(bson.M{"reactions": -1})

	loadedCursor, err := repo.collection.Find(applicationContext, bson.D{}, sort)

	if err != nil {
		log.Fatal(err)
		return []Opener{}, errRepo
	}

	defer loadedCursor.Close(applicationContext)

	var openers []Opener

	if err = loadedCursor.All(applicationContext, &openers); err != nil {
		return []Opener{}, err
	}

	return openers, nil
}

func (repo *openerRepo) UpdateReactionBy(openerName string, quantity int) error {
	_, err := repo.collection.UpdateOne(
		applicationContext,
		bson.M{"_id": openerName},
		bson.M{"$inc": bson.M{"reactions": quantity}},
	)

	if err != nil {
		log.Fatal(errRepo)
		return errRepo
	}

	return errRepo
}
