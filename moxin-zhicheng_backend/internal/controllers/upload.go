package controllers

import (
	"moxin-zhicheng/utils"

	"github.com/gin-gonic/gin"
)

func UploadArticleImage(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(400, gin.H{"code": 400, "msg": "未发现图片文件"})
		return
	}

	// 调用之前写的 utils.UploadImage 传到 OSS
	url, err := utils.UploadImage(file)
	if err != nil {
		c.JSON(500, gin.H{"code": 500, "msg": "云端存储失败"})
		return
	}

	// 返回 OSS 的完整 URL 给前端编辑器
	c.JSON(200, gin.H{
		"code": 200,
		"url":  url,
		"msg":  "墨宝已入云端",
	})
}
