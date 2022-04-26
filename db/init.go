package db

import (
	"context"
	. "github.com/sxyazi/bendan/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

var ctx = context.TODO()
var client *mongo.Client

func init() {
	var err error
	client, err = mongo.Connect(ctx, options.Client().ApplyURI(Config("db")))
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
}
