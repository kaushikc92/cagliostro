package db

import (
	"context"
	"time"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Puzzle struct {
	Fen string `bson:"fen"`
	Moves string `bson:"moves"`
	Rating int `bson:"rating"`
}

func GetPuzzle(minRating int, maxRating int) (*Puzzle, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(DatabaseUrl))
	if err != nil {
		return &Puzzle{}, err
	}
	ctx, _ := context.WithTimeout(context.Background(), 1000000*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		return &Puzzle{}, err
	}
	defer client.Disconnect(ctx)

	collection := client.Database("cagliostro").Collection("puzzles")
	matchStage := bson.D{
		{ "$match", bson.M{
			"rating": bson.M{
				"$gte": minRating,
				"$lt": maxRating,
			},
		}},
	}
	sampleStage := bson.D{{ "$sample", bson.M{"size": 1} }}

	puzzleCursor, err := collection.Aggregate(ctx, mongo.Pipeline{matchStage, sampleStage})
	if err != nil {
		return &Puzzle{}, err
	}
	var puzzle Puzzle

	if puzzleCursor.Next(ctx) {
		err := puzzleCursor.Decode(&puzzle)
		if err != nil {
			return &puzzle, err
		}
	}
	return &puzzle, nil
}
