package main

import (
	"flag"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/Haupc/attendance/core"
	"github.com/Haupc/attendance/utils"
	"github.com/joho/godotenv"
	"github.com/yanzay/tbot/v2"
	"k8s.io/klog/v2"
)

func init() {
	klog.InitFlags(nil)
	// By default klog writes to stderr. Setting logtostderr to false makes klog
	// write to a log file.
	flag.Set("logtostderr", "false")
	flag.Set("log_file", "attendance.log")
	flag.Set("alsologtostderr", "true")
	flag.Set("v", "3")
	flag.Parse()
	klog.Flush()
}
func main() {
	klog.V(3).Infoln("starting server...")
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	BOT_TOKEN := os.Getenv("TELEGRAM_TOKEN")

	rand.Seed(time.Now().UnixNano())
	go core.StoreUser()
	go func() {
		for {
			for _, user := range core.UserMap {
				go func(u *core.User) {
					if u.CanAttend() {
						u.Attendance()
						if run, ok := u.Result[time.Now().Day()]; !ok || !run {
							klog.V(3).Infoln("Next job will run on next minute")
						}
					}
				}(user)
			}
			time.Sleep(time.Minute)
		}
	}()

	// bot, err := tbot.NewServer(os.Getenv("TELEGRAM_TOKEN"))
	bot := tbot.New(BOT_TOKEN)
	core.SetClient(bot.Client())
	bot.HandleMessage("/test", core.TestHandler)
	bot.HandleMessage(utils.GetCommandPattern(core.ADD_COMMAND), core.AddUserHandler)
	bot.HandleMessage(core.CHECK_COMMAND, core.CheckHandler)
	bot.HandleMessage(core.HELP_COMMAND, core.UsageHandler)
	bot.HandleMessage(core.START_COMMAND, core.UsageHandler)
	log.Fatal(bot.Start())
}
