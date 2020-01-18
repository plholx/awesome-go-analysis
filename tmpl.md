# Awesome Go Info

go语言开源项目列表，项目分类及GitHub上的开源项目数据完全来自于[awesome-go](https://github.com/avelino/awesome-go) 的[README.md](https://github.com/avelino/awesome-go/blob/master/README.md)文件，通过调用GitHub的API获取仓库信息，展示项目的star数、watch数等，方便查看go语言开源项目的一些相关信息。

_该文件仅包含[awesome-go](https://github.com/avelino/awesome-go) [README.md](https://github.com/avelino/awesome-go/blob/master/README.md)文件中列出的在GitHub上开源的优秀项目，不罗列其它golang相关的网站_
_该文件中的GitHub仓库信息数据会在每天凌晨1点左右更新,当前数据更新于{{.GenTime}}_

{{range .Categorys}}{{.Spaces}}- [{{.Name}}](#{{.CategoryHtmlId}})
{{end}}

{{range .GoRepos}}
    {{- if .Category}}
        {{- "\n\n"}}{{.TitleMarks}} {{.Name}}
        {{if .Description}}
            {{- "\n*"}}{{.Description}}*
        {{- end}}
        {{- if .WithReposTable}}
            {{- "\n\n|"}} Go_repository    | Stars      | Watchers   | Created_at | Latest_push | Description |
            {{- "\n|"}} :--------- | ---------:| ---------:|:---------:|:---------:|:--------- |
        {{- end}}
    {{- else if .Repo}}
         {{- "\n|"}} [{{.RepoName}}]({{.RepoHtmlURL -}}) | {{.RepoStargazersCount}} | {{.RepoSubscribersCount}} | {{.RepoCreatedAtStr}} | {{.TimeSince}} | {{.Description}} |
    {{- end}}
{{- end}}

> 这仅仅是个人练手学习go语言的一个小项目，欢迎指点 <plholx@126.com> ^_^
> 更专业的go开源项目分析请移步 [Awesome Go](https://go.libhunt.com/)
