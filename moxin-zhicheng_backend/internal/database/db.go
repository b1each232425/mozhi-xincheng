package database

import (
	"fmt"
	"log/slog"
	"moxin-zhicheng/internal/config"
	"moxin-zhicheng/internal/logger"
	model "moxin-zhicheng/internal/models"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	glog "gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() {
	c := config.Conf.Database
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.User,
		c.Password,
		c.Host,
		c.Port,
		c.DBName,
	)
	var level glog.LogLevel
	if config.Conf.AppMode == "dev" {
		level = glog.Info
	} else {
		level = glog.Error
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: glog.Default.LogMode(level),
	})

	if err != nil {
		logger.Error("数据库连接失败", err, slog.String("dsn", dsn))
		panic("Failed to connect to database")
	}
	sqlDB, err := db.DB()
	if err != nil {
		logger.Error("获取 sql.DB 失败", err)
		panic(err)
	}

	// --- 关键配置：连接池优化 ---
	// 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(10)
	// 设置打开数据库连接的最大数量
	sqlDB.SetMaxOpenConns(100)
	// 设置连接可复用的最大时间
	sqlDB.SetConnMaxLifetime(time.Hour)

	DB = db
	err = DB.AutoMigrate(
		&model.Poetry{},
		&model.PoetryTag{},
		&model.PoetryTagRelation{},
		&model.Article{})
	if err != nil {
		logger.Error("数据库自动迁移失败", err)
		panic(err)
	}
	logger.Info("数据库初始化成功", slog.String("db", c.DBName))
}
