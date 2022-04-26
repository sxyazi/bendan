package types

import "time"

type Config struct {
	Token string `json:"token"`
	Db    string `json:"db"`
}

type Log struct {
	Tag       string    `bson:"tag"`
	Content   string    `bson:"content"`
	CreatedAt time.Time `bson:"createdAt"`
}

type PinnedMessage struct {
	Id         int       `bson:"_id"`
	ChatId     int64     `bson:"chatId"`
	FromUserId int64     `bson:"fromUserId"`
	CreatedAt  time.Time `bson:"createdAt"`
}
