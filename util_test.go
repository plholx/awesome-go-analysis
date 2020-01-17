package main_test

import (
	"testing"
	"time"

	aga "awesome-go-analysis"
)

func TestDelay(t *testing.T) {
	timeX := time.Now().Add(time.Second * 5)
	t.Log("当前时间[阻塞前]：" + time.Now().Format("2006-01-02 15:04:05"))
	aga.Delay(timeX.Hour(), timeX.Minute(), timeX.Second(), 0, false)
	t.Log("当前时间[阻塞后]：" + time.Now().Format("2006-01-02 15:04:05"))
}

func TestTimeSince(t *testing.T) {
	timex := time.Now().Add(time.Second * 2 * (-1))
	thenStr := aga.TimeSince(&timex)
	t.Log(thenStr)
}

func TestGitPush(t *testing.T) {
	aga.GitPush("测试自动推送功能")
}
