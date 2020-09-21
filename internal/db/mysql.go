package db

import (
	"errors"
	"fmt"

	"github.com/go-xorm/xorm"
	"github.com/thunes/internal/config"

	_ "github.com/go-sql-driver/mysql"
)

var (
	ConnectMaxRetry = 10
)

func initMySQL() (err error) {
	srvInfo, ok := config.DB["thunes"] // dbservers.thunes
	if !ok {
		return errors.New("db config not found")
	}

	conn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?parseTime=True&charset=utf8mb4",
		srvInfo.User,
		srvInfo.Pwd,
		srvInfo.IP,
		srvInfo.Port,
		srvInfo.Database)

	_mysql, err = xorm.NewEngineGroup("mysql", []string{conn, conn})

	return
}
