package core

import (
	"net/http"
	"net/http/cookiejar"
	"strings"

	"github.com/Haupc/attendance/utils"
)

func AddUserToChan(username, password, chatId string) string {
	username = strings.TrimSpace(username)
	password = strings.TrimSpace(password)
	cookieJar, _ := cookiejar.New(nil)
	token, name := utils.Login(&http.Client{
		Jar: cookieJar,
	}, username, password)
	if token != "" {
		UserChan <- User{
			Username:              username,
			Password:              password,
			ChatId:                chatId,
			Result:                map[int]bool{},
			IsFixedTimeAttendance: false,
			nextTimeAttendance:    20,
		}
	}
	return name
}

func StoreUser() {
	for {
		u := <-UserChan
		UserMap[u.ChatId] = &u
	}
}
