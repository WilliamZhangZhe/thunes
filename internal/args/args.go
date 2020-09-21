// Package args 解析程序启动参数信息

package args

import (
	"flag"
	"fmt"
	"os"

	"github.com/thunes/pkg/errwrap"
)

type Args struct {
	CfgPath  string // 配置文件地址
	Graceful bool   // 是否是优雅重启
}

// Parse 解析程序启动参数，例如-c为配置文件地址，-gr为优雅重启
func (A *Args) Parse() (err error) {
	cfgPath := flag.String(
		"c",
		"./thunes.toml",
		fmt.Sprintf("thunes service config file, default %s", "./thunes.toml"))

	gracefulRestart := flag.Bool(
		"gr",
		false,
		"graceful restart app",
	)

	flag.Parse()

	A.CfgPath = *cfgPath
	if _, err := os.Stat(A.CfgPath); err != nil {
		return errwrap.WithContext(err, "stat config file: "+A.CfgPath)
	}

	A.Graceful = *gracefulRestart

	return
}
