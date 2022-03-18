package db

import (
	"context"
	"time"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type RepertoirePosition struct {
	Fen string `bson:"fen,omitempty"`
	BestMove string `bson:"bestMove,omitempty"`
	Depth int `bson:"depth,omitempty"`
	MyMove string `bson:"myMove,omitempty"`
}

func GetRepertoirePosition(fen string) (*RepertoirePosition, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(DatabaseUrl))
	if err != nil {
		return &RepertoirePosition{}, err
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		return &RepertoirePosition{}, err
	}
	defer client.Disconnect(ctx)

	collection := client.Database("cagliostro").Collection("repertoire")
	var result RepertoirePosition
	err = collection.FindOne(ctx, bson.M{"fen": fen}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			doesntExistErr := &ErrRecordDoesntExist{
				Err: errors.New("This record does not exist in the database"),
			}
			return nil, doesntExistErr
		} else {
			return nil, err
		}
	} else {
		return &result, nil
	}
}

func UpsertRepertoirePosition(rpos RepertoirePosition) error {
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
	filter := bson.M{"fen": rpos.Fen}
	update := bson.D{{"$set", rpos}}
	_, err = collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return err
	} else {
		return nil
	}
}
