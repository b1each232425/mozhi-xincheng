package routes

import (
	"moxin-zhicheng/internal/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	chenxiang := r.Group("/chenxiang")
	{
		chenxiang.GET("/list",
			controllers.SearchPoetry)
		chenxiang.GET("/starTags",
			controllers.GetStarTags)
		chenxiang.GET("/singlePoem/:id",
			controllers.SinglePoem)
	}

	luobi := r.Group("/luobi")
	{
		// 文章相关
		luobi.POST("/article", controllers.CreateArticle)
		luobi.GET("/articles", controllers.GetArticleList)
		luobi.GET("/article/:id", controllers.GetArticleDetail)
		// 图片上传
		luobi.POST("/upload/image", controllers.UploadArticleImage)
	}
}
