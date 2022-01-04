package core

import (
	"fmt"

	"github.com/yanzay/tbot/v2"
)

var CHECK_COMMAND = "/check"

func CheckHandler(t *tbot.Message) {
	user, ok := UserMap[t.Chat.ID]
	if !ok {
		botClient.SendMessage(t.Chat.ID, "You need to add a user first!!")
		UsageHandler(t)
		return
	}

	botClient.SendMessage(t.Chat.ID, fmt.Sprintf("Today attendance status: %v", !user.CanAttend()))
	botClient.SendMessage(t.Chat.ID, fmt.Sprintf("next attendance will be on 7h%v", user.nextTimeAttendance))
}
