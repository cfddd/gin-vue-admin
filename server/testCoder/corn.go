package main

import (
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/utils/timer"
	"time"
)

func Timer() (err error) {

	//corn框架中的定时器
	t := timer.NewTimerTask()
	t.AddTaskByFunc("ClearDB", "*/1 * * * *", func() {
		fmt.Println("cfddfc")
	})
	//time.Sleep(time.Second * 100)
	return err
}

func main() {
	Timer()
	time.Sleep(time.Second * 100)
	return
}
