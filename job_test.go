package main_test

import (
	"fmt"
	analysis "github.com/plholx/awesome-go-analysis"
	"log"
	"regexp"
	"testing"
)

func TestReg(t *testing.T){
	line := "*Libraries for manipulating audio.*"
	// 分类的描述正则，如：*Software written in Go.*
	categoryDescReg := regexp.MustCompile(`(\s*)\*([^⬆]*)\*$`)
	if categoryDescReg.MatchString(line) {
		groups := categoryDescReg.FindStringSubmatch(line)
		fmt.Printf("正则匹配结果: %s\n", groups[2])
	}

	line = "* [flac](https://github.com/mewkiz/flac) - Native Go FLAC encoder/decoder with support for FLAC streams."
	// 判断是否包含链接正则，如：
	containsLinkReg        := regexp.MustCompile(`^\s*\* \[.*\]\(.*\)`)
	if containsLinkReg.MatchString(line) {
		groups := containsLinkReg.FindStringSubmatch(line)
		fmt.Printf("正则匹配结果: %s\n", groups[0])
	}

	line = ""
	// github 链接正则
	gitHubURLReg           := regexp.MustCompile(`https://github.com/(.+?)/([a-zA-Z0-9_\-\.]+)(.*)$`)
	if gitHubURLReg.MatchString(line) {
		groups := gitHubURLReg.FindStringSubmatch(line)
		fmt.Printf("正则匹配结果: %s\n", groups[0])
	}
}

func TestUpdateAwesomeGoRepo(t *testing.T){
	readmeFilePath, err := analysis.UpdateAwesomeGoRepo()
	if err != nil{
		log.Printf("更新仓库失败: %+v", err)
		return
	}
	log.Println(readmeFilePath)
}

func TestParseReadmeFile(t *testing.T){
	analysis.InitDB()
	err := analysis.ParseReadmeFile("awesome-go/README.md", "xxx")
	if err != nil{
		log.Printf("解析文件异常: %+v", err)
	}
}

func TestConvertToCategoryHtmlId(t *testing.T){
	fmt.Println(analysis.ConvertToCategoryHtmlId("Libraries and123 tools for manipulating XML."))
}

func TestGetCategoryInfoByName(t *testing.T){
	analysis.InitDB()
	id := analysis.GetCategoryInfoByName("Job Scheduler")
	fmt.Printf("类目id: %d", id)
}

func TestGenerateMd(t *testing.T){
	analysis.InitDB()
	analysis.GenerateMd("temp");
}

func TestGetAGI(t *testing.T){
	analysis.InitDB()
	id, err := analysis.GetAGI("mewkiz/flac", true, false)
	if err != nil{
		log.Printf("%+v\n", err)
		return
	}
	log.Printf("id: %d\n", id)
}