package mongodb

import (
	"collyproject/config"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

var (
	client *mongo.Database
)

func init() {
	// Set client options
	clientOptions := options.Client().ApplyURI(config.MongoConf["connect"])


	// Connect to MongoDB
	if clt, err := mongo.Connect(context.TODO(), clientOptions); err != nil {
		log.Fatalln(err)
	} else {
		if err := clt.Ping(context.TODO(), nil); err != nil {
			log.Fatalln(err)
		}
		client = clt.Database(config.MongoConf["database"])

	}
	//collection := client.Collection("users")

	fmt.Println("Connected to MongoDB!")
}

type Collection struct {
	collection *mongo.Collection
	database *mongo.Database
}

func GetConnection() *mongo.Database {
	return client
}