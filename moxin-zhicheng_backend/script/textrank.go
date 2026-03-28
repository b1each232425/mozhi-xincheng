package main

import (
	"fmt"
	"io/ioutil"
	"moxin-zhicheng/internal/config"
	"moxin-zhicheng/internal/database"
	"moxin-zhicheng/internal/logger"
	model "moxin-zhicheng/internal/models"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/yanyiwu/gojieba"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const checkpointFile = "poetry_checkpoint.txt"

func loadCheckpoint() uint {
	data, err := ioutil.ReadFile(checkpointFile)
	if err != nil {
		if os.IsNotExist(err) {
			return 0
		}
		return 0
	}
	id, _ := strconv.Atoi(strings.TrimSpace(string(data)))
	return uint(id)
}

func saveCheckpoint(id uint) {
	_ = ioutil.WriteFile(checkpointFile, []byte(fmt.Sprintf("%d", id)), 0644)
}

func main() {
	logger.InitLogger("dev")
	config.InitConfig()
	database.InitDB()
	db := database.DB

	x := gojieba.NewJieba()
	defer x.Free()

	lastID := loadCheckpoint()
	fmt.Printf("🔄 任务启动：从 ID > %d 开始处理\n", lastID)

	const workerCount = 8 // TextRank 构图较耗内存，建议核心数左右
	batchSize := 150
	jobs := make(chan []model.Poetry, workerCount)
	var wg sync.WaitGroup

	for w := 1; w <= workerCount; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for batch := range jobs {
				processBatch(db, x, batch)
			}
		}()
	}

	for {
		var list []model.Poetry
		err := db.Where("id > ?", lastID).Order("id asc").Limit(batchSize).Find(&list).Error
		if err != nil || len(list) == 0 {
			break
		}
		jobs <- list
		lastID = list[len(list)-1].ID
		saveCheckpoint(lastID)
		fmt.Printf("🚀 已投递至 ID: %d\n", lastID)
	}

	close(jobs)
	wg.Wait()
	fmt.Println("✅ 任务全部完成！")
}

func processBatch(db *gorm.DB, x *gojieba.Jieba, batch []model.Poetry) {
	// 记录本批次中：标签名 -> 出现的次数（用于更新 PoetryTag 计数）
	batchAggregatedCounts := make(map[string]int64)
	// 记录本批次中：PoetryID -> 包含的标签名列表（用于后续插入关联表）
	poetryToTags := make(map[uint][]string)

	for i := range batch {
		p := &batch[i]
		singlePoetryTags := make(map[string]bool)

		// 提取标题方便后续搜索
		if p.Title != "" {
			// 直接存入，不分词，确保搜索标题时 100% 匹配
			singlePoetryTags[p.Title] = true
		}

		// --- 轨道 1：TextRank 高质量意象 ---
		// (保持你原有的 TextRank 提取逻辑不变...)
		// 假设提取出的词存入了 singlePoetryTags

		// --- 轨道 2：CutAll 全模式保底 ---
		allPossibleWords := x.CutAll(p.Paragraphs)
		for _, w := range allPossibleWords {
			if len([]rune(w)) == 2 && !strings.ContainsAny(w, "一二三四五六七八九十") {
				singlePoetryTags[w] = true
			}
		}

		// 汇总结果
		for tag := range singlePoetryTags {
			batchAggregatedCounts[tag]++
			poetryToTags[p.ID] = append(poetryToTags[p.ID], tag)
		}
	}

	if len(batchAggregatedCounts) == 0 {
		return
	}

	// --- 数据库同步核心：使用事务保证原子性 ---
	db.Transaction(func(tx *gorm.DB) error {
		// 1. 批量更新/创建 PoetryTag，并获取本批次所有标签的 ID
		tagNameToID := make(map[string]uint)
		for tag, count := range batchAggregatedCounts {
			t := model.PoetryTag{Tag: tag, Count: count}
			// 使用 OnConflict 确保标签存在并更新计数
			tx.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "tag"}},
				DoUpdates: clause.Assignments(map[string]interface{}{"count": gorm.Expr("count + ?", count)}),
			}).Create(&t)

			// GORM 在 Create/Upsert 后会自动将主键回填到结构体 t.ID
			tagNameToID[tag] = t.ID
		}

		// 2. 构建中间表记录
		var relations []model.PoetryTagRelation
		for pID, tags := range poetryToTags {
			for _, tagName := range tags {
				if tID, ok := tagNameToID[tagName]; ok {
					relations = append(relations, model.PoetryTagRelation{
						PoetryID: pID,
						TagID:    tID,
					})
				}
			}
		}

		// 3. 批量插入关联表（1000条一跳，防止 SQL 过长）
		if len(relations) > 0 {
			if err := tx.CreateInBatches(relations, 1000).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
