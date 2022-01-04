package core

import (
	"github.com/yanzay/tbot/v2"
	"k8s.io/klog/v2"
)

func TestHandler(m *tbot.Message) {
	klog.V(3).Infoln("test function")
	botClient.SendMessage(m.Chat.ID, "test pass!!")
}
