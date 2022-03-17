package db

const DatabaseUrl string = "mongodb://mongo-0.mongo,mongo-1.mongo:27017"

type Position struct {
	Fen string `bson:"fen,omitempty"`
	BestMove string `bson:"bestMove,omitempty"`
	Depth string `bson:"depth,omitempty"`
	MyMove string `bson:"myMove,omitempty"`
}
