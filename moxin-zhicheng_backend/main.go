package main

import (
	"github.com/gin-contrib/pprof"
	"log"
	"moxin-zhicheng/internal/config"
	"moxin-zhicheng/internal/database"
	"moxin-zhicheng/internal/logger"
	"net/http"

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
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("服务器启动失败：%v", err)
	}
	// 4. 定义一个简单的测试接口
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
}
