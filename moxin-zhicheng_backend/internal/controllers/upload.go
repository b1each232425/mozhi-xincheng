package controllers

import (
	"moxin-zhicheng/internal/response"
	"moxin-zhicheng/utils"

	"github.com/gin-gonic/gin"
)

func UploadArticleImage(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		response.BadRequest(c, "未发现图片文件")
		return
	}

	// 调用之前写的 utils.UploadImage 传到 OSS
	url, err := utils.UploadImage(file)
	if err != nil {
		response.ServerError(c, "云端存储失败")
		return
	}

	// 返回 OSS 的完整 URL 给前端编辑器
	response.SuccessWithMsg(c, "墨宝已入云端", gin.H{"url": url})
}
