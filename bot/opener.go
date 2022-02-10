package bot

type Opener struct {
	Name        string `bson:"_id"`
	Description string `bson:"description"`
	Reactions   int32  `bson:"reactions"`
	ImageURL    string `bson:"image_url"`
}

type OpenerRepository interface {
	FindById(id string) (Opener, error)
	GetRandomOpener() (Opener, error)
	UpdateReactionBy(openerName string, quantity int) error
	GetLeaderboard() ([]Opener, error)
}
