package main_test

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os/exec"
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

func TestPing(t *testing.T) {
	out, err := aga.Ping("raw.githubusercontent.com")
	if err != nil {
		t.Log(err)
	}
	t.Log("返回内容：\n", out)
}

func TestTail(t *testing.T) {
	out, err := aga.Tail("out.log")
	if err != nil {
		t.Log(err)
	}
	t.Log("返回内容：\n", out)
}

func TestTailf(t *testing.T) {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	cmd := exec.CommandContext(ctx, "tail", "-f", "out.log")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return
	}
	cmd.Start()
	reader := bufio.NewReader(stdout)
	for {
		line, err2 := reader.ReadString('\n')
		if err2 != nil || io.EOF == err2 {
			fmt.Println("退出循环")
			break
		}
		fmt.Print(line)
		if line == "shutdown\n" {
			fmt.Println("准备退出。。。")
			cancelFunc()
		}
	}
	fmt.Println("结束tail")
	cmd.Wait()
}
