package main

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"sync"
)

var db *sqlx.DB
var dbonce sync.Once

func InitDB() {
	dbonce.Do(func() {
		db = sqlx.MustOpen("sqlite3", "awesomego.db")
		// 创建表 awesome_go_info，先判断库中是否已存在该表
		tableCount := 0
		db.Get(&tableCount, `select count(*) from sqlite_master where type = 'table' and name = 'awesome_go_info'`)
		if tableCount == 0{
			log.Println("创建表: awesome_go_info")
			db.MustExec(awesome_go_info_schema)
		}
		log.Println("数据库初始化成功")
	})
}

const awesome_go_info_schema = `
create table awesome_go_info
(
    id                     integer
        primary key autoincrement,
    repo_name              varchar(50) default '',
    repo_full_name         varchar(100) default '',
    repo_owner             varchar(50) default '',
    repo_html_url          varchar(500) default '',
    repo_description       varchar(1000) default '',
    repo_created_at        datetime,
    repo_pushed_at         datetime,
    repo_homepage          varchar(500) default '',
    repo_size              integer default 0,
    repo_forks_count       integer default 0,
    repo_stargazers_count  integer default 0,
    repo_subscribers_count integer default 0,
    repo_open_issues_count integer default 0,
    repo_license_name      varchar(100) default '',
    repo_license_spdx_id   varchar(50) default '',
    repo_license_url       varchar(500) default '',
    parent_id              integer default 0,
    repo                   bool    default false,
    category               bool    default false,
    name                   varchar(50) default '',
    description            varchar(1000) default '',
    homepage               varchar(500) default '',
    category_html_id       varchar(100) default '',
    created_at             datetime,
    updated_at             datetime
);
`