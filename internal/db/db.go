// Package db 提供数据库资源相关操作封装，主要包括mysql数据库连接资源
package db

import (
	"time"

	"github.com/go-xorm/xorm"
)

var (
	_mysql *xorm.EngineGroup
)

// Init 初始化数据库相关资源（数据库连接池）
func Init() (err error) {
	if err = initMySQL(); err != nil {
		return
	}

	for i := 0; i < ConnectMaxRetry && err != nil; i++ {
		err = _mysql.Ping()
		if err == nil {
			break
		}

		time.Sleep(time.Second)
	}

	_mysql.ShowSQL(true)
	_mysql.SetMaxOpenConns(20)

	return
}

// Release 释放数据库相关资源
func Relese() {
	if _mysql != nil {
		_mysql.Close()
	}
}

// Cli 获取DB client
func Cli() *xorm.EngineGroup {
	return _mysql
}
