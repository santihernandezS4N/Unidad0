package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)


func HandleRequest() {
	myRouter := mux.NewRouter().StrictSlash(true)
	//test
	myRouter.HandleFunc("/", apiStatus).Methods("GET")

	//basic crud
	myRouter.HandleFunc("/", createUser).Methods("POST")

	myRouter.HandleFunc("/all", getUsers).Methods("GET")

	myRouter.HandleFunc("/id/{id}", getUserById).Methods("GET")

	myRouter.HandleFunc("/name/{name}", getUserByName).Methods("GET")

	myRouter.HandleFunc("/id/{id}", deleteUserById).Methods("DELETE")

	myRouter.HandleFunc("/id/{id}/name/{name}", updateNamebyID).Methods("PATCH")

	fmt.Println("Port 8080 is listening")
	log.Fatal(http.ListenAndServe(":8080", myRouter))

}

