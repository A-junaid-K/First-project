package main

import (
	"os"

	"github.com/first_project/database"
	loadenv "github.com/first_project/loadEnv"
	"github.com/first_project/routes"
	"github.com/gin-gonic/gin"
)

func init() {
	loadenv.Loadenv()
	database.Connetdb()
	database.SyncDatabase()
}

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	router := gin.New()
	router.Use(gin.Logger())

	routes.AuthRoutes(router)
	routes.UserRouter(router)

	router.Run(":" + port)
}
