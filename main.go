package main

import (
	log "github.com/sirupsen/logrus"
	"ojeommu/controller"
	"os"

	"github.com/gin-gonic/gin"
)

func ServeStaticFiles(r *gin.Engine) {
	r.Static("/css", "./static/css")
	r.Static("/js", "./static/js")
	r.Static("/assets", "./assets")
	r.StaticFile("/favicon.ico", "./assets/favicon.ico")
	r.StaticFile("/robots.txt", "./static/robots.txt")
	r.LoadHTMLGlob("templates/**/*")
}

func main() {
	routeHttp := gin.Default()
	routeHttp.Use(gin.Logger())
	ServeStaticFiles(routeHttp)
	// Initialize the routes
	controller.InitRoutes(routeHttp)

	// HTTP
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}
	routeHttp.Run(":" + port)
}
