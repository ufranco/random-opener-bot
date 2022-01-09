package bot

type Account struct {
	ID             string `bson:"_id"`
	FavoriteOpener string `bson:"favorite_opener"`
}

type AccountRepository interface {
	Register(account Account) error
	FindById(userId string) (Account, error)
	UpdateFavoriteOpener(userId string, newFavoriteOpener string) error
}
