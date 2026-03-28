package main

import (
	"log"
	"moxin-zhicheng/internal/config"
	"moxin-zhicheng/internal/database"
	"moxin-zhicheng/internal/logger"
	"moxin-zhicheng/internal/redis"
	"moxin-zhicheng/routes"

	"github.com/gin-contrib/pprof"

	"github.com/gin-gonic/gin"
)

func main() {
	logger.InitLogger("dev")
	// 1. 初始化配置 (Viper)
	config.InitConfig()
	// 2. 初始化数据库 (GORM)
	database.InitDB()

	// 3. 初始化 Gin 引擎
	// gin.Default() 默认包含了 Logger 和 Recovery 中间件
	r := gin.Default()
	pprof.Register(r)
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Token")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})
	redis.InitRedis()
	defer redis.CloseRedis()
	routes.SetupRoutes(r)
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("服务器启动失败：%v", err)
	}

}
