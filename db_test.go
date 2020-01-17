package main_test

import (
	"testing"

	"github.com/jinzhu/gorm"

	aga "awesome-go-analysis"
)

func TestUpdateAGIDescription(t *testing.T) {
	aga.InitDB()
	aga.UpdateAGIDescription("abc", 1)
}

func TestUpdateAGIGithubInfo(t *testing.T) {
	aga.InitDB()
	awesomeGoInfo := aga.AwesomeGoInfo{
		Model: gorm.Model{
			ID: 1,
		},
		RepoInfo: aga.RepoInfo{RepoHtmlURL: "http://www.github.com"},
	}
	aga.UpdateAGIGithubInfo(&awesomeGoInfo)
}

func TestGetAGI(t *testing.T) {
	aga.InitDB()
	agi, _ := aga.GetAGI("a", false, false)
	t.Log("$$$:", agi)
}

func TestGetAGIByCategoryHtmlId(t *testing.T) {
	aga.InitDB()
	agi, _ := aga.GetAGIByCategoryHtmlId("1")
	t.Log("$$$:", agi)
}

func TestModifyAGIParentIdByCategoryHtmlId(t *testing.T) {
	aga.InitDB()
	aga.ModifyAGIParentIdByCategoryHtmlId(1, "7")
}

func TestSaveAGI(t *testing.T) {
	aga.InitDB()
	agi := new(aga.AwesomeGoInfo)
	t.Log("插入条数：", aga.SaveAGI(agi))
	t.Log("返回主键：", agi.ID)
}

func TestSaveGRR(t *testing.T) {
	aga.InitDB()
	grr := new(aga.GithubRepoRecord)
	t.Log("插入条数：", aga.SaveGRR(grr))
	t.Log("返回主键：", grr.ID)
}

func TestGetAGITree(t *testing.T) {
	aga.InitDB()
	agis, err := aga.GetAGITree(false)
	if err != nil {
		t.Log(err)
	}
	t.Log("返回数据量", len(agis))
	if len(agis) > 0 {
		t.Log("名字：", agis[0].Name)
	}
}

func TestGetAgiReposByParentId(t *testing.T) {
	aga.InitDB()
	agis, err := aga.GetAgiReposByParentId(100000)
	if err != nil {
		t.Log(err)
	}
	if agis == nil {
		t.Log("agis为空：", agis)
	}
	t.Log("返回数据量", len(agis))
	if len(agis) > 0 {
		t.Log("名字：", agis[0].Name)
	}
}

func TestSlice(t *testing.T) {
	var agis []*aga.AwesomeGoInfo
	t.Log(agis, agis == nil)
	for i := 0; i < 2; i++ {
		tempAgi := aga.AwesomeGoInfo{
			Name: string(i),
		}
		agis = append(agis, &tempAgi)
	}
	t.Log("名字1：", agis[0].Name)
	t.Log("名字2：", agis[1].Name)

	var ints []int
	t.Log("数字切片", ints)
	intm := make([]int, 0)
	t.Log("make数字切片", intm)
	var inta [1]int
	t.Log("数字数组", inta)
}
