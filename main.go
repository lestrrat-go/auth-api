package main

import (
	"log"
	"net/http"

	"tripleoak/auth-api/rest"
	"tripleoak/auth-api/services"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	services.MongoClient.Init()
	defer services.MongoClient.Close()

	http.Handle("/", rest.Router())

	log.Println("[+] Auth API running on port 8080")
	http.ListenAndServe(":8080", nil)
}
