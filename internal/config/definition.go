package config

import "fmt"

// DbServer db server
type DbServer struct {
	IP       string `mapstructure:"ip"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Pwd      string `mapstructure:"pwd"`
	Database string `mapstructure:"database"`
}

// APIService  rest api service
type APIService struct {
	Host string
	Port int
}

// API 返回接口地址
func (S *APIService) API(URI string) string {
	return fmt.Sprintf("%s/%s", S.Addr(), URI)
}

// Addr 返回服务的域名地址
func (S *APIService) Addr() string {
	if S.Port > 0 {
		return fmt.Sprintf("%s:%d", S.Host, S.Port)
	}

	return S.Host
}
