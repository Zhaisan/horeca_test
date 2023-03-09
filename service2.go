package main

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
)

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !isValidEmail(user.Email) {
		http.Error(w, "Invalid email", http.StatusBadRequest)
		return
	}

	if !isValidPassword(user.Password) {
		http.Error(w, "Invalid password", http.StatusBadRequest)
		return
	}

	db, err := getDB()
	if err != nil {
		http.Error(w, "Error connecting to database", http.StatusInternalServerError)
		return
	}
	defer db.Client().Disconnect(context.Background())

	collection := db.Collection("users")
	filter := bson.M{"email": user.Email}
	count, err := collection.CountDocuments(context.Background(), filter)
	if err != nil {
		http.Error(w, "Error querying database", http.StatusInternalServerError)
		return
	}
	if count > 0 {
		http.Error(w, "Email already exists", http.StatusBadRequest)
		return
	}

	salt, err := getSalt()
	if err != nil {
		http.Error(w, "Error generating salt", http.StatusInternalServerError)
		return
	}

	passwordHash := hashPassword(salt, user.Password)

	newUser := User{
		Email:    user.Email,
		Salt:     salt,
		Password: passwordHash,
	}
	_, err = collection.InsertOne(context.Background(), newUser)
	if err != nil {
		http.Error(w, "Error inserting", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}


func getUserHandler(w http.ResponseWriter, r *http.Request) {
	email := chi.URLParam(r, "email")

	db, err := getDB()
	if err != nil {
		http.Error(w, "Error connecting to database", http.StatusInternalServerError)
		return
	}
	defer db.Client().Disconnect(context.Background())

	collection := db.Collection("users")
	filter := bson.M{"email": email}
	var user User
	err = collection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(user)
}

func hashPassword(salt string, password string) string {
	saltedPassword := salt + password
	return fmt.Sprintf("%x", md5.Sum([]byte(saltedPassword)))
}