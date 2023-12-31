package services

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	Username string `bson:"username"`
	Password string `bson:"password"`
}

type AddNoteType struct {
	Tittle  string `bson:"tittle"`
	Content string `bson:"content"`
	Writer  string `bson:"writer"`
}

type Note struct {
	Tittle  string             `bson:"tittle"`
	Content string             `bson:"content"`
	Writer  string             `bson:"writer"`
	ID      primitive.ObjectID `bson:"_id"`
}

type NoteDeleteType struct {
	ID   primitive.ObjectID `bson:"_id"`
	User string             `bson:"username"`
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
var collUsers = db.Collection("users")
var collNotes = db.Collection("notes")

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
	cursor, err := collUsers.Find(context.TODO(), filter)

	if err != nil {
		return users, errors.New("didn't find names")
	}

	for cursor.Next(context.TODO()) {
		var result User
		if err := cursor.Decode(&result); err != nil {
			return users, errors.New("could't decode the package")
		}
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
	cursor := collUsers.FindOne(context.TODO(), filter).Decode(&result)
	fmt.Printf("%+v\n", cursor)

	if cursor == mongo.ErrNoDocuments {
		return User{}, errors.New("didn't find username")
	}

	return User{Username: result.Username, Password: result.Password}, nil
}

func AddUserDB(newUser User) error {

	allUsers, dataErr := GetUsersDB()

	if dataErr != nil {
		return errors.New("failed to add user")
	}

	for i := 0; i < len(allUsers); i++ {
		if allUsers[i].Username == newUser.Username {
			return errors.New("username already in use")
		}
	}

	_, err := collUsers.InsertOne(context.TODO(), newUser)

	if err != nil {
		return errors.New("failed to add user")
	}

	return nil
}

func DeleteUserDB(deleteUser User) error {
	_, err := collUsers.DeleteOne(context.TODO(), bson.M{"username": deleteUser.Username})
	if err != nil {
		return errors.New("couldn't delete user")
	}
	return nil
}

func LoginDB(loginUser User) error {
	foundUser, dataErr := GetUserByUsernameDB(loginUser.Username)

	if dataErr != nil {
		return errors.New("didn't find username from database")
	}

	if loginUser.Password == foundUser.Password {
		return nil
	}

	return errors.New("wrong password")
}

func GetNotesDB() ([]*Note, error) {
	notes := []*Note{}
	filter := bson.D{{}}
	cursor, err := collNotes.Find(context.TODO(), filter)

	if err != nil {
		return notes, errors.New("didn't find notes")
	}

	for cursor.Next(context.TODO()) {
		var result Note
		if err := cursor.Decode(&result); err != nil {
			return notes, errors.New("could't decode the package")
		}
		notes = append(notes, &result)
	}
	if err := cursor.Err(); err != nil {
		return notes, errors.New("error")
	}

	return notes, nil
}

func AddNoteDB(newNote AddNoteType) error {
	_, err := collNotes.InsertOne(context.TODO(), newNote)

	if err != nil {
		return errors.New("failed to add note")
	}

	return nil
}

func DeleteNoteDB(deleteNote NoteDeleteType) error {
	_, err := collNotes.DeleteOne(context.TODO(), bson.M{"_id": deleteNote.ID})
	if err != nil {
		return errors.New("couldn't delete note")
	}
	return nil
}

func GetNoteByIDDB(id primitive.ObjectID) (Note, error) {
	filter := bson.D{{Key: "_id", Value: id}}
	var result Note
	cursor := collNotes.FindOne(context.TODO(), filter).Decode(&result)
	fmt.Printf("%+v\n", cursor)

	if cursor == mongo.ErrNoDocuments {
		return Note{}, errors.New("didn't find note id")
	}

	return Note{Tittle: result.Tittle, Content: result.Content, Writer: result.Writer,
		ID: result.ID}, nil
}
