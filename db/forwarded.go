package db

import (
	"errors"
	"github.com/sxyazi/bendan/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

func AddForwarded(fm *types.ForwardedMessage) error {
	fm.CreatedAt = time.Now()
	coll := Db().Collection("forwarded")

	if _, err := coll.InsertOne(ctx, fm); err == nil {
		return nil
	} else if mongo.IsDuplicateKeyError(err) {
		err = coll.FindOneAndReplace(ctx, bson.M{"id": fm.Id, "chatId": fm.ChatId}, fm).Err()
		if err != nil {
			log.Println("AddForwarded error:", err)
		}
		return err
	} else {
		log.Println("AddForwarded error:", err)
		return err
	}
}

func UpdateForwarded(id int, chatId int64, update any) error {
	_, err := Db().
		Collection("forwarded").
		UpdateMany(ctx, bson.M{"id": id, "chatId": chatId}, bson.M{"$set": update})
	return err
}

func UpdateForwardedByGroup(group string, chatId int64, update any) error {
	_, err := Db().
		Collection("forwarded").
		UpdateMany(ctx, bson.M{"photoGroup": group, "chatId": chatId}, bson.M{"$set": update})
	return err
}

func GetForwarded(id int, chatId int64) ([]*types.ForwardedMessage, error) {
	var record *types.ForwardedMessage
	err := Db().
		Collection("forwarded").
		FindOne(ctx, bson.M{"id": id, "chatId": chatId}).
		Decode(&record)

	if err != nil {
		log.Println("GetForwarded error:", err)
		return nil, err
	}
	return GetForwardedByGroup(record.PhotoGroup, chatId)
}

func GetForwardedByGroup(group string, chatId int64) ([]*types.ForwardedMessage, error) {
	opt := options.Find().SetSort(bson.M{"id": 1})
	cur, err := Db().
		Collection("forwarded").
		Find(ctx, bson.M{"photoGroup": group, "chatId": chatId}, opt)

	if err != nil {
		log.Println("GetForwardedByGroup error:", err)
		return nil, err
	}

	var records []*types.ForwardedMessage
	if err = cur.All(ctx, &records); err != nil {
		return nil, err
	}
	return records, nil
}

func GetForwardedGroupOne(group string, chatId int64) (*types.ForwardedMessage, error) {
	var record *types.ForwardedMessage
	err := Db().
		Collection("forwarded").
		FindOne(ctx, bson.M{"photoGroup": group, "chatId": chatId}).
		Decode(&record)

	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		log.Println("GetForwardedGroupOne error:", err)
		return nil, err
	}
	return record, nil
}
