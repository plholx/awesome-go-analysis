package main

import (
	"time"

	"github.com/jinzhu/gorm"
)

// AwesomeGoInfo awesome-go项目仓库最新信息表
//
// Repo 是否是一个项目
// Category 是否是一类目
// Name 类目或github仓库亦或某官网名称
// Description 描述
// Homepage 官网主页地址
// CategoryHtmlId 源README.md文件中类别的锚点id
type AwesomeGoInfo struct {
	gorm.Model
	RepoInfo

	ParentId       uint   `gorm:"default:0"`
	Repo           bool   `gorm:"default:false"`
	Category       bool   `gorm:"default:false"`
	Name           string `gorm:"type:varchar(50)"`
	Description    string `gorm:"type:varchar(1000)"`
	Homepage       string `gorm:"type:varchar(500)"`
	CategoryHtmlId string `gorm:"type:varchar(100)"`

	Depth            uint   `gorm:"-"`
	Spaces           string `gorm:"-"`
	TitleMarks       string `gorm:"-"`
	RepoCreatedAtStr string `gorm:"-"`
	RepoPushedAtStr  string `gorm:"-"`
	WithReposTable   bool   `gorm:"-"`
	TimeSince        string `gorm:"-"`
}

// GithubRepoRecord awesome-go项目仓库信息记录表
//
// AgiId awesome_go_info表主键
type GithubRepoRecord struct {
	gorm.Model
	AgiId uint `gorm:"default:0"`
	RepoInfo
}

// RepoInfo 仓库信息
//
// RepoName 项目仓库名称
// RepoFullName 项目仓库完整名称
// RepoOwner 项目作者
// RepoHtmlURL 项目在github上的地址
// RepoDescription 项目描述
// RepoCreatedAt 项目创建时间
// RepoPushedAt 项目最后推送时间
// RepoHomepage 项目官网主页
// RepoForksCount 项目fork数
// RepoStargazersCount 项目star数
// RepoSubscribersCount 项目watch数
type RepoInfo struct {
	RepoName             string `gorm:"type:varchar(50)"`
	RepoFullName         string `gorm:"type:varchar(100)"`
	RepoOwner            string `gorm:"type:varchar(50)"`
	RepoHtmlURL          string `gorm:"type:varchar(500)"`
	RepoDescription      string `gorm:"type:varchar(1000)"`
	RepoCreatedAt        *time.Time
	RepoPushedAt         *time.Time
	RepoHomepage         string `gorm:"type:varchar(500)"`
	RepoSize             uint   `gorm:"default:0"`
	RepoForksCount       uint   `gorm:"default:0"`
	RepoStargazersCount  uint   `gorm:"default:0"`
	RepoSubscribersCount uint   `gorm:"default:0"`
	RepoOpenIssuesCount  uint   `gorm:"default:0"`
	RepoLicenseName      string `gorm:"type:varchar(100)"`
	RepoLicenseSpdxId    string `gorm:"type:varchar(50)"`
	RepoLicenseURL       string `gorm:"type:varchar(500)"`
}
