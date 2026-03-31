package controllers

import (
	"moxin-zhicheng/internal/database"
	model "moxin-zhicheng/internal/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/microcosm-cc/bluemonday"
)

type ArticleReq struct {
	Title    string `json:"title"`
	Content  string `json:"content"`
	Summary  string `json:"summary"`
	IsPublic bool   `json:"is_public"`
}

// CreateArticle 发表文章/日记
func CreateArticle(c *gin.Context) {
	var ArcReq ArticleReq
	if err := c.ShouldBindJSON(&ArcReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	article := model.Article{
		Title:   ArcReq.Title,
		Content: ArcReq.Content,
		Summary: ArcReq.Summary,
		// Season 字段在这里被“过滤”掉了，因为数据库模型里没有它
	}

	// 允许基本的 HTML 标签（p, br, img, div, span 等），剔除 <script>
	p := bluemonday.UGCPolicy()
	p.AllowAttrs("src", "alt").OnElements("img") // 允许图片展示
	article.Content = p.Sanitize(article.Content)

	//if article.Summary == "" {
	//	// 简单处理：过滤 HTML 标签后取前 150 字
	//	plainText := p.Sanitize(article.Content) // 实际上可以用更彻底的去标签库
	//	runes := []rune(plainText)
	//	if len(runes) > 50 {
	//		article.Summary = string(runes[:50]) + "..."
	//	} else {
	//		article.Summary = string(runes)
	//	}
	//}

	// 强制字数检查（安全防御）
	if len(article.Content) > 30000 { // 字节长度，粗略对应1万汉字+HTML
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "内容过长，请控制在1万字内"})
		return
	}

	if err := database.DB.Create(&article).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "保存失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": article, "msg": "落笔成功"})
}

// GetArticleList 获取文章列表（高性能分页版）
func GetArticleList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	offset := (page - 1) * pageSize

	var articles []model.Article
	var total int64

	// 关键优化：使用 .Omit("content") 排除大文本字段
	// 这样查询 10 条数据只会占用极小的内存，不会因为 Content 字段撑爆服务器
	query := database.DB.Model(&model.Article{}).Omit("content")

	// 统计总数
	query.Count(&total)

	// 分页查询
	err := query.Order("created_at DESC").Limit(pageSize).Offset(offset).Find(&articles).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "查询失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"list":  articles,
			"total": total,
		},
	})
}

// GetArticleDetail 获取单篇文章详情
func GetArticleDetail(c *gin.Context) {
	id := c.Param("id")
	var article model.Article

	// 详情页才加载 Content 字段
	if err := database.DB.First(&article, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "msg": "文章已随风散去"})
		return
	}

	// 增加阅读量（异步增加更好，这里演示简单写法）
	database.DB.Model(&article).UpdateColumn("view_count", article.ViewCount+1)

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": article})
}
