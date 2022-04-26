package db

import (
	"github.com/sxyazi/bendan/types"
	"log"
	"time"
)

func AddLog(tag, content string) (interface{}, error) {
	_log := &types.Log{
		Tag:       tag,
		Content:   content,
		CreatedAt: time.Now(),
	}

	one, err := db.Collection("logs").InsertOne(ctx, _log)
	if err != nil {
		log.Println("AddLog error:", err)
		return 0, err
	}

	return one.InsertedID, nil
}
