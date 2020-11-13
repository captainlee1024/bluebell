// Package mysql provides ...
package mysql

import (
	"fmt"

	"github.com/captainlee1024/bluebell/settings"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

var db *sqlx.DB

// Init 初始化 MySQL 连接
func Init(cfg *settings.MySQLConfig) (err error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true",
		//viper.GetString("mysql.user"),
		//viper.GetString("mysql.password"),
		//viper.GetString("mysql.host"),
		//viper.GetInt("mysql.port"),
		//viper.GetString("mysql.dbname"),
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DbName,
	)

	db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		zap.L().Error("connect DB failed", zap.Error(err))
		return
	}

	//db.SetMaxOpenConns(viper.GetInt("mysql.max_open_conns"))
	//db.SetMaxIdleConns(viper.GetInt("mysql.max_idel_conns"))
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	return
}

// Close 关闭 MySQL 连接
func Close() {
	_ = db.Close()
}
