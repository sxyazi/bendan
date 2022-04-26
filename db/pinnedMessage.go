package db

import (
	"fmt"
	"github.com/sxyazi/bendan/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

func AddPinnedMessage(data *types.PinnedMessage) (any, error) {
	data.CreatedAt = time.Now()
	one, err := client.Database("pinned_messages").Collection(fmt.Sprintf("chat_%d", data.ChatId)).InsertOne(ctx, data)
	if err != nil {
		log.Println("AddPinnedMessage error:", err)
		return 0, err
	}

	return one.InsertedID, nil
}

func RemovePinnedMessage(id int, chatId int64) (int64, error) {
	one, err := client.Database("pinned_messages").Collection(fmt.Sprintf("chat_%d", chatId)).DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		log.Println("RemovePinnedMessage error:", err)
		return 0, err
	}

	return one.DeletedCount, nil
}

func GetPinnedMessages(chatId int64) ([]*types.PinnedMessage, error) {
	findOptions := options.Find()
	findOptions.SetSort(bson.M{"$natural": -1})

	fmt.Println(fmt.Sprintf("chat_%d", chatId))
	cur, err := client.Database("pinned_messages").Collection(fmt.Sprintf("chat_%d", chatId)).Find(ctx, bson.M{}, findOptions)
	if err != nil {
		log.Println("GetPinnedMessages error:", err)
		return nil, err
	}

	var data []*types.PinnedMessage
	for cur.Next(ctx) {
		var result types.PinnedMessage
		if err := cur.Decode(&result); err != nil {
			return nil, err
		}

		data = append(data, &result)
	}

	return data, nil
}
