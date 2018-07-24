package mysql

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/guobin8205/api_demo/utils/config"
)

func NewMysqlCon(user, pwd, host, port, database, charset string, maxidle, maxactive int) *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s",
		user,
		pwd,
		host,
		port,
		database,
		charset,
	)
	db, err := gorm.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	db.DB().SetMaxIdleConns(maxidle)
	db.DB().SetMaxOpenConns(maxactive)

	//是否打印sql日志,默认打印
	db.LogMode(true)
	logMode := conf.MustString("log.sql.mode", "1")
	if logMode != "1" {
		db.LogMode(false)
	}

	return db
}
