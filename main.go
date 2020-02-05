package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	_ "awesome-go-analysis/conf"
)

func main() {
	InitDB()
	// 启动README.md文件解析任务
	StartMDParseJob()
	// 使用gin
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello Gopher!")
	})

	// 路由组使用 gin.BasicAuth() 中间件
	// gin.Accounts 是 map[string]string 的一种快捷方式
	authorized := router.Group("/admin", gin.BasicAuth(gin.Accounts{
		"admin": viper.GetString("pwd"),
	}))
	// 测试github的raw.githubusercontent.com是否可ping通
	authorized.GET("/ping", func(c *gin.Context) {
		output, err := Ping("raw.githubusercontent.com")
		if err != nil {
			c.String(http.StatusOK, err.Error())
			return
		}
		c.String(http.StatusOK, output)
	})
	// 查看日志文件
	authorized.GET("/tail", func(c *gin.Context) {
		output, err := Tail("out.log")
		if err != nil {
			c.String(http.StatusOK, err.Error())
			return
		}
		c.String(http.StatusOK, output)
	})
	// 手动执行文件下载及解析任务
	authorized.GET("/mexec", func(c *gin.Context) {
		go func() {
			log.Println("手动执行业务开始：", time.Now().Format("2006-01-02 15:04:05"))
			ParseJob()
			log.Println("手动执行业务结束：", time.Now().Format("2006-01-02 15:04:05"))
		}()
		c.String(http.StatusOK, "任务手动调用成功，详情请查看日志 /tail")
	})
	// 创建server对象
	srv := &http.Server{
		Addr:    ":4000",
		Handler: router,
	}
	go func() {
		// 服务连接
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}

// StartMDParseJob 启动README.md文件解析任务
func StartMDParseJob() {
	go func() {
		for {
			log.Println("业务开始执行：", time.Now().Format("2006-01-02 15:04:05"))
			ParseJob()
			log.Println("业务执行结束：", time.Now().Format("2006-01-02 15:04:05"))
			log.Println("等待下次执行......")

			Delay(0, 0, 0, 0, true)
		}
	}()
}

func ParseJob() {
	if filePath, has := CurREADMEFile(); !has {
		// 获取avelino/awesome-go项目中最新的README.md文件
		filePath, err := DownloadREADMEFile()
		if err != nil {
			log.Println("文件下载失败", err)
		} else {
			// 解析awesome-go中的README.md文件,并存入数据库中
			err := ParseREADMEFile(viper.GetString("token"), filePath)
			if err != nil {
				log.Println("文件解析失败", err)
			} else {
				// 生成README.md
				GenerateMd()
				// 自动push,测试push
				GitPush("自动推送新生成的README.md文件")
			}
		}

	} else {
		log.Printf("%s文件已解析，不再重复解析", filePath)
	}
}
