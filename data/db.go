package data

import (
	"context"
	"fmt"
	"log"

	"github.com/Smart-Pot/pkg"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Collection of comments
var collection *mongo.Collection

// DatabaseConnection :
func DatabaseConnection() {

	clientOptions := options.Client().ApplyURI(pkg.Config.Database.Addr)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to Database!")
	collection = client.Database(pkg.Config.Database.DBName).Collection("users")
}
