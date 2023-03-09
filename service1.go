package main

import (
	"encoding/json"
	"math/rand"
	"net/http"
)
func generateSaltHandler(w http.ResponseWriter, r *http.Request) {
	salt := generateRandomString(12)
	response := map[string]string{
		"salt": salt,
	}
	json.NewEncoder(w).Encode(response)
}

func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}

func getSalt() (string, error) {
	res, err := http.Post("http://localhost:8000/generate-salt", "application/json", nil)
	if err != nil {
		return "", err
	}

	var response struct {
		Salt string `json:"salt"`
	}

	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return "", err
	}

	return response.Salt, nil
}
