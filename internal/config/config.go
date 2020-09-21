package config

import (
	"fmt"

	"github.com/spf13/viper"
	"github.com/thunes/pkg/errwrap"
)

var (
	// Self 本服务配置
	Self = APIService{}

	// DB 数据库配置信息
	DB = map[string]DbServer{}
)

// Load 加载配置信息
func Load(path string) (err error) {
	var (
		data = viper.New()
	)

	data.SetConfigType("toml")
	data.SetConfigFile(path)

	err = errwrap.WithContext(data.ReadInConfig(), "读取配置文件")
	if err != nil {
		return
	}

	if err = loadKey(data, "self", &Self); err != nil {
		return
	}

	if err = loadKey(data, "dbservers", &DB); err != nil {
		return
	}

	return
}

func loadKey(data *viper.Viper, key string, target interface{}) error {
	ctx := fmt.Sprintf("读取%s", key)
	return errwrap.WithContext(data.UnmarshalKey(key, target), ctx)
}
