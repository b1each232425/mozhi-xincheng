package controllers

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

// HandleVideoUpload 处理视频上传并调用 FFmpeg 标准化
func HandleVideoUpload(c *gin.Context) {
	// 1. 获取上传文件
	file, err := c.FormFile("video")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "未接收到视频文件"})
		return
	}

	// 2. 创建本地存储目录 (建议放在根目录的 uploads)
	uploadDir := "./uploads/raw"
	processedDir := "./uploads/processed"
	os.MkdirAll(uploadDir, os.ModePerm)
	os.MkdirAll(processedDir, os.ModePerm)

	// 3. 保存原始文件
	inputPath := filepath.Join(uploadDir, file.Filename)
	if err := c.SaveUploadedFile(file, inputPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存原始视频失败"})
		return
	}

	// 4. 定义处理后的输出路径
	outputFilename := fmt.Sprintf("std_%s", file.Filename)
	outputPath := filepath.Join(processedDir, outputFilename)

	// 5. 调用 FFmpeg (执行标准化：1080p, 30fps, 无音频)
	// 这里的指令专门为下一步的 Python 视觉识别优化
	cmd := exec.Command("ffmpeg", "-i", inputPath,
		"-vf", "scale=1920:1080:force_original_aspect_ratio=decrease,pad=1920:1080:(ow-iw)/2:(oh-ih)/2",
		"-r", "30",
		"-c:v", "libx264",
		"-preset", "fast",
		"-crf", "23",
		"-an",
		"-y",
		outputPath,
	)

	// 执行并捕获错误
	if out, err := cmd.CombinedOutput(); err != nil {
		fmt.Printf("FFmpeg Error: %s\n", string(out))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "视频标准化处理失败"})
		return
	}

	// 6. 返回结果
	// 注意：outputPath 之后将作为参数通过 gRPC 发送给 Python 服务
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "视频处理完成，准备进入视觉识别阶段",
		"data": gin.H{
			"origin":    inputPath,
			"processed": outputPath,
		},
	})
}
