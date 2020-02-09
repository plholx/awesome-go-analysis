package main_test

import (
	"io/ioutil"
	"testing"

	"github.com/spf13/viper"

	aga "awesome-go-analysis"
)

func TestDownloadREADMEFile(t *testing.T) {
	filePath, err := aga.DownloadREADMEFile()
	if err != nil {
		t.Error(err)
	}
	t.Log(filePath)
}

func TestGenerateMd(t *testing.T) {
	aga.InitDB()
	aga.GenerateMd()
}

func TestGitHubApiReq(t *testing.T) {
	viper.SetDefault("token", "abc")
	resp, err := aga.GitHubApiReq("GET", "https://api.github.com/rate_limit", nil)
	if err != nil {
		t.Log(err)
		return
	}
	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Log(err)
		return
	}
	t.Log("响应数据：", string(bytes))
}
