package main

import (
	"fmt"
	"time"
)

type User struct {
	username string `json:"username"`
}

var m = map[string]User{
	"1": User{
		username: "1",
	},
	"2": User{
		username: "2",
	},
	"3": User{
		username: "3",
	},
}

func main() {
	for _, v := range m {
		v.username = "k"
	}
	i := 0
	for i < 10 {

		go func() {
			fmt.Println(i)
		}()
		i++
	}

	time.Sleep(time.Second)

}
