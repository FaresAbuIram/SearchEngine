package routes

import (
	"searchEngine/controllers"
	"searchEngine/docs"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Setup(router *gin.Engine) {
	docs.SwaggerInfo.Title = "Swagger Example API"
	docs.SwaggerInfo.Description = "This is a sample server Petstore server."
	docs.SwaggerInfo.Version = "2.0"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.Schemes = []string{"http"}

	router.Use(static.Serve("/", static.LocalFile("./website/dist", true)))
	router.POST("/createNewResource", controllers.Createresource)
	// router.POST("/get/:token", controllers.GetToken)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
