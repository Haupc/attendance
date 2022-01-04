package core

import "github.com/yanzay/tbot/v2"

var (
	HELP_COMMAND  = "/help"
	START_COMMAND = "/start"
)

func UsageHandler(t *tbot.Message) {
	botClient.SendMessage(t.Chat.ID, "usage:\n/help - show this usage\n/check - check today attendance status\n/add <email> <password> - add user to attendance")
}
