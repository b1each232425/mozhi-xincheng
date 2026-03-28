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
	}
}
