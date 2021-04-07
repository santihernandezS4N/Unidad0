package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"os"
)

var ctx = context.TODO()
var uri = os.Getenv("DB_CONN")

func apiStatus(w http.ResponseWriter, r *http.Request) {
	clientOptions := options.Client().ApplyURI(uri) // Connect to //MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Panic(err)
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Panic(err)
	} else {
		fmt.Fprint(w, "Connected to MongoDB!")
	}
}
func getHash(pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	return string(hash)
}
func createUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var user User
	clientOptions := options.Client().ApplyURI(uri)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Panic(err)
	}

	collection := client.Database("Test").Collection("user")

	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Panic(err)
	}
	user.Password = getHash([]byte(user.Password))
	insertResult, err := collection.InsertOne(ctx, user)
	if err != nil {
		log.Panic(err)
		fmt.Fprintf(w, "User with this ID already exist")
	} else {
		json.NewEncoder(w).Encode(insertResult.InsertedID)
	}

}

func getUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	clientOptions := options.Client().ApplyURI(uri)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Panic(err)
	}

	collection := client.Database("Test").Collection("user")

	var users []*User
	usr, err := collection.Find(ctx, bson.D{{}})
	if err != nil {
		log.Panic(err)
		fmt.Fprint(w, "No data found")
	}

	for usr.Next(ctx) {
		var s User
		err := usr.Decode(&s)
		if err != nil {
			log.Panic(err)
		}
		users = append(users, &s)
	}
	if err := usr.Err(); err != nil {
		log.Panic(err)
	} else {
		json.NewEncoder(w).Encode(users)
	}
}

func getUserById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	clientOptions := options.Client().ApplyURI(uri)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Panic(err)
	}

	collection := client.Database("Test").Collection("user")

	var user User
	vars := mux.Vars(r)
	thisId := vars["id"]
	filter := bson.M{"_id": thisId}
	err = collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		log.Panic(err)
		fmt.Fprintf(w, "error retrieving user userid : %s", thisId)
	} else {
		json.NewEncoder(w).Encode(user)
	}
}

func getUserByName(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	clientOptions := options.Client().ApplyURI(uri)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Panic(err)
	}

	collection := client.Database("Test").Collection("user")

	vars := mux.Vars(r)
	thisName := vars["name"]
	filter := bson.M{"name": thisName}

	var users []*User
	usr, err := collection.Find(ctx, filter)
	if err != nil {
		log.Panic(err)
		fmt.Fprint(w, "No data found")
	}

	for usr.Next(ctx) {
		var s User
		err := usr.Decode(&s)
		if err != nil {
			log.Panic(err)
		}
		users = append(users, &s)
	}
	if err := usr.Err(); err != nil {
		log.Panic(err)
	} else {
		json.NewEncoder(w).Encode(users)
	}
}

func deleteUserById(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Panic(err)
	}

	vars := mux.Vars(r)
	thisId := vars["id"]
	filter := bson.M{"_id": thisId}
	collection := client.Database("Test").Collection("user")
	res, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		log.Fatal(err)
	} else {
		if res.DeletedCount == 0 {
			fmt.Fprint(w, "No data found")
		} else {
			fmt.Fprintf(w, "Succesfully deleted user with id %s", thisId)
		}
	}

}

func updateNamebyID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Panic(err)
	}
	vars := mux.Vars(r)
	thisId := vars["id"]
	thisName := vars["name"]
	filter := bson.M{"_id": thisId}
	collection := client.Database("Test").Collection("user")

	after := options.After // for returning updated document
	returnOpt := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
	}
	update := bson.D{{"$set", bson.D{{"name", thisName}}}}

	updateResult := collection.FindOneAndUpdate(context.TODO(), filter, update, &returnOpt)
	var result primitive.M
	updateResult.Decode(&result)
	json.NewEncoder(w).Encode(result)
}

