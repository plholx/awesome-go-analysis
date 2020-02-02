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
	// 启动README.md文件解析任务
	StartMDParseJob()

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello Gopher!"))
	})
	// 测试github的raw.githubusercontent.com是否可ping通
	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		output, err := Ping("raw.githubusercontent.com")
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		w.Write([]byte(output))
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
				filePath, err := DownloadREADMEFile()
				if err != nil {
					log.Println("文件下载失败", err)
				} else {
					//解析awesome-go中的README.md文件,并存入数据库中
					err := ParseREADMEFile(viper.GetString("token"), filePath)
					if err != nil {
						log.Println("文件解析失败", err)
					} else {
						//生成README.md
						GenerateMd()
						//自动push,测试push
						GitPush("自动推送新生成的README.md文件")
					}
				}

			} else {
				log.Printf("%s文件已解析，不再重复解析", filePath)
			}

			log.Println("业务执行结束：", time.Now().Format("2006-01-02 15:04:05"))
			log.Println("等待下次执行......")

			Delay(0, 0, 0, 0, true)
		}
	}()
}
