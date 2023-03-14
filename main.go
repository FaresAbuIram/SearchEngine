package main

import (
	"searchEngine/database"
	"searchEngine/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {

	database.Connect()

	router := gin.Default()

	router.Use(cors.Default())
	routes.Setup(router)
	router.Run()
}
