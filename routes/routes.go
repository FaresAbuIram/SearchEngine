package routes

import (
	"searchEngine/controllers"
	"searchEngine/database"
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

	resourcedatabase := database.NewResourceService()
	resourceService := controllers.NewResourceService(resourcedatabase)

	router.Use(static.Serve("/", static.LocalFile("./website/dist", true)))
	router.POST("/createNewResource", resourceService.CreateResource)
	router.POST("/search", resourceService.Search)
	router.GET("/resource/:id", resourceService.GetResource)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
