package controllers

import (
	"encoding/json"
	"moxin-zhicheng/internal/database"
	model "moxin-zhicheng/internal/models"
	"moxin-zhicheng/internal/redis"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

const DailyTagsKey = "chenxiang:StarTags"

func GetStarTags(c *gin.Context) {
	ctx := c.Request.Context()

	//  先看 Redis 里有没有缓存
	val, err := redis.Get(ctx, DailyTagsKey)
	if err == nil {
		var tags []string
		if json.Unmarshal([]byte(val), &tags) == nil {
			c.JSON(http.StatusOK, gin.H{"code": 200, "data": tags})
			return
		}
	}

	var starTags []string
	err = database.DB.Model(&model.PoetryTag{}).
		Order("RAND()"). // 让数据库随机排序
		Limit(10).       // 只取 10 条
		Pluck("tag", &starTags).Error
	if err != nil {
		starTags = []string{"清风", "明月", "江南", "边塞", "风雪", "前世", "蜀相", "老臣", "青史", "忠言"}
	}

	selected := starTags
	// 3. [关键] 计算到今天午夜的剩余时间，作为过期时间
	now := time.Now()
	// 今天凌晨 3 点
	target := time.Date(now.Year(), now.Month(), now.Day(), 3, 0, 0, 0, now.Location())

	// 如果当前已经过了今天的 3 点，目标就是明天的 3 点
	if now.After(target) {
		target = target.AddDate(0, 0, 1)
	}

	ttl := target.Sub(now)

	// 4. 存入 Redis
	data, _ := json.Marshal(selected)
	redis.Set(ctx, DailyTagsKey, data, ttl)

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": selected})
}

// SearchPoetry 综合搜索接口
func SearchPoetry(c *gin.Context) {
	keyword := c.Query("keyword")
	if keyword == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "搜索关键词不能为空"})
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	} // 限制单页最大 100 条，保护服务器
	offset := (page - 1) * pageSize

	db := database.DB
	var poems []model.Poetry
	var total int64
	// 1. 核心路径：查询标签表（此时标签表已包含：意象、完整标题、完整作者）
	var tag model.PoetryTag
	// 精确匹配标签名，利用 B-Tree 索引极速定位
	if err := db.Where("tag = ?", keyword).First(&tag).Error; err == nil {
		db.Model(&model.PoetryTagRelation{}).Where("tag_id = ?", tag.ID).Count(&total)
		// 2. 通过关联表获取所有相关的诗词 ID
		var poetryIDs []uint
		db.Model(&model.PoetryTagRelation{}).
			Where("tag_id = ?", tag.ID).
			Limit(pageSize).
			Offset(offset).
			Pluck("poetry_id", &poetryIDs)

		if len(poetryIDs) > 0 {
			// 3. 批量查询诗词详情，使用 ID 列表查询是数据库性能最优的操作
			db.Where("id IN ?", poetryIDs).Limit(20).Find(&poems)
		}
	}

	// 右模糊搜索作者和标题
	// 注意：这里去掉了 title 之前的 %，仅保留右模糊以利用索引最左匹配原则
	if len(poems) == 0 {
		query := db.Model(&model.Poetry{}).Where("title LIKE ? OR author LIKE ?", keyword+"%", keyword+"%")

		// 获取模糊匹配的总数
		query.Count(&total)

		// 获取分页数据
		err := query.Limit(pageSize).Offset(offset).Find(&poems).Error

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "搜索系统繁忙"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"list":  poems,
			"total": total,
			"page":  page,
			"size":  pageSize,
		},
	})
}
