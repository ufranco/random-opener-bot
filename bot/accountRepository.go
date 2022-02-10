package bot

import (
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type accountRepo struct {
	collection *mongo.Collection
}

func NewAccountRepository(collection *mongo.Collection) AccountRepository {
	return &accountRepo{collection}
}

func (repo *accountRepo) Register(account Account) error {
	_, err := repo.collection.InsertOne(applicationContext, account)

	if err != nil {
		log.Fatal(err.Error())
		return errRepo
	}

	return nil
}

func (repo *accountRepo) FindById(userId string) (Account, error) {
	account := Account{}

	err := repo.collection.FindOne(applicationContext, bson.M{"_id": userId}).Decode(&account)

	if err != nil {
		return account, err
	}

	return account, nil
}

func (repo *accountRepo) UpdateFavoriteOpener(userId string, openerName string) error {
	result := repo.collection.FindOneAndUpdate(
		applicationContext,
		bson.M{"_id": userId},
		bson.M{"$set": bson.M{"favorite_opener": openerName}},
	)

	if result == nil {
		log.Fatal(errRepo)
		return errRepo
	}

	return nil
}
