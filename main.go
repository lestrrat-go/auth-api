package main

import (
	"log"
	"net/http"

	"tripleoak/auth-api/rest"
	"tripleoak/auth-api/services"

	"github.com/joho/godotenv"
)

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE, PUT")
		w.Header().Set("Access-Control-Allow-Headers", "append,delete,entries,foreach,get,has,keys,set,values,Authorization")
		next.ServeHTTP(w, r)
	})
}

func main() {
	godotenv.Load()
	services.MongoClient.Init()
	defer services.MongoClient.Close()

	http.Handle("/", CORS(rest.Router()))

	log.Println("[+] Auth API running on port 8080")
	http.ListenAndServe(":8080", nil)
}
