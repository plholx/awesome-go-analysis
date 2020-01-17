package conf

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	pflag.StringP("token", "t", "xxx", "GitHub API access_token")
	pflag.StringP("dsurl", "d", "./data_test.db", "数据库地址")
	pflag.StringP("rpath", "r", "rfiles", "生成的README路径")
	pflag.Parse()

	// 将命令行参数绑定到viper对象
	viper.BindPFlags(pflag.CommandLine)
}
