package types

import "time"

type Config struct {
	Token string `json:"token"`
	Db    string `json:"db"`
}

type PinnedMessage struct {
	Id        int       `bson:"_id"`
	ChatId    int64     `bson:"chatId"`
	CreatedAt time.Time `bson:"createdAt"`
}
