package db

import (
	"context"
	"time"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)



const DatabaseUrl string = "mongodb://mongo-0.mongo,mongo-1.mongo:27017"

type ErrPositionDoesntExist struct {
	Err error
}

func (err *ErrPositionDoesntExist) Error() string {
	return err.Err.Error()
}


type Position struct {
	Fen string `bson:"fen,omitempty"`
	BestMove string `bson:"bestMove,omitempty"`
	Depth int `bson:"depth,omitempty"`
	MyMove string `bson:"myMove,omitempty"`
}

func GetPosition(fen string) (*Position, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(DatabaseUrl))
	if err != nil {
		return &Position{}, err
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		return &Position{}, err
	}
	defer client.Disconnect(ctx)

	collection := client.Database("cagliostro").Collection("repertoire")
	var result Position
	err = collection.FindOne(ctx, bson.M{"fen": fen}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			doesntExistErr := &ErrPositionDoesntExist{
				Err: errors.New("This position does not exist in the database"),
			}
			return nil, doesntExistErr
		} else {
			return nil, err
		}
	} else {
		return &result, nil
	}
}

func UpsertPosition(pos Position) error {
	client, err := mongo.NewClient(options.Client().ApplyURI(DatabaseUrl))
	if err != nil {
		return err
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		return err
	}
	defer client.Disconnect(ctx)

	collection := client.Database("cagliostro").Collection("repertoire")
	opts := options.Update().SetUpsert(true)
	filter := bson.M{"fen": pos.Fen}
	update := bson.D{{"$set", pos}}
	_, err = collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return err
	} else {
		return nil
	}
}
