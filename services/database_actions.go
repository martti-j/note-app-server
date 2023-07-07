package services

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type user struct {
	Username string `bson:"username"`
	Password string `bson:"password"`
}

func getDotEnvVariable(key string) string {

	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

var uri string = getDotEnvVariable("MONGODB_URI")

// Use the SetServerAPIOptions() method to set the Stable API version to 1
var serverAPI *options.ServerAPIOptions = options.ServerAPI(options.ServerAPIVersion1)
var opts *options.ClientOptions = options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

// Create a new client and connect to the server
var client, err = mongo.Connect(context.TODO(), opts)

func CheckConnectionToDB() {
	if err != nil {
		panic(err)
	}

	if err := client.Database("noteApp").RunCommand(context.TODO(), bson.D{{Key: "ping", Value: 1}}).Err(); err != nil {
		panic(err)
	}
	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")
}

func GetUsernames(c *gin.Context) {
	users := []*user{}

	db := client.Database("noteApp")
	coll := db.Collection("users")

	filter := bson.D{{}}

	cursor, err := coll.Find(context.TODO(), filter)

	if err != nil {
		panic(err)
	}

	for cursor.Next(context.TODO()) {
		var result user
		if err := cursor.Decode(&result); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%+v\n", result)
		users = append(users, &result)
	}
	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println(users)

	c.IndentedJSON(http.StatusOK, users)
}
