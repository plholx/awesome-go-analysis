package main

import "time"

// AwesomeGoInfo awesome-go项目仓库信息表
//
// Repo 是否是一个项目
// Category 是否是一类目
// Name 类目或github仓库亦或某官网名称
// Description 描述
// Homepage 官网主页地址
// CategoryHtmlId 源README.md文件中类别的锚点id
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
type AwesomeGoInfo struct {
	ID                   uint      `db:"id"`
	ParentId             uint      `db:"parent_id"`
	Repo                 bool      `db:"repo"`
	Category             bool      `db:"category"`
	Name                 string    `db:"name"`
	Description          string    `db:"description"`
	Homepage             string    `db:"homepage"`
	CategoryHtmlId       string    `db:"category_html_id"`
	RepoName             string    `db:"repo_name"`
	RepoFullName         string    `db:"repo_full_name"`
	RepoOwner            string    `db:"repo_owner"`
	RepoHtmlURL          string    `db:"repo_html_url"`
	RepoDescription      string    `db:"repo_description"`
	RepoCreatedAt        *time.Time `db:"repo_created_at"`
	RepoPushedAt         *time.Time `db:"repo_pushed_at"`
	RepoHomepage         string    `db:"repo_homepage"`
	RepoSize             uint      `db:"repo_size"`
	RepoForksCount       uint      `db:"repo_forks_count"`
	RepoStargazersCount  uint      `db:"repo_stargazers_count"`
	RepoSubscribersCount uint      `db:"repo_subscribers_count"`
	RepoOpenIssuesCount  uint      `db:"repo_open_issues_count"`
	RepoLicenseName      string    `db:"repo_license_name"`
	RepoLicenseSpdxId    string    `db:"repo_license_spdx_id"`
	RepoLicenseURL       string    `db:"repo_license_url"`
	CreatedAt            *time.Time `db:"created_at"`
	UpdatedAt            *time.Time `db:"updated_at"`

	Depth            uint
	Spaces           string
	TitleMarks       string
	RepoCreatedAtStr string
	RepoPushedAtStr  string
	WithReposTable   bool
	TimeSince        string
}
