package main

import (
	"fmt"
	"moxin-zhicheng/internal/config"
	"moxin-zhicheng/internal/database"
	"moxin-zhicheng/internal/logger"
	model "moxin-zhicheng/internal/models"
	"moxin-zhicheng/utils"
	"time"

	"gorm.io/gorm"
)

func main() {
	config.InitConfig()
	logger.InitLogger("dev")
	database.InitDB()
	utils.InitQdrant()
	// 1. 使用 FindInBatches 分批处理，防止 30W 数据一次性撑爆 16GB 内存
	batchSize := 5
	err := database.DB.Model(&model.Poetry{}).
		Where("translation IS NULL OR translation = ''").
		FindInBatches(&[]model.Poetry{}, batchSize, func(tx *gorm.DB, batch int) error {
			fmt.Printf("--- 正在处理第 %d 批数据 ) ---\n", batch)

			for _, p := range *tx.Statement.Dest.(*[]model.Poetry) {
				// 步骤 1: 本地 RAG 检索 (E 盘 Qdrant)
				references := utils.SearchSimilarCorpus(p.Paragraphs)

				// 步骤 2: 本地 AI 推理 (Ollama)
				// 建议在执行前运行 $env:OLLAMA_MAX_VRAM=0 强制内存运行
				aiRes := utils.GenerateWithRAG(p, references)

				// 步骤 3: 回填并标记
				p.Translation = "[AI补全] " + aiRes.Translation
				p.Annotation = aiRes.Annotation
				database.DB.Save(&p)

				// 给显卡喘息时间，防止 4050 持续满载高温
				time.Sleep(100 * time.Millisecond)
			}
			return nil
		}).Error

	if err != nil {
		fmt.Println("任务执行异常:", err)
	}
}
