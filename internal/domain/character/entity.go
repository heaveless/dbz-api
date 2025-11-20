package character

type CharacterEntity struct {
	Id          int64  `bson:"_id" json:"id"`
	Name        string `bson:"name" json:"name"`
	Ki          string `bson:"ki" json:"ki"`
	MaxKi       string `bson:"maxKi" json:"maxKi"`
	Race        string `bson:"race" json:"race"`
	Gender      string `bson:"gender" json:"gender"`
	Image       string `bson:"image" json:"image"`
	Affiliation string `bson:"affiliation" json:"affiliation"`
}
