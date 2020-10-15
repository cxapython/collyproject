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

var Conf = config.Config
func init() {
	// Set client options
	clientOptions := options.Client().ApplyURI(Conf.GetString("MONGO.URL"))


	// Connect to MongoDB
	if clt, err := mongo.Connect(context.TODO(), clientOptions); err != nil {
		log.Fatalln(err)
	} else {
		if err := clt.Ping(context.TODO(), nil); err != nil {
			log.Fatalln(err)
		}
		client = clt.Database(Conf.GetString("MONGO.DB_NAME"))

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