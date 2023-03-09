package main

import (
	"context"
	"github.com/go-chi/chi"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
)


func getDB() (*mongo.Database, error) {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	db := client.Database("mydb")
	return db, nil
}


func main() {
	router1 := chi.NewRouter()
	router1.Post("/generate-salt", generateSaltHandler)

	router2 := chi.NewRouter()
	router2.Post("/create-user", createUserHandler)
	router2.Get("/get-user/{email}", getUserHandler)

	go func() {
		log.Fatal(http.ListenAndServe(":8000", router1))
	}()

	log.Fatal(http.ListenAndServe(":8001", router2))

}
