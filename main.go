package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/spf13/viper"

	_ "awesome-go-analysis/conf"
)

func main() {
	InitDB()
	agis, err1 := GetAGITree(false)
	if err1 != nil {
		log.Println(err1)
	}
	log.Println("返回数据量", len(agis))
	if len(agis) > 0 {
		log.Println("名字：", agis[0].Name)
	}
	// 启动README.md文件解析任务
	StartMDParseJob()

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello Gopher!"))
	})

	// 创建server对象
	server := &http.Server{
		Addr:    ":4000",
		Handler: mux,
	}

	// 创建系统信号接收器
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	go func() {
		<-quit

		if err := server.Shutdown(context.Background()); err != nil {
			log.Fatal("Shutdown server:", err)
		}
	}()

	log.Println("Starting Http server...")
	err := server.ListenAndServe()
	if err != nil {
		if err == http.ErrServerClosed {
			log.Print("Serer closed under request")
		} else {
			log.Fatal("Server closed unexpected")
		}
	}
}

// StartMDParseJob 启动README.md文件解析任务
func StartMDParseJob() {
	go func() {
		for {
			log.Println("业务开始执行：", time.Now().Format("2006-01-02 15:04:05"))

			if filePath, has := CurREADMEFile(); !has {
				//获取avelino/awesome-go项目中最新的README.md文件
				filePath, _ := DownloadREADMEFile()
				//解析awesome-go中的README.md文件,并存入数据库中
				ParseREADMEFile(viper.GetString("token"), filePath)
				//生成README.md
				GenerateMd()
				//自动push,测试push
				GitPush("自动推送新生成的README.md文件")
			} else {
				log.Printf("%s文件已解析，不再重复解析", filePath)
			}

			log.Println("业务执行结束：", time.Now().Format("2006-01-02 15:04:05"))
			log.Println("等待下次执行......")

			Delay(0, 0, 0, 0, true)
		}
	}()
}
