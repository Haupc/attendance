package core

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"

	"github.com/Haupc/attendance/utils"
	"k8s.io/klog/v2"
)

var (
	UserMap  map[string]*User
	UserChan chan User
)

func init() {
	UserMap = map[string]*User{}
	UserChan = make(chan User)
}

type User struct {
	Username              string `json:"username"`
	Password              string `json:"password"`
	ChatId                string
	Result                map[int]bool
	nextTimeAttendance    int
	IsFixedTimeAttendance bool
}

func (u *User) doAttendance(client *http.Client, token string) {
	data := url.Values{}
	data.Set("_csrf-backend", token)
	data.Add("Attendance[type]", utils.GetAttendanceType())
	data.Add("Attendance[date_input]", time.Now().Format(utils.TIME_FORMAT))
	// data.Add("Attendance[date_input]", "2021-06-09")

	data.Add("Attendance[note]", "Làm việc online")
	req, err := http.NewRequest("POST", utils.ATTENDANCE_PAGE, strings.NewReader(data.Encode()))
	if err != nil {
		klog.V(3).Infof("Attendance setting err: %v", err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	response, err := client.Do(req)
	if err != nil {
		klog.V(3).Infof("Attendance request err: %v", err)
	}
	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)

	alerts := utils.AlertRe.FindAllString(string(body), -1)
	if len(alerts) > 0 && !strings.Contains(alerts[0], "bạn đã nộp xác nhận") {
		klog.V(3).Infof("Attendance err: %v\n", alerts)
		u.Result[time.Now().Day()] = false
		return
	}

	logString := fmt.Sprintf("Attendance success!!!!!!\n\n")
	botClient.SendMessage(u.ChatId, logString)

	u.Result[time.Now().Day()] = true
	//TODO: check fixed time
	u.nextTimeAttendance = utils.RandomTime()
	u.resetNextDay()
	logString = fmt.Sprintf("Next job will run on 7:%d %s\n\n", u.nextTimeAttendance, time.Now().Add(24*time.Hour).Format(utils.TIME_FORMAT))
	klog.V(3).Infoln(logString)
	botClient.SendMessage(u.ChatId, logString)
}

func (u *User) Attendance() {
	if time.Now().Weekday() == time.Sunday {
		klog.V(3).Infof("\nIt's Sunday, time for relax\n\n")
		u.Result[time.Now().Day()] = true
		u.resetNextDay()
		return
	}
	klog.V(3).Infof("Starting Attendance %s...\n", time.Now().Format(utils.TIME_FORMAT))
	cookieJar, _ := cookiejar.New(nil)

	client := &http.Client{
		Jar: cookieJar,
	}
	token, _ := utils.Login(client, u.Username, u.Password)
	if token == "" {
		klog.V(3).Infof("Cannel do attendance\n")
		return
	}
	u.doAttendance(client, token)
}

func (u *User) resetNextDay() {
	u.Result[time.Now().Add(24*time.Hour).Day()] = false
}

func (u *User) CanAttend() bool {
	run, ok := u.Result[time.Now().Day()]
	return (!ok || !run && (u.nextTimeAttendance <= time.Now().Minute() && time.Now().Hour() == 7 || time.Now().Hour() > 7))
}
