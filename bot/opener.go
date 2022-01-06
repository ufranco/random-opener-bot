package bot

type Opener struct {
	ID          string `json:"_id" bson:"_id"`
	Name        string `json:"name" bson:"name"`
	Description string `json:"description" bson:"description"`
	Reactions   int32  `json:"reactions" bson:"reactions"`
	ImageURL    string `json:"image_url" bson:"image_url"`
}

type Repository interface {
	GetLeaderboard() ([]Opener, error)
	GetRandomOpener() (Opener, error)
}
