package core

import (
	"fmt"
	"strings"

	"github.com/yanzay/tbot/v2"
	"k8s.io/klog/v2"
)

var ADD_COMMAND = "/add"

func AddUserHandler(t *tbot.Message) {
	klog.V(3).Infoln(t.Text)
	cmd := strings.Split(strings.TrimSpace(t.Text), " ")
	if len(cmd) != 3 {
		botClient.SendMessage(t.Chat.ID, "wrong format\nusage: /add <email> <password>")
		return
	}
	email := strings.TrimSpace(cmd[1])
	password := strings.TrimSpace(cmd[2])
	name := AddUserToChan(email, password, t.Chat.ID)
	if name == "" {
		botClient.SendMessage(t.Chat.ID, "wrong email/password")
		return
	}
	botClient.SendMessage(t.Chat.ID, fmt.Sprintf("welcome %s\nuser add success!!", name))
}
