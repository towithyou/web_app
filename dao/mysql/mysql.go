package mysql

import (
	"fmt"

	"github.com/towithyou/web_app/settings"

	"go.uber.org/zap"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

func Init(cfg *settings.MySqlConfig) (err error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
		cfg.Charset,
	)
	// 也可以使用MustConnect连接不成功就panic
	db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		zap.L().Error("connect DB failed, err", zap.Error(err))
		return
	}
	db.SetMaxOpenConns(cfg.MaxConn)
	db.SetMaxIdleConns(cfg.MaxIdleConn)
	return
}

func Close() {
	_ = db.Close()
}
