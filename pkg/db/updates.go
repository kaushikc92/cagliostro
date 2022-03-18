package db

import (
	"context"
	"time"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UpdatePosition struct {
	Fen string `bson:"fen,omitempty"`
	Depth int `bson:"depth,omitempty"`
}

func GetUpdatePosition(fen string) (*UpdatePosition, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(DatabaseUrl))
	if err != nil {
		return &UpdatePosition{}, err
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		return &UpdatePosition{}, err
	}
	defer client.Disconnect(ctx)

	collection := client.Database("cagliostro").Collection("updates")
	var result UpdatePosition
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

func UpsertUpdatePosition(upos UpdatePosition) error {
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

	collection := client.Database("cagliostro").Collection("updates")
	opts := options.Update().SetUpsert(true)
	filter := bson.M{"fen": upos.Fen}
	update := bson.D{{"$set", upos}}
	_, err = collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return err
	} else {
		return nil
	}
}
