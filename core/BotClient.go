package core

import "github.com/yanzay/tbot/v2"

var botClient *tbot.Client

func GetClient() *tbot.Client {
	return botClient
}

func SetClient(client *tbot.Client) {
	botClient = client
}
