package main_test

import (
	"TikTok/biz/dao"
	"fmt"
	"github.com/go-co-op/gocron"
	"testing"
	"time"
)

func TestAddLike(t *testing.T) {
	fmt.Println(t)
	dao.Init()
	timezone, _ := time.LoadLocation("Asia/Shanghai")
	s := gocron.NewScheduler(timezone)

	// 每秒执行一次
	s.Every(1).Seconds().Do(func() {
		go cron1()
	})

	// 每秒执行一次
	s.Every(1).Second().Do(func() {
		go cron2()
	})

	s.StartBlocking()

}
func cron1() {
	fmt.Println("cron1")
}

func cron2() {
	fmt.Println("cron2")
}
