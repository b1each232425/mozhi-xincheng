package controllers

import (
	"moxin-zhicheng/internal/database"
	model "moxin-zhicheng/internal/models"
	"moxin-zhicheng/internal/response"
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
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "未登录")
		return
	}
	var ArcReq ArticleReq
	if err := c.ShouldBindJSON(&ArcReq); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	article := model.Article{
		Title:    ArcReq.Title,
		Content:  ArcReq.Content,
		Summary:  ArcReq.Summary,
		AuthorID: userID.(uint),
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
		response.BadRequest(c, "内容过长，请控制在1万字内")
		return
	}

	if err := database.DB.Create(&article).Error; err != nil {
		response.ServerError(c, "保存失败")
		return
	}

	response.SuccessWithMsg(c, "落笔成功", article)
}

// GetArticleList 获取文章列表
func GetArticleList(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "未登录")
		return
	}
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
	err := query.Order("created_at DESC").
		Where("AuthorID=?", userID).
		Limit(pageSize).
		Offset(offset).
		Find(&articles).
		Error
	if err != nil {
		response.ServerError(c, "查询失败")
		return
	}

	response.Success(c, gin.H{
		"list":  articles,
		"total": total,
	})
}

// GetArticleDetail 获取单篇文章详情
func GetArticleDetail(c *gin.Context) {
	id := c.Param("id")
	var article model.Article

	// 详情页才加载 Content 字段
	if err := database.DB.First(&article, id).Error; err != nil {
		response.NotFound(c, "文章已随风散去")
		return
	}

	// 增加阅读量（异步增加更好，这里演示简单写法）
	database.DB.Model(&article).UpdateColumn("view_count", article.ViewCount+1)

	response.Success(c, article)
}

func SwitchPublic(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "未登录")
		return
	}
	id := c.Param("id")
	var article model.Article
	if err := database.DB.Where("id = ? AND author_id = ?", id, userID).First(&article).Error; err != nil {
		response.NotFound(c, "文章不存在或无权限")
		return
	}

	newStatus := article.IsPublic ^ 1
	database.DB.Model(&article).UpdateColumn("is_public", newStatus)
	response.SuccessWithMsg(c, "切换成功", gin.H{"is_public": newStatus})
}
