package database

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

var (
	SearchCollection *mongo.Collection
	Ctx              = context.TODO()
	ErrorLogger      *log.Logger
	InfoLogger       *log.Logger
	DebugLogger      *log.Logger
)

func Connect() {

	ErrorLogger = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	InfoLogger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)
	DebugLogger = log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)

	InfoLogger.Println("database,", "connection.go,", "Connect() Func")

	// create client for our database
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://mongo:27017"))
	if err != nil {
		ErrorLogger.Println("database doesn't exist: ", err)
		return
	}

	// connect to our database within 10 seconds
	Ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(Ctx)
	if err != nil {
		ErrorLogger.Println("timeout while connecting database: ", err)
		return
	}

	// create search engine database
	SearchDatabase := client.Database("searchEngine")
	//create resources collection
	SearchCollection = SearchDatabase.Collection("resources")

	// create an index in our database that sort the tags
	// this will make the search faster
	indexModel := mongo.IndexModel{
		Keys: bson.M{
			"tags": 1,
		},
		Options: options.Index().SetName("resourcesIndex"),
	}

	_, err = SearchCollection.Indexes().CreateOne(context.Background(), indexModel)
	if err != nil {
		ErrorLogger.Println("cannot create an index in adatabase: ", err)
		return
	}

	InfoLogger.Println("Successfully setup the database")
}
