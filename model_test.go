package main_test

import (
	"testing"

	"github.com/jinzhu/gorm"

	aga "awesome-go-analysis"
)

func TestModelStruct(t *testing.T) {
	db, err := gorm.Open("sqlite3", "./data_test.db")
	if err != nil {
		t.Error(err)
	}
	defer db.Close()

	if !db.HasTable(&aga.AwesomeGoInfo{}) {
		db.CreateTable(&aga.AwesomeGoInfo{})
	}
	if !db.HasTable(&aga.GithubRepoRecord{}) {
		db.CreateTable(&aga.GithubRepoRecord{})
	}

	githubRepoRecord := aga.GithubRepoRecord{RepoInfo: aga.RepoInfo{RepoName: "测试"}}
	db.Create(&githubRepoRecord)
	t.Log("主键: ", githubRepoRecord.ID)

	awesomeGoInfo := aga.AwesomeGoInfo{}
	db.Create(&awesomeGoInfo)
	db.Model(&awesomeGoInfo).Update("description", "xxx")

	db.LogMode(true)
	tempAwesomeGoInfo := new(aga.AwesomeGoInfo)
	db.Where("description = ?", "x").Where("name = ?", "b").Find(tempAwesomeGoInfo)
	t.Log("**********************", tempAwesomeGoInfo)
	db.Where("id = ?", 2).Find(tempAwesomeGoInfo)
	t.Log("**********************", tempAwesomeGoInfo)
}
