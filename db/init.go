package db

import (
	"context"
	. "github.com/sxyazi/bendan/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
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
		client, err = mongo.Connect(options.Client().ApplyURI(Config("db_uri")))
		if err != nil {
			log.Println("Database initialization failed:", err)
			log.Println("Database is disabled! Some database features may cause bot errors when called!")
			return
		}

		db = client.Database(Config("db_name"))
	})
	return db
}

func Indexes() {

	if Db() == nil {
		return
	}

	// Indexes for replied
	Db().Collection("replied").Indexes().DropAll(ctx)
	Db().Collection("replied").Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "chatId", Value: 1}, {Key: "id", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	})

	// Indexes for pinned
	Db().Collection("pinned").Indexes().DropAll(ctx)
	Db().Collection("pinned").Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "chatId", Value: 1}, {Key: "id", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	})

	// Indexes for forwarded
	Db().Collection("forwarded").Indexes().DropAll(ctx)
	Db().Collection("forwarded").Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "chatId", Value: 1}, {Key: "id", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	})
}
