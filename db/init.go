package db

import (
	"context"
	. "github.com/sxyazi/bendan/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"sync"
)

var ctx = context.TODO()
var client *mongo.Client
var db *mongo.Database
var connectOnce = sync.Once{}

func Db() *mongo.Database {
	connectOnce.Do(func() {
		var err error
		client, err = mongo.Connect(ctx, options.Client().ApplyURI(Config("db_uri")))
		if err != nil {
			log.Fatal(err)
		}

		db = client.Database(Config("db_name"))
	})
	return db
}

func Indexes() {
	// Indexes for replied
	Db().Collection("replied").Indexes().DropAll(ctx)
	Db().Collection("replied").Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys:    bson.D{{"chatId", 1}, {"id", 1}},
			Options: options.Index().SetUnique(true),
		},
	})

	// Indexes for pinned
	Db().Collection("pinned").Indexes().DropAll(ctx)
	Db().Collection("pinned").Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys:    bson.D{{"chatId", 1}, {"id", 1}},
			Options: options.Index().SetUnique(true),
		},
	})

	// Indexes for forwarded
	Db().Collection("forwarded").Indexes().DropAll(ctx)
	Db().Collection("forwarded").Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys:    bson.D{{"chatId", 1}, {"id", 1}},
			Options: options.Index().SetUnique(true),
		},
	})
}
