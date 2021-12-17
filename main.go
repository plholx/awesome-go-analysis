package main

import (
	"flag"
)


func main(){
	var path, token string
	flag.StringVar(&path, "p", ".", "生成的README.md文件路径")
	flag.StringVar(&token, "t", "xxx", "GitHub API access_token")
	flag.Parse()

	InitDB()
	// 启动README.md文件解析任务
	StartReadmeParseJob(path, token)
	signal := make(chan int)
	<-signal
}