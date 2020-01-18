package main

import (
	"database/sql"
	"log"
	"sync"

	"github.com/spf13/viper"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var db *gorm.DB
var dbErr error
var dbonce sync.Once

func InitDB() {
	dbonce.Do(func() {
		db, dbErr = gorm.Open("sqlite3", viper.GetString("dsurl"))
		if dbErr != nil {
			log.Fatal(dbErr)
		}
		db.LogMode(true)
		// 创建表
		if !db.HasTable(&AwesomeGoInfo{}) {
			db.CreateTable(&AwesomeGoInfo{})
			log.Println("创建表 awesome_go_infos")
		}
		if !db.HasTable(&GithubRepoRecord{}) {
			db.CreateTable(&GithubRepoRecord{})
			log.Println("创建表 github_repo_records")
		}
		log.Println("数据库初始化成功")
	})
}

// UpdateAGIDescription 更新awesome_go_info表description
func UpdateAGIDescription(description string, id uint) {
	awesomeGoInfo := AwesomeGoInfo{
		Description: description,
		Model: gorm.Model{
			ID: id,
		},
	}
	db.Model(awesomeGoInfo).Update("description", awesomeGoInfo.Description)
}

// UpdateAGIGithubInfo 更新仓库github仓库相关信息
func UpdateAGIGithubInfo(agi *AwesomeGoInfo) {
	db.Model(agi).Updates(map[string]interface{}{
		"repo_html_url":          agi.RepoHtmlURL,
		"repo_description":       agi.RepoDescription,
		"repo_pushed_at":         agi.RepoPushedAt,
		"repo_homepage":          agi.RepoHomepage,
		"repo_size":              agi.RepoSize,
		"repo_forks_count":       agi.RepoForksCount,
		"repo_stargazers_count":  agi.RepoStargazersCount,
		"repo_subscribers_count": agi.RepoSubscribersCount,
		"repo_open_issues_count": agi.RepoOpenIssuesCount,
		"repo_license_name":      agi.RepoLicenseName,
		"repo_license_spdx_id":   agi.RepoLicenseSpdxId,
		"repo_license_url":       agi.RepoLicenseURL,
		"name":                   agi.Name,
		"description":            agi.Description,
		"homepage":               agi.Homepage,
		"parent_id":              agi.ParentId,
	})
}

func errNR(s *gorm.DB) error {
	if s != nil && s.RowsAffected < 1 {
		return sql.ErrNoRows
	}
	return nil
}

// GetAGI 获取awesome_go_info表数据
func GetAGI(name string, repo bool, category bool) (agi *AwesomeGoInfo, err error) {
	agi = new(AwesomeGoInfo)
	err = errNR(db.Where("name = ? and repo = ? and category = ?", name, repo, category).First(agi))
	return
}

// GetAGIByCategoryHtmlId 根据CategoryHtmlId获取awesome_go_info表分类数据(类别以锚点id为唯一标识)
func GetAGIByCategoryHtmlId(categoryHtmlId string) (agi *AwesomeGoInfo, err error) {
	agi = new(AwesomeGoInfo)
	err = errNR(db.Where("category_html_id = ? and category = true", categoryHtmlId).First(agi))
	return
}

// ModifyAGIParentIdByCategoryHtmlId 更新awesome_go_info表parent_id
func ModifyAGIParentIdByCategoryHtmlId(parentId uint, categoryHtmlId string) int64 {
	agi := new(AwesomeGoInfo)
	return db.Model(agi).Where("category_html_id = ? and category = true", categoryHtmlId).Update("parent_id", parentId).RowsAffected
}

// SaveAGI awesome_go_info表插入数据
func SaveAGI(agi *AwesomeGoInfo) int64 {
	return db.Create(agi).RowsAffected
}

// SaveGRR github_repo_record表插入数据
func SaveGRR(grr *GithubRepoRecord) int64 {
	return db.Create(grr).RowsAffected
}

// GetAGITree 获取awesome_go_info信息树
func GetAGITree(all bool) (agis []*AwesomeGoInfo, err error) {
	sqlStr := `WITH RECURSIVE subordinates AS (
		SELECT id, parent_id, repo_name, repo_full_name, repo_owner, repo_html_url, repo_description, repo_created_at, repo_pushed_at, repo_homepage, repo_size, repo_forks_count, repo_stargazers_count, repo_subscribers_count, repo_open_issues_count, repo_license_name, repo_license_spdx_id, repo_license_url, repo, category, name, description, homepage, category_html_id, 1 depth FROM awesome_go_infos WHERE parent_id = 0 and category=true and deleted_at is null
		UNION ALL
		SELECT  r.id, r.parent_id, r.repo_name, r.repo_full_name, r.repo_owner, r.repo_html_url, r.repo_description, r.repo_created_at, r.repo_pushed_at, r.repo_homepage, r.repo_size, r.repo_forks_count, r.repo_stargazers_count, r.repo_subscribers_count, r.repo_open_issues_count, r.repo_license_name, r.repo_license_spdx_id, r.repo_license_url, r.repo, r.category, r.name, r.description, r.homepage, r.category_html_id, s.depth+1 depth FROM awesome_go_infos r
		inner JOIN subordinates s ON r.parent_id = s.id and r.category=true and r.deleted_at is null order by depth desc
		) SELECT id, parent_id, repo_name, repo_full_name, repo_owner, repo_html_url, repo_description, repo_created_at, repo_pushed_at, repo_homepage, repo_size, repo_forks_count, repo_stargazers_count, repo_subscribers_count, repo_open_issues_count, repo_license_name, repo_license_spdx_id, repo_license_url, repo, category, name, description, homepage, category_html_id, depth FROM subordinates`
	rows, err := db.Raw(sqlStr).Rows()
	defer rows.Close()
	if err != nil {
		log.Println(err)
		return
	}

	for rows.Next() {
		tmpAGI := new(AwesomeGoInfo)
		db.ScanRows(rows, tmpAGI)
		tmpAGI.Spaces = getSpace(int64(tmpAGI.Depth) - 1)
		tmpAGI.TitleMarks = getTitleMarks(int64(tmpAGI.Depth))
		agis = append(agis, tmpAGI)
		// 需要同时查询仓库信息
		if all {
			repos, err := GetAgiReposByParentId(tmpAGI.ID)
			if err == nil && repos != nil && len(repos) > 0 {
				tmpAGI.WithReposTable = true
				agis = append(agis, repos...)
			}
		}
	}
	return
}

// GetAgisByParentId 根据父id查询AwesomeGoInfo表中的github仓库信息
func GetAgiReposByParentId(parentId uint) (agis []*AwesomeGoInfo, err error) {
	sqlStr := `select
	id, parent_id, repo_name, repo_full_name, repo_owner, repo_html_url, repo_description, repo_created_at, repo_pushed_at, repo_homepage, repo_size, repo_forks_count, repo_stargazers_count, repo_subscribers_count, repo_open_issues_count, repo_license_name, repo_license_spdx_id, repo_license_url, repo, category, name, description, homepage, category_html_id
	from awesome_go_infos
	where parent_id = ? and repo = true order by repo_stargazers_count desc`
	rows, err := db.Raw(sqlStr, parentId).Rows()
	defer rows.Close()
	if err != nil {
		log.Println(err)
		return
	}
	for rows.Next() {
		tmpAGI := new(AwesomeGoInfo)
		db.ScanRows(rows, tmpAGI)
		tmpAGI.RepoCreatedAtStr = tmpAGI.RepoCreatedAt.Format("2006-01-02")
		tmpAGI.RepoPushedAtStr = tmpAGI.RepoPushedAt.Format("2006-01-02 15:04:05")
		tmpAGI.TimeSince = TimeSince(tmpAGI.RepoPushedAt)
		agis = append(agis, tmpAGI)
	}
	return
}

func getSpace(count int64) (s string) {
	for i := int64(0); i < count; i++ {
		s += "    "
	}
	return s
}

func getTitleMarks(count int64) (s string) {
	for i := int64(0); i < count; i++ {
		s += "#"
	}
	return s
}
