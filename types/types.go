package types

import "time"

// RepliedMessage only represents the replied message from this bot,
// under certain cases (e.g. reply to channel post that needs to be forwarded).
type RepliedMessage struct {
	ID        int   `bson:"id"`
	ChatID    int64 `bson:"chatId"`
	RepliedTo int   `bson:"repliedTo"`
}

type PinnedMessage struct {
	ID        int       `bson:"id"`
	ChatID    int64     `bson:"chatId"`
	CreatedAt time.Time `bson:"createdAt"`
}

type ForwardedMessage struct {
	ID int `bson:"id"`

	// Due to Telegram's limitation of only can retrieve up to 2 of the messages in a context,
	// (i.e. current message, and the message it's replied to), so we need store the text for interactivity.
	// But that data is still under your control, cause it's only stored in your database where the `Forward` feature depends on self-deploying.
	Text string `bson:"text"`

	ChatID     int64  `bson:"chatId"`
	PhotoID    string `bson:"photoId"`
	PhotoGroup string `bson:"photoGroup"`

	TweetID      string    `bson:"tweetId"`
	TweetURL     string    `bson:"tweetUrl"`
	TweetPhotoID string    `bson:"tweetPhotoId"`
	TootID       string    `bson:"tootId"`
	TootURL      string    `bson:"tootUrl"`
	TootPhotoID  string    `bson:"tootPhotoId"`
	CreatedAt    time.Time `bson:"createdAt"`
}
