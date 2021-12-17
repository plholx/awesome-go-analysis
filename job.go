package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"text/template"
	"time"
)

// StartReadmeParseJob 启动README.md文件解析任务
func StartReadmeParseJob(path, token string) {
	go func() {
		for {
			log.Println("业务开始执行：", time.Now().Format("2006-01-02 15:04:05"))
			ReadmeParseJob(path, token)
			log.Println("业务执行结束：", time.Now().Format("2006-01-02 15:04:05"))
			log.Println("等待下次执行......")

			Delay(0, 0, 0, 0, true)
		}
	}()
}

// ReadmeParseJob 解析README.md文件任务
func ReadmeParseJob(path, token string) {
	// 获取avelino/awesome-go项目中最新的README.md文件
	readmeFilePath, err := UpdateAwesomeGoRepo()
	if err != nil {
		log.Printf("更新https://github.com/avelino/awesome-go.git失败: %+v\n", err)
		return
	}
	fmt.Println(readmeFilePath)
	// 解析awesome-go中的README.md文件,并存入数据库中
	err = ParseReadmeFile(readmeFilePath, token)
	if err != nil {
		log.Println("解析README.md失败", err)
		return
	}
	// 生成README.md
	GenerateMd(path)
	// 自动push,测试push
	GitPush("自动推送新生成的README.md文件")
}

func GitPush(comment string){
	cmd := exec.Command(git, add, ".")
	output, err := cmd.CombinedOutput()
	fmt.Println(string(output))
	if err != nil {
		log.Printf(" git add .异常: %+v", err)
		return
	}

	cmd = exec.Command(git, commit, "-m", comment)
	output, err = cmd.CombinedOutput()
	fmt.Println(string(output))
	if err != nil {
		log.Printf(" git commit 异常: %+v", err)
		return
	}

	cmd = exec.Command(git, push)
	output, err = cmd.CombinedOutput()
	fmt.Println(string(output))
	if err != nil {
		log.Printf(" git push 异常: %+v", err)
		return
	}
}

type data struct {
	Categorys []*AwesomeGoInfo
	GoRepos   []*AwesomeGoInfo
	GenTime   string
}

const (
	README_TMPL      = "tmpl.md"
	README           = "README.md"
)

//GenerateMd 生成README.md文件
func GenerateMd(path string) (err error) {
	agis, err := GetAGITree(false)
	if err != nil {
		return
	}
	allAGIs, err := GetAGITree(true)
	if err != nil {
		return
	}

	t := template.Must(template.ParseFiles(README_TMPL))
	f, err := os.Create(path + string(os.PathSeparator) + README)
	if err != nil {
		return
	}
	data := &data{
		Categorys: agis,
		GoRepos:   allAGIs,
	}
	data.GenTime = time.Now().Format("2006-01-02 15:04:05")
	return t.Execute(f, data)
}

// GetAGITree 获取awesome_go_info信息树
func GetAGITree(all bool) (agis []*AwesomeGoInfo, err error) {
	sqlStr := `WITH RECURSIVE subordinates AS (
		SELECT id, parent_id, repo_name, repo_full_name, repo_owner, repo_html_url, repo_description, repo_created_at, repo_pushed_at, repo_homepage, repo_size, repo_forks_count, repo_stargazers_count, repo_subscribers_count, repo_open_issues_count, repo_license_name, repo_license_spdx_id, repo_license_url, repo, category, name, description, homepage, category_html_id, 1 depth FROM awesome_go_info WHERE parent_id = 0 and category=true
		UNION ALL
		SELECT  r.id, r.parent_id, r.repo_name, r.repo_full_name, r.repo_owner, r.repo_html_url, r.repo_description, r.repo_created_at, r.repo_pushed_at, r.repo_homepage, r.repo_size, r.repo_forks_count, r.repo_stargazers_count, r.repo_subscribers_count, r.repo_open_issues_count, r.repo_license_name, r.repo_license_spdx_id, r.repo_license_url, r.repo, r.category, r.name, r.description, r.homepage, r.category_html_id, s.depth+1 depth FROM awesome_go_info r
		inner JOIN subordinates s ON r.parent_id = s.id and r.category=true order by depth desc
		) SELECT id, parent_id, repo_name, repo_full_name, repo_owner, repo_html_url, repo_description, repo_created_at, repo_pushed_at, repo_homepage, repo_size, repo_forks_count, repo_stargazers_count, repo_subscribers_count, repo_open_issues_count, repo_license_name, repo_license_spdx_id, repo_license_url, repo, category, name, description, homepage, category_html_id, depth FROM subordinates`
	rows, err := db.Queryx(sqlStr)
	if err != nil {
		log.Println(err)
		return
	}

	for rows.Next() {
		var tmpAGI AwesomeGoInfo
		rows.StructScan(&tmpAGI)
		tmpAGI.Spaces = getSpace(int64(tmpAGI.Depth) - 1)
		tmpAGI.TitleMarks = getTitleMarks(int64(tmpAGI.Depth))
		agis = append(agis, &tmpAGI)
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
	from awesome_go_info
	where parent_id = ? and repo = true order by repo_stargazers_count desc`
	rows, err := db.Queryx(sqlStr, parentId)
	if err != nil {
		log.Println(err)
		return
	}
	for rows.Next() {
		var tmpAGI AwesomeGoInfo
		rows.StructScan(&tmpAGI)
		tmpAGI.RepoCreatedAtStr = tmpAGI.RepoCreatedAt.Format("2006/01/02")
		tmpAGI.RepoPushedAtStr = tmpAGI.RepoPushedAt.Format("2006-01-02 15:04:05")
		tmpAGI.TimeSince = TimeSince(tmpAGI.RepoPushedAt)
		agis = append(agis, &tmpAGI)
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

const (
	git, clone, add, commit, pull, push, baseDir = `git`, `clone`, `add`, `commit`, `pull`, `push`, `awesome-go`
	awesomeGoRepoURL = `https://github.com/avelino/awesome-go.git`
)

// UpdateAwesomeGoRepo 更新https://github.com/avelino/awesome-go
func UpdateAwesomeGoRepo() (readmeFilePath string, err error){
	_, err = os.Stat("awesome-go")
	if !(err == nil || os.IsExist(err)) {
		log.Printf("awesome-go仓库不存在，开始clone...\n")
		cmd := exec.Command(git, clone, awesomeGoRepoURL)
		output, err := cmd.CombinedOutput()
		fmt.Println(string(output))
		if err != nil {
			return "", err
		}
		log.Println("clone完成")
		// 新clone的仓库无需进一步更新
		return "awesome-go/README.md", nil
	}
	cmd := exec.Command(git, pull)
	cmd.Dir = baseDir
	output, err := cmd.CombinedOutput()
	fmt.Println(string(output))
	if err != nil {
		return "", err
	}
	return"awesome-go/README.md", nil
}

const (
	githubReposAPI     = "https://api.github.com/repos/:owner/:repo"
	githubRateLimitAPI = "https://api.github.com/rate_limit"
	githubDomain       = "https://github.com"
)

var (
	// Contents部分类目列表正则，如：- [Awesome Go](#awesome-go)
	categoryListReg = regexp.MustCompile(`(\s*)- \[(.*)\]\(#(.*)\)`)
	// 扫描中遇到的分类正则，如：## Audio and Music
	categoryReg            = regexp.MustCompile(`#+ (.+)`)
	// 分类的描述正则，如：*Software written in Go.*
	categoryDescReg = regexp.MustCompile(`(\s*)\*([^⬆]+)\*$`)
	// 小分类正则，如：* Testing Frameworks
	littleCategoryReg      = regexp.MustCompile(`(\s*)\* ([a-zA-Z\s]+)$`)
	// 判断是否包含链接正则，如：* [flac](https://github.com/mewkiz/flac) - Native Go FLAC encoder/decoder with support for FLAC streams.
	containsLinkReg        = regexp.MustCompile(`^\s*\* \[.*\]\(.*\)`)

	// 仅包含链接，如* [flac](https://github.com/mewkiz/flac)
	onlyLinkReg            = regexp.MustCompile(`(\s*)\* \[(.*)\]\((.+)\)$`)
	// 链接包含描述， 如：* [flac](https://github.com/mewkiz/flac) - Native Go FLAC encoder/decoder with support for FLAC streams.
	linkWithDescriptionReg = regexp.MustCompile(`(\s*)\* \[(.*?)\]\((.+?)\) - (\S.*[\.\!])`)
	// github 链接正则
	gitHubURLReg           = regexp.MustCompile(`https://github.com/(.+?)/([a-zA-Z0-9_\-\.]+)(.*)$`)

	// 特殊字符正则
	specialCharactersReg = regexp.MustCompile(`[^a-zA-Z0-9_\-\.]+`)
)
// ParseReadmeFile 解析README.md文件
func ParseReadmeFile(filePath, token string) (err error) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return
	}
	lines := strings.Split(string(content), "\n")
	categoryIds := make(map[int]uint)
	var tmpCategoryId, linkCategoryId uint
	var count int64
	for _, line := range lines {
		if categoryListReg.MatchString(line){
			// 处理类目信息
			groups  := categoryListReg.FindStringSubmatch(line)
			SaveBigCategory(categoryIds, len(groups[1]), groups[2],  groups[3])
		}else if categoryReg.MatchString(line){
			// 遇到分类
			groups := categoryReg.FindStringSubmatch(line)
			tmpCategoryId = GetCategoryInfoByName(groups[1])
			linkCategoryId = tmpCategoryId
		} else if categoryDescReg.MatchString(line) {
			// 分类描述
			if tmpCategoryId == 0{
				continue
			}
			groups := categoryDescReg.FindStringSubmatch(line)
			ModifyCategoryDescription(tmpCategoryId, groups[2])
		} else if littleCategoryReg.MatchString(line) {
			// 小分类
			if tmpCategoryId == 0{
				continue
			}
			groups := littleCategoryReg.FindStringSubmatch(line)
			name := groups[2]
			categoryHtmlId := ConvertToCategoryHtmlId(name)
			linkCategoryId = SaveCategory(tmpCategoryId, name, categoryHtmlId)
		} else if containsLinkReg.MatchString(line) && strings.Contains(line, githubDomain){
			// 含有链接,且为GitHub仓库
			var githubRepoLink, repoDescription string
			if onlyLinkReg.MatchString(line) {
				subMatchs := onlyLinkReg.FindStringSubmatch(line)
				githubRepoLink = subMatchs[3]
			} else if linkWithDescriptionReg.MatchString(line) {
				subMatchs := linkWithDescriptionReg.FindStringSubmatch(line)
				githubRepoLink = subMatchs[3]
				repoDescription = subMatchs[4]
			}
			if !gitHubURLReg.MatchString(githubRepoLink) {
				continue
			}
			groups := gitHubURLReg.FindStringSubmatch(githubRepoLink)
			extraStr := groups[3]
			// 非仓库地址
			if len(extraStr) > 0 && strings.HasPrefix(extraStr, "/") {
				continue
			}
			repoOwner, repoName := groups[1], groups[2]
			name := repoOwner + "/" + repoName
			ok, err := GitHubAPIReqControl(token)
			if err != nil {
				log.Println(err)
				continue
			}
			if !ok {
				log.Println("API请求速率限制")
				continue
			}
			count++
			log.Printf("%d 开始请求仓库(%s)信息\n", count, name)
			tmpAGI, err := GetRepoInfo(repoOwner, repoName, token)
			if err != nil {
				log.Println(err)
				continue
			}
			tmpAGI.ParentId = linkCategoryId
			if repoDescription != ""{
				tmpAGI.Description = repoDescription
			}
			// GitHub仓库地址中的特殊字符会自动去掉，导致README.md中的名字与库中的不一致，因此存储以GitHub接口返回的full_name为准
			id, e := GetAGI(tmpAGI.RepoFullName, true, false)
			if e != nil {
				_, err = db.Exec(`insert into awesome_go_info (repo_name, repo_full_name, repo_owner, repo_html_url, repo_description, repo_created_at, repo_pushed_at, repo_homepage, repo_size, repo_forks_count, repo_stargazers_count, repo_subscribers_count, repo_open_issues_count, repo_license_name, repo_license_spdx_id, repo_license_url, parent_id, repo, category, name, description, homepage, category_html_id, created_at, updated_at)
values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, datetime('now'), datetime('now'))`, tmpAGI.RepoName, tmpAGI.RepoFullName, tmpAGI.RepoOwner, tmpAGI.RepoHtmlURL, tmpAGI.RepoDescription, tmpAGI.RepoCreatedAt, tmpAGI.RepoPushedAt, tmpAGI.RepoHomepage, tmpAGI.RepoSize, tmpAGI.RepoForksCount, tmpAGI.RepoStargazersCount, tmpAGI.RepoSubscribersCount, tmpAGI.RepoOpenIssuesCount, tmpAGI.RepoLicenseName, tmpAGI.RepoLicenseSpdxId, tmpAGI.RepoLicenseURL, tmpAGI.ParentId, tmpAGI.Repo, tmpAGI.Category, tmpAGI.Name, tmpAGI.Description, tmpAGI.Homepage, tmpAGI.CategoryHtmlId)
				if err != nil{
					log.Printf("新增仓库信息(%s)异常: %+v\n", name, err)
				}
			} else {
				_, err = db.Exec(`update awesome_go_info set repo_name = ?, repo_full_name = ?, repo_owner = ?, repo_html_url = ?, repo_description = ?, repo_created_at = ?, repo_pushed_at = ?, repo_homepage = ?, repo_size = ?, repo_forks_count = ?, repo_stargazers_count = ?, repo_subscribers_count = ?, repo_open_issues_count = ?, repo_license_name = ?, repo_license_spdx_id = ?, repo_license_url = ?, parent_id = ?, repo = ?, category = ?, name = ?, description = ?, homepage = ?, category_html_id = ?, updated_at = datetime('now') where id = ?`, tmpAGI.RepoName, tmpAGI.RepoFullName, tmpAGI.RepoOwner, tmpAGI.RepoHtmlURL, tmpAGI.RepoDescription, tmpAGI.RepoCreatedAt, tmpAGI.RepoPushedAt, tmpAGI.RepoHomepage, tmpAGI.RepoSize, tmpAGI.RepoForksCount, tmpAGI.RepoStargazersCount, tmpAGI.RepoSubscribersCount, tmpAGI.RepoOpenIssuesCount, tmpAGI.RepoLicenseName, tmpAGI.RepoLicenseSpdxId, tmpAGI.RepoLicenseURL, tmpAGI.ParentId, tmpAGI.Repo, tmpAGI.Category, tmpAGI.Name, tmpAGI.Description, tmpAGI.Homepage, tmpAGI.CategoryHtmlId, id)
				if err != nil{
					log.Printf("修改仓库信息(%s)异常: %+v\n", name, err)
				}
			}
		}
	}
	return nil
}

// GetAGI 获取awesome_go_info表数据
func GetAGI(name string, repo bool, category bool) (id uint, err error) {
	err = db.Get(&id, `select id from awesome_go_info where name = ? and repo = ? and category = ?`, name, repo, category)
	return
}

// ModifyCategoryDescription 修改分类描述信息
func ModifyCategoryDescription(id uint, description string){
	_, err := db.Exec(`update awesome_go_info set description = ? where id = ?`, description, id)
	if err != nil{
		log.Printf("修改分类描述信息异常： %+v", err)
	}
}

// GetCategoryInfoByName 根据README.md重的分类名称获取分类信息
func GetCategoryInfoByName(name string) uint{
	var awesomeGoInfo AwesomeGoInfo
	categoryHtmlId := ConvertToCategoryHtmlId(name)
	err  := db.Get(&awesomeGoInfo, `select * from awesome_go_info where category_html_id = ? and category = true`, categoryHtmlId)
	if err != nil{
		log.Printf("根据分类名称 %s 获取分类信息失败: %+v", name, err)
		return 0
	}
	return awesomeGoInfo.ID
}

// SaveGithubRepo 保存 github 仓库信息
func SaveGithubRepo(parentID uint, name, categoryHtmlId string){
	var repo AwesomeGoInfo
	err := db.Get(&repo, `select * from awesome_go_info where category_html_id = ? and category = true`, categoryHtmlId)
	if err != nil && err != sql.ErrNoRows{
		log.Printf("获取仓库信息异常: %+v", err)
		return
	}
	if err == sql.ErrNoRows{
		_, err = db.Exec(`insert into awesome_go_info (parent_id, repo, category, name, category_html_id, created_at, updated_at) values (?, ?, ?, ?, ?, datetime('now'), datetime('now'))`, parentID, false, true, name, categoryHtmlId)
		if err != nil{
			log.Printf("新增仓库信息异常: %+v", err)
			return
		}
	}else{
		_, err = db.Exec(`update awesome_go_info set parent_id = ? where category_html_id = ? and category = true`, parentID, categoryHtmlId)
		if err != nil{
			log.Printf("修改仓库信息异常: %+v", err)
			return
		}
	}
}

// SaveCategory 保存类目信息
func SaveBigCategory(categoryIds map[int]uint, spaceCount int, name, categoryHtmlId string){
	parentID := categoryIds[spaceCount-4]
	categoryIds[spaceCount] = SaveCategory(parentID, name, categoryHtmlId)
}

// SaveCategory 保存类目信息
func SaveCategory(parentID uint, name, categoryHtmlId string) uint{
	var category AwesomeGoInfo
	err := db.Get(&category, `select * from awesome_go_info where category_html_id = ? and category = true`, categoryHtmlId)
	if err != nil && err != sql.ErrNoRows{
		log.Printf("获取类目信息异常: %+v", err)
		return 0
	}
	if err == sql.ErrNoRows{
		_, err = db.Exec(`insert into awesome_go_info (parent_id, repo, category, name, category_html_id, created_at, updated_at) values (?, ?, ?, ?, ?, datetime('now'), datetime('now'))`, parentID, false, true, name, categoryHtmlId)
		if err != nil{
			log.Printf("新增类目信息异常: %+v", err)
			return 0
		}
	}else{
		_, err = db.Exec(`update awesome_go_info set parent_id = ? where category_html_id = ? and category = true`, parentID, categoryHtmlId)
		if err != nil{
			log.Printf("修改类目信息异常: %+v", err)
			return 0
		}
	}
	// 获取类目id
	err = db.Get(&category, `select * from awesome_go_info where category_html_id = ? and category = true`, categoryHtmlId)
	if err != nil && err != sql.ErrNoRows{
		log.Printf("获取类目信息异常: %+v", err)
		return 0
	}
	if err == sql.ErrNoRows{
		log.Printf("类目 %s 未新增成功", name)
		return 0
	}
	return category.ID
}

// ConvertToCategoryHtmlId 转换成CategoryHtmlId
func ConvertToCategoryHtmlId(categoryName string) (id string) {
	id = specialCharactersReg.ReplaceAllString(categoryName, "-")
	id = strings.Trim(id, "-")
	return strings.ToLower(id)
}