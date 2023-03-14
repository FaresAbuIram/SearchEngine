package database

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

var (
	SearchCollection *mongo.Collection
	Ctx              = context.TODO()
)

func Connect() {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://mongo:27017"))
	if err != nil {
		log.Fatal(err)
	}

	Ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(Ctx)
	if err != nil {
		log.Fatal(err)
	}

	SearchDatabase := client.Database("searchEngine")
	SearchCollection = SearchDatabase.Collection("resources")

	indexModel := mongo.IndexModel{
		Keys: bson.M{
			"tags":  1,
		},
		Options: options.Index().SetName("resourcesIndex"),
	}

	_, err = SearchCollection.Indexes().CreateOne(context.Background(), indexModel)
	if err != nil {
		log.Fatal(err)
	}
}
