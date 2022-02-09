package security

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"strings"
	"time"

	"tripleoak/auth-api/services"

	"github.com/google/uuid"
	"github.com/lestrrat/go-jwx/jwa"
	"github.com/lestrrat/go-jwx/jwt"
)

func GenerateJWT(email string) (string, error) {
	token := jwt.New()

	// Payload
	token.Set("iss", "TripleOak")
	token.Set("sub", email)
	token.Set("aud", "https://api.tripleoak.pt/")
	token.Set("exp", time.Now().Add(15*time.Minute).Unix())
	token.Set("iat", time.Now().Unix())
	guid := uuid.New().String()
	token.Set("jti", guid)
	// Claims
	token.Set("admin", strings.HasSuffix(email, "@tripleoak.pt"))

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return "", err
	}

	keyString := base64.StdEncoding.EncodeToString(x509.MarshalPKCS1PrivateKey(key))
	collection := services.MongoClient.Client().Database("tripleoak").Collection("users")

	_, err = collection.UpdateOne(context.Background(), map[string]string{"email": email}, map[string]interface{}{"$set": map[string]interface{}{"privKey": keyString}})
	if err != nil {
		return "", err
	}

	tokenString, err := token.Sign(jwa.RS512, key)
	if err != nil {
		return "", err
	}

	return string(tokenString), nil
}

func VerifyJWT(tokenString string) (*jwt.Token, error) {
	token, err := jwt.ParseString(tokenString)
	if err != nil {
		return nil, err
	}

	collection := services.MongoClient.Client().Database("tripleoak").Collection("users")
	subject, _ := token.Get("sub")
	var result map[string]interface{}
	err = collection.FindOne(context.Background(), map[string]string{"email": subject.(string)}).Decode(&result)
	if err != nil {
		return nil, err
	}

	privKeyBytes, err := base64.StdEncoding.DecodeString(result["privKey"].(string))
	if err != nil {
		return nil, err
	}

	privKey, err := x509.ParsePKCS1PrivateKey(privKeyBytes)
	if err != nil {
		return nil, err
	}

	token, err = jwt.ParseString(tokenString, jwt.WithVerify(jwa.RS512, &privKey.PublicKey))
	if err != nil {
		return nil, err
	}

	return token, nil
}
