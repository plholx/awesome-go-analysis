package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"strings"
	"time"
)

// GitHubApiReq GitHub请求
func GitHubApiReq(token, method, url string, body io.Reader) (resp *http.Response, err error) {
	client := http.Client{}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return
	}
	req.Header.Add("Authorization", " token "+token)
	resp, err = client.Do(req)
	return
}

// GitHubAPIReqControl 控制是否可以进行GitHub的API调用
// 通过访问rate_limit API，校验access_token是否有效，及通过剩余次数控制API的访问时机
func GitHubAPIReqControl(accessToken string) (ok bool, err error) {
	res, err := GitHubApiReq(accessToken, "GET", githubRateLimitAPI, nil)
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
	res, err := GitHubApiReq(accessToken, "GET", apiURL, nil)
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
			RepoName:             repoMap["name"].(string),
			RepoFullName:         repoMap["full_name"].(string),
			RepoHtmlURL:          repoMap["html_url"].(string),
			RepoForksCount:       uint(repoMap["forks_count"].(float64)),
			RepoStargazersCount:  uint(repoMap["stargazers_count"].(float64)),
			RepoSubscribersCount: uint(repoMap["subscribers_count"].(float64)),
			RepoOpenIssuesCount:  uint(repoMap["open_issues_count"].(float64)),
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
