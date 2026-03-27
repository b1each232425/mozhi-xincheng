package main

import (
	"encoding/json"
	"io/ioutil"
	"log/slog"
	"moxin-zhicheng/internal/config"
	"moxin-zhicheng/internal/database"
	"moxin-zhicheng/internal/logger"
	"moxin-zhicheng/internal/models"
	"path/filepath"
	"strings"
)

type ImportStats struct {
	SuccessFiles int
	FailedFiles  int
	TotalRecords int
	SkippedFiles []string // 记录具体哪些文件没被识别成功
}

var stats = ImportStats{}

type RawPoetry struct {
	Title      string   `json:"title"`
	Author     string   `json:"author"`
	Paragraphs []string `json:"paragraphs"`
	Content    []string `json:"content"`
	Para       []string `json:"para"`
	Chapter    string   `json:"chapter"`
	Rhythmic   string   `json:"rhythmic"`
}

func main() {
	logger.InitLogger("dev")
	config.InitConfig()
	database.InitDB()

	// 1. 定义需要扫描的目录及其对应的类型标签
	// 这里的路径要对应你本地 clone 下来的文件夹名称
	features := map[string]string{
		"tang":        "tang",
		"song":        "song",
		"chuci":       "chuci",
		"caocao":      "caocao",
		"shijing":     "shijing",
		"huajianji":   "huajianji",
		"lunyu":       "lunyu",
		"wudai":       "wudai",
		"nalanxingde": "nalanxingde",
	}

	pattern := filepath.Join("script", "chinese_poetry_Data", "*.json")
	files, _ := filepath.Glob(pattern)

	if len(files) == 0 {
		logger.Warn("未发现匹配文件", slog.String("path", pattern))
		return
	}

	for _, file := range files {
		// 跳过元数据文件
		if strings.Contains(file, "authors") {
			continue
		}

		fileName := filepath.Base(file) // 获取纯文件名，如 poet.tang.0.json
		pType := "unknown"              // 默认标记为未知

		// 3. 核心：遍历映射表，根据文件名特征匹配标签
		for key, label := range features {
			if strings.Contains(fileName, key) {
				pType = label
				break
			}
		}

		// 4. 监控：如果识别失败（即没有匹配到任何关键字），打印警告
		if pType == "unknown" {
			logger.Warn("🚩 无法识别朝代的文件", slog.String("file", fileName))
		}

		logger.Info("正在导入", slog.String("file", fileName), slog.String("type", pType))
		importFile(file, pType)
	}

	logger.Info("🎉 数据库灌顶完成！")
}

func importFile(filePath string, pType string) {
	fileName := filepath.Base(filePath)
	if strings.Contains(fileName, ".tang.") {
		pType = "tang"
	} else if strings.Contains(fileName, ".song.") {
		pType = "song"
	} else if strings.Contains(fileName, "huajianji.") {
		pType = "huajianji"
	}

	byteValue, err := ioutil.ReadFile(filePath)
	if err != nil {
		logger.Error("读取文件失败", err, slog.String("file", fileName))
		stats.FailedFiles++
		stats.SkippedFiles = append(stats.SkippedFiles, filePath+" (读取失败)")
		return
	}

	var raws []RawPoetry
	if err := json.Unmarshal(byteValue, &raws); err != nil {
		logger.Error("解析JSON失败", err, slog.String("file", fileName))
		stats.FailedFiles++
		stats.SkippedFiles = append(stats.SkippedFiles, filePath+" (JSON格式错误)")
		return
	}

	var poetries []model.Poetry
	for _, r := range raws {
		var poetryContent []string
		if len(r.Paragraphs) > 0 {
			poetryContent = r.Paragraphs // 优先取 paragraphs
		} else if len(r.Content) > 0 {
			poetryContent = r.Content // 备选取 content
		} else if len(r.Para) > 0 {
			poetryContent = r.Para
		}

		if r.Author == "" && strings.Contains(filePath, "曹操") {
			r.Author = "曹操"
		}
		formattedText := strings.Join(poetryContent, "\n")
		poetries = append(poetries, model.Poetry{
			Title:      r.Title,
			Author:     r.Author,
			Paragraphs: formattedText,
			Type:       pType,
			Chapter:    r.Chapter,
			Rhythmic:   r.Rhythmic,
		})
	}
	if len(poetries) > 0 {
		if err := database.DB.CreateInBatches(poetries, 500).Error; err != nil {
			logger.Error("入库失败", err)
			stats.FailedFiles++
		} else {
			stats.SuccessFiles++
			stats.TotalRecords += len(poetries)
			logger.Info("导入成功", slog.String("file", fileName), slog.Int("count", len(poetries)))
		}
	} else {
		// 如果一个文件解析完竟然一条有效数据都没有，记录下来
		stats.SkippedFiles = append(stats.SkippedFiles, filePath+" (无有效内容)")
	}

	// 依然采用分批插入，每批 500 条
	if err := database.DB.CreateInBatches(poetries, 500).Error; err != nil {
		logger.Error("入库失败", err)
	}
}
