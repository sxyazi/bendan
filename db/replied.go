package db

import (
	"github.com/sxyazi/bendan/types"
	"go.mongodb.org/mongo-driver/bson"
	"log"
)

func AddReplied(rm *types.RepliedMessage) error {
	_, err := Db().Collection("replied").InsertOne(ctx, rm)
	if err != nil {
		log.Println("AddReplied error:", err)
		return err
	}

	return nil
}

func GetReplied(id int, chatId int64) (*types.RepliedMessage, error) {
	var record *types.RepliedMessage
	err := Db().
		Collection("replied").
		FindOne(ctx, bson.M{"id": id, "chatId": chatId}).
		Decode(&record)

	if err != nil {
		log.Println("GetReplied error:", err)
		return nil, err
	}
	return record, nil
}
