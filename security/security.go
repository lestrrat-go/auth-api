package security

import (
	"context"
	"net/http"

	"tripleoak/auth-api/services"
)

func Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("X-API-Key")
		if tokenString == "" {
			tokenString = r.URL.Query().Get("api_key")
			if tokenString == "" {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Missing API key"))
				return
			}
		}

		collection := services.MongoClient.Client().Database("financial-data").Collection("api-keys")
		var result map[string]map[string]interface{}
		err := collection.FindOne(context.Background(), map[string]interface{}{"_id": tokenString}).Decode(&result)
		if err != nil || result[tokenString]["disabled"] == true {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Invalid API key"))
			return
		}

		next.ServeHTTP(w, r)
	})
}
