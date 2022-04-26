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
var db *mongo.Database

func init() {
	var err error
	client, err = mongo.Connect(ctx, options.Client().ApplyURI(Config("db_uri")))
	if err != nil {
		log.Fatal(err)
	}

	db = client.Database(Config("db_name"))
}
