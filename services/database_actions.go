package services

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
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

var db = client.Database("noteApp")
var coll = db.Collection("users")

func CheckConnectionToDB() {
	if err != nil {
		panic(err)
	}

	if err := client.Database("noteApp").RunCommand(context.TODO(), bson.D{{Key: "ping", Value: 1}}).Err(); err != nil {
		panic(err)
	}
	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")
}

func GetUsersDB() ([]*User, error) {
	users := []*User{}
	filter := bson.D{{}}
	cursor, err := coll.Find(context.TODO(), filter)

	if err != nil {
		return users, errors.New("didn't find names")
	}

	for cursor.Next(context.TODO()) {
		var result User
		if err := cursor.Decode(&result); err != nil {
			return users, errors.New("could't decode the package")
		}
		fmt.Printf("%+v\n", result)
		users = append(users, &result)
	}
	if err := cursor.Err(); err != nil {
		return users, errors.New("error")
	}

	return users, nil
}

func GetUserByUsernameDB(username string) (User, error) {
	filter := bson.D{{Key: "username", Value: username}}
	var result User
	cursor := coll.FindOne(context.TODO(), filter).Decode(&result)
	fmt.Printf("%+v\n", cursor)

	if cursor == mongo.ErrNoDocuments {
		return User{}, errors.New("didn't find username")
	}

	return User{Username: result.Username, Password: result.Password}, nil
}

func AddUserDB(username string, password string) error {
	newUser := User{
		Username: username, Password: password,
	}

	_, err := coll.InsertOne(context.TODO(), newUser)

	if err != nil {
		return errors.New("failed to add user")
	}

	return nil
}
