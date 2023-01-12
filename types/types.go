package types

import "time"

// RepliedMessage only represents the replied message from this bot,
// under certain cases (e.g. reply to channel post that needs to be forwarded).
type RepliedMessage struct {
	Id        int   `bson:"id"`
	ChatId    int64 `bson:"chatId"`
	RepliedTo int   `bson:"repliedTo"`
}

type PinnedMessage struct {
	Id        int       `bson:"id"`
	ChatId    int64     `bson:"chatId"`
	CreatedAt time.Time `bson:"createdAt"`
}

type ForwardedMessage struct {
	Id int `bson:"id"`

	// Due to Telegram's limitation of only can retrieve up to 2 of the messages in a context,
	// (i.e. current message, and the message it's replied to), so we need store the text for interactivity.
	// But that data is still under your control, cause it's only stored in your database where the `Forward` feature depends on self-deploying.
	Text string `bson:"text"`

	ChatId     int64  `bson:"chatId"`
	PhotoId    string `bson:"photoId"`
	PhotoGroup string `bson:"photoGroup"`

	TweetId      string    `bson:"tweetId"`
	TweetUrl     string    `bson:"tweetUrl"`
	TweetPhotoId string    `bson:"tweetPhotoId"`
	TootId       string    `bson:"tootId"`
	TootUrl      string    `bson:"tootUrl"`
	TootPhotoId  string    `bson:"tootPhotoId"`
	CreatedAt    time.Time `bson:"createdAt"`
}
