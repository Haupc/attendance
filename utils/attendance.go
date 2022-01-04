package utils

import (
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strings"
	"time"

	"k8s.io/klog/v2"
)

var (
	Re      = regexp.MustCompile(`(?m)csrf-token.*content="(.*)">`)
	AlertRe = regexp.MustCompile(`(?m)icon fa fa-ban".*`)
	NameRe  = regexp.MustCompile(`(?m)"hidden-xs".*`)
)

const (
	LOGIN_PAGE      = "https://hr.ewings.vn/site/login"
	ATTENDANCE_PAGE = "https://hr.ewings.vn/attendance/create"
	TIME_FORMAT     = "2006-01-02"
	// cronPattern       = "TZ=Asia/Ho_Chi_Minh 0 8 * * *"
)

func GetAttendanceType() string {
	if time.Now().Weekday() == time.Saturday {
		return "1"
	}
	return "0"
}

func RandomTime() int {
	klog.V(3).Infof("Geting random time: 7h00 -> 7h45")
	return rand.Intn(45)
}

func getLoginPage(client *http.Client) string {
	response, err := client.Get(LOGIN_PAGE)
	if err != nil {
		klog.V(3).Infof("get login page err: %v", err)
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		klog.V(3).Infof("read body from login page err: %v", err)
	}
	token := Re.FindAllStringSubmatch(string(body), -1)[0][1]
	return token
}

func Login(client *http.Client, username, password string) (string, string) {
	token := getLoginPage(client)
	data := url.Values{}
	data.Set("LoginForm[username]", username)
	data.Add("LoginForm[password]", password)
	data.Add("_csrf-backend", token)
	req, err := http.NewRequest("POST", LOGIN_PAGE, strings.NewReader(data.Encode()))
	if err != nil {
		klog.V(3).Infof("login setting err: %v", err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	response, err := client.Do(req)
	if err != nil {
		klog.V(3).Infof("login request err: %v", err)
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		klog.V(3).Infof("read body login err: %v", err)
	}
	if strings.Contains(string(body), "/site/logout") {
		klog.V(3).Infoln("Login success!!!!!!!!")
		name := NameRe.FindAllString(string(body), -1)
		klog.V(3).Infoln("Human name: ", name)
		token := Re.FindAllStringSubmatch(string(body), -1)[0][1]
		return token, name[0]
	} else {
		klog.V(3).Infof("Login Fail!!!!!!!!!!!!")
		return "", ""
	}
}

func GetCommandPattern(c string) string {
	return c + ".*"
}

func GetName(username string, password string) string {
	cookieJar, _ := cookiejar.New(nil)
	_, name := Login(&http.Client{
		Jar: cookieJar,
	}, username, password)
	return name
}
