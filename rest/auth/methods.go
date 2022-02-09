package auth

import (
	"context"
	"errors"
	"tripleoak/auth-api/security"
	"tripleoak/auth-api/services"

	"golang.org/x/crypto/bcrypt"
)

func Login(user string, password string) (string, error) {
	collection := services.MongoClient.Client().Database("tripleoak").Collection("users")
	var result map[string]interface{}
	err := collection.FindOne(context.Background(), map[string]interface{}{"email": user}).Decode(&result)
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(result["password"].(string)), []byte(password))
	if err != nil {
		return "", err
	}

	jwt, err := security.GenerateJWT(user)
	if err != nil {
		return "", err
	}

	return jwt, nil
}

func Logout(user string) error {
	collection := services.MongoClient.Client().Database("tripleoak").Collection("users")
	_, err := collection.UpdateOne(context.Background(), map[string]string{"email": user}, map[string]interface{}{"$unset": map[string]interface{}{"privKey": ""}})
	return err
}

func Signup(user string, password string) error {
	collection := services.MongoClient.Client().Database("tripleoak").Collection("users")

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	var result map[string]interface{}
	err = collection.FindOne(context.Background(), map[string]interface{}{"email": user}).Decode(&result)
	if err == nil {
		return errors.New("user already exists")
	}

	_, err = collection.InsertOne(context.Background(), map[string]string{"email": user, "password": string(bytes)})
	if err != nil {
		return err
	}

	return nil
}
