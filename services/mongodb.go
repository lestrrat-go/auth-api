package services

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Client struct {
	client *mongo.Client
}

func (mongoCli *Client) Init() {
	mongo_user := os.Getenv("MONGODB_USER")
	mongo_pass := os.Getenv("MONGODB_PASS")
	mongo_host := os.Getenv("MONGODB_HOST")
	mongo_db := os.Getenv("MONGODB_DB")
	clientOptions := options.Client().
		ApplyURI("mongodb+srv://" + mongo_user + ":" + mongo_pass + "@" + mongo_host + "/" + mongo_db + "?retryWrites=true&w=majority")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var err error
	mongoCli.client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = mongoCli.client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func (mongoCli *Client) Close() {
	mongoCli.client.Disconnect(context.Background())
}

func (mongoCli *Client) Client() *mongo.Client {
	return mongoCli.client
}

var MongoClient Client = Client{}
