package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"regexp"
	"strings"
	"text/template"
	"time"

	"github.com/spf13/viper"
)

const (
	sourceFileURL      = "https://raw.githubusercontent.com/avelino/awesome-go/master/README.md"
	githubReposAPI     = "https://api.github.com/repos/:owner/:repo?access_token=OAUTH-TOKEN"
	githubRateLimitAPI = "https://api.github.com/rate_limit?access_token=OAUTH-TOKEN"
	githubDomain       = "https://github.com"

	RFILES_PATH      = "rfiles"
	RFILES_FILE_NAME = "README-FORMAT_DATE.md"
	RFILE            = RFILES_PATH + string(os.PathSeparator) + RFILES_FILE_NAME
	README_TMPL      = "tmpl.md"
	README           = "README.md"
)

// init 创建必要的文件路径
func init() {
	os.MkdirAll(RFILES_PATH, os.ModePerm)
}

var (
	reCategoryLi          = regexp.MustCompile(`(\s*)- \[(.*)\]\(#(.*)\)`)
	reCategory            = regexp.MustCompile(`#+ (.+)`)
	reCategoryDescription = regexp.MustCompile(`(\s*)\*(.*)\*$`)
	reContainsLink        = regexp.MustCompile(`\* \[.*\]\(.*\)`)
	reOnlyLink            = regexp.MustCompile(`(\s*)\* \[(.*)\]\((.+)\)$`)
	reLinkWithDescription = regexp.MustCompile(`(\s*)\* \[(.*?)\]\((.+?)\) - (\S.*[\.\!])`)
	reLittleCategory      = regexp.MustCompile(`(\s*)\* ([a-zA-Z\s]*)$`)
	reGitHubURL           = regexp.MustCompile(`https://github.com/(.+?)/([a-zA-Z0-9_\-\.]+)(.*)$`)

	reSpecialCharacters = regexp.MustCompile(`[^a-zA-Z0-9_\-\.]+`)
)

// CurREADMEFile 获取当天的README.md文件
// 返回文件路径filePath及该文件是否存在exist信息
func CurREADMEFile() (filePath string, exist bool) {
	filePath = strings.Replace(RFILE, "FORMAT_DATE", time.Now().Format("20060102"), -1)
	_, err := os.Stat(filePath)
	return filePath, err == nil || os.IsExist(err)
}

// DownloadREADMEFile 下载awesome-go项目中的README.md文件
// 返回文件路径
func DownloadREADMEFile() (filePath string, err error) {
	filePath, exist := CurREADMEFile()
	// 文件已存在
	if exist {
		return
	}
	// 下载文件
	res, err := http.Get(sourceFileURL)
	if err != nil {
		return
	}
	defer res.Body.Close()
	// 保存文件
	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	writer.Write(bytes)
	writer.Flush()
	return
}

// ParseREADMEFile 解析awesome-go中的README.md文件,并存入数据库中
func ParseREADMEFile(accessToken string, filePath string) (err error) {
	input, err := ioutil.ReadFile(filePath)
	if err != nil {
		return
	}

	lines := strings.Split(string(input), "\n")
	categoryIds := make(map[int]uint)
	var tmpCategoryId uint = 0
	var linkCategoryId uint = 0
	var count int64 = 0
	name := ""
	for _, line := range lines {
		// 分类目录
		if reCategoryLi.MatchString(line) {
			subMatchs := reCategoryLi.FindStringSubmatch(line)
			name = subMatchs[2]
			categoryHtmlId := subMatchs[3]
			spaceCount := len(subMatchs[1])

			agi, e := GetAGIByCategoryHtmlId(categoryHtmlId)
			if e != nil {
				agi = &AwesomeGoInfo{
					ParentId:       categoryIds[spaceCount-8],
					Repo:           false,
					Category:       true,
					Name:           name,
					CategoryHtmlId: categoryHtmlId,
				}
				SaveAGI(agi)
			} else {
				ModifyAGIParentIdByCategoryHtmlId(categoryIds[spaceCount-8], categoryHtmlId)
			}
			categoryIds[spaceCount-4] = agi.ID
		} else if reCategory.MatchString(line) {
			// 遇到分类
			subMatchs := reCategory.FindStringSubmatch(line)
			name = subMatchs[1]
			categoryHtmlId := getCategoryHtmlId(name)

			agi, e := GetAGIByCategoryHtmlId(categoryHtmlId)
			if e != nil {
				log.Printf("分类%s不存在", name)
				continue
			}
			tmpCategoryId = agi.ID
			linkCategoryId = tmpCategoryId
		} else if reCategoryDescription.MatchString(line) {
			// 分类描述
			subMatchs := reCategoryDescription.FindStringSubmatch(line)
			description := subMatchs[2]
			UpdateAGIDescription(description, tmpCategoryId)
		} else if reLittleCategory.MatchString(line) {
			// 小分类
			subMatchs := reLittleCategory.FindStringSubmatch(line)
			name = subMatchs[2]
			categoryHtmlId := getCategoryHtmlId(name)
			agi, e := GetAGIByCategoryHtmlId(categoryHtmlId)
			if e != nil {
				agi = &AwesomeGoInfo{
					ParentId:       tmpCategoryId,
					Repo:           false,
					Category:       true,
					Name:           name,
					CategoryHtmlId: categoryHtmlId,
				}
				SaveAGI(agi)
			} else {
				ModifyAGIParentIdByCategoryHtmlId(tmpCategoryId, categoryHtmlId)
			}
			linkCategoryId = agi.ID
		} else if reContainsLink.MatchString(line) && strings.Contains(line, githubDomain) {
			// 含有链接,且为GitHub仓库
			githubRepoLink := ""
			repoDescription := ""
			if reOnlyLink.MatchString(line) {
				subMatchs := reOnlyLink.FindStringSubmatch(line)
				githubRepoLink = subMatchs[3]
			} else if reLinkWithDescription.MatchString(line) {
				subMatchs := reLinkWithDescription.FindStringSubmatch(line)
				githubRepoLink = subMatchs[3]
				repoDescription = subMatchs[4]
			}
			if githubRepoLink != "" && reGitHubURL.MatchString(githubRepoLink) {
				subMatchs := reGitHubURL.FindStringSubmatch(githubRepoLink)
				extraStr := subMatchs[3]
				// 非仓库地址
				if len(extraStr) > 0 && strings.HasPrefix(extraStr, "/") {
					continue
				}

				repoOwner, repoName := subMatchs[1], subMatchs[2]
				name = repoOwner + "/" + repoName
				ok, err := GitHubAPIReqControl(accessToken)
				if err != nil {
					log.Println(err)
					continue
				}
				if !ok {
					log.Println("API请求速率限制")
					continue
				}
				count++
				log.Println(count, " 开始请求仓库信息", name)
				tmpAGI, err := GetRepoInfo(repoOwner, repoName, accessToken)
				if err != nil {
					log.Println(err)
					continue
				}
				tmpAGI.ParentId = linkCategoryId
				if repoDescription != "" {
					tmpAGI.Description = repoDescription
				}

				// GitHub仓库地址中的特殊字符会自动去掉，导致README.md中的名字与库中的不一致，因此存储以GitHub接口返回的full_name为准
				agi, e := GetAGI(tmpAGI.RepoFullName, true, false)
				if e != nil {
					SaveAGI(tmpAGI)
				} else {
					tmpAGI.ID = agi.ID
					UpdateAGIGithubInfo(tmpAGI)
				}

				grr := GithubRepoRecord{
					RepoInfo: tmpAGI.RepoInfo,
				}
				SaveGRR(&grr)
				log.Print("请求完成。")
			}
		}
	}
	log.Printf("解析%s文件完毕", filePath)
	return
}

// GitHubAPIReqControl 控制是否可以进行GitHub的API调用
func GitHubAPIReqControl(accessToken string) (ok bool, err error) {
	// 通过访问rate_limit API，校验access_token是否有效，及通过剩余次数控制API的访问时机
	apiURL := strings.Replace(githubRateLimitAPI, "OAUTH-TOKEN", accessToken, -1)
	res, err := http.Get(apiURL)
	if err != nil {
		return
	}
	defer res.Body.Close()
	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}
	var v interface{}
	err = json.Unmarshal(bytes, &v)
	if err != nil {
		return
	}
	if rateLimitMap, ok := v.(map[string]interface{}); ok {
		//resources信息
		if resourcesMap, ok := rateLimitMap["resources"].(map[string]interface{}); ok {
			if coreMap, ok := resourcesMap["core"].(map[string]interface{}); ok {
				var remaining float64
				var limit float64
				if remainingFloat64, ok := coreMap["remaining"].(float64); ok {
					remaining = math.RoundToEven(remainingFloat64)
				}
				if limitFloat64, ok := coreMap["limit"].(float64); ok {
					limit = math.RoundToEven(limitFloat64)
				}
				if remaining > 0 && limit > 0 {
					time.Sleep(3600 * 1e9 / time.Duration(limit))
					return true, nil
				}
			}
		}
	}
	return
}

// 获取github仓库信息
func GetRepoInfo(repoOwner, repoName, accessToken string) (agi *AwesomeGoInfo, err error) {
	repoFullName := repoOwner + "/" + repoName
	apiURL := strings.Replace(githubReposAPI, ":owner/:repo", repoFullName, -1)
	apiURL = strings.Replace(apiURL, "OAUTH-TOKEN", accessToken, -1)
	res, err := http.Get(apiURL)
	if err != nil {
		return
	}
	defer res.Body.Close()
	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}
	var v interface{}
	err = json.Unmarshal(bytes, &v)
	if err != nil {
		return
	}
	if repoMap, ok := v.(map[string]interface{}); ok {
		if fullName, ok := repoMap["full_name"].(string); !ok || fullName == "" {
			return agi, fmt.Errorf("从GitHub获取%s仓库信息失败", repoFullName)
		}
		agi = &AwesomeGoInfo{
			RepoInfo: RepoInfo{
				RepoName:             repoMap["name"].(string),
				RepoFullName:         repoMap["full_name"].(string),
				RepoHtmlURL:          repoMap["html_url"].(string),
				RepoForksCount:       uint(repoMap["forks_count"].(float64)),
				RepoStargazersCount:  uint(repoMap["stargazers_count"].(float64)),
				RepoSubscribersCount: uint(repoMap["subscribers_count"].(float64)),
				RepoOpenIssuesCount:  uint(repoMap["open_issues_count"].(float64)),
			},
			Repo:     true,
			Category: false,
			Name:     repoMap["full_name"].(string),
		}
		//可能不存在的一些信息，需要进行判断
		if homepage, ok := repoMap["homepage"].(string); ok {
			agi.RepoHomepage = homepage
			agi.Homepage = homepage
		}
		if description, ok := repoMap["description"].(string); ok {
			agi.RepoDescription = description
			agi.Description = description
		}
		//作者信息
		if ownerMap, ok := repoMap["owner"].(map[string]interface{}); ok {
			agi.RepoOwner = ownerMap["login"].(string)
		}
		//证书信息
		if licenseMap, ok := repoMap["license"].(map[string]interface{}); ok {
			if licenseName, ok := licenseMap["name"].(string); ok {
				agi.RepoLicenseName = licenseName
			}
			if licenseSpdxId, ok := licenseMap["spdx_id"].(string); ok {
				agi.RepoLicenseSpdxId = licenseSpdxId
			}
			if licenseURL, ok := licenseMap["url"].(string); ok {
				agi.RepoLicenseURL = licenseURL
			}
		}
		//日期时间信息
		if createAtStr, ok := repoMap["created_at"].(string); ok {
			createAt, err := time.Parse(time.RFC3339, createAtStr)
			if err == nil {
				agi.RepoCreatedAt = &createAt
			}
		}
		if pushedAtStr, ok := repoMap["pushed_at"].(string); ok {
			pushedAt, err := time.Parse(time.RFC3339, pushedAtStr)
			if err == nil {
				agi.RepoPushedAt = &pushedAt
			}
		}
		//size可能带小数
		if sizeFloat64, ok := repoMap["size"].(float64); ok {
			size := math.RoundToEven(sizeFloat64)
			agi.RepoSize = uint(size)
		}
	}
	return
}

func getCategoryHtmlId(categoryName string) (id string) {
	id = reSpecialCharacters.ReplaceAllString(categoryName, "-")
	id = strings.Trim(id, "-")
	return strings.ToLower(id)
}

type data struct {
	Categorys []*AwesomeGoInfo
	GoRepos   []*AwesomeGoInfo
	GenTime   string
}

//GenerateMd 生成README.md文件
func GenerateMd() (err error) {
	agis, err := GetAGITree(false)
	if err != nil {
		return
	}
	allAGIs, err := GetAGITree(true)
	if err != nil {
		return
	}

	t := template.Must(template.ParseFiles(README_TMPL))
	f, err := os.Create(viper.GetString("rpath") + string(os.PathSeparator) + README)
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
