package mysql

import (
	"NothingBlog/settings"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Db *gorm.DB

// 数据库初始化
func Init(cfg *settings.MysqlConfig, appModel string) (err error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Schema,
	)
	// 使用gorm打开mysql数据库
	Db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	//db, err = sqlx.Connect("mysql", dsn)

	if err != nil {
		zap.L().Error("Connect mysql failed...", zap.Error(err))
		return
	}
	sqlDb, err := Db.DB()
	if err != nil {
		zap.L().Error("get db failed...", zap.Error(err))
		return err
	}

	sqlDb.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDb.SetMaxOpenConns(cfg.MaxOpenConns)
	// 设置数据库连接池的最大连接数和空闲连接数
	//db.SetMaxIdleConns(cfg.MaxIdleConns)
	//db.SetMaxOpenConns(cfg.MaxOpenConns)
	return
}

func Close() {
	//db.Close()
	sqlDb, _ := Db.DB()
	sqlDb.Close()
}
