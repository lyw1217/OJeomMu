package main

import (
	"ojeommu/controller"

	"github.com/gin-gonic/gin"
)

var r *gin.Engine

func main() {
	r = gin.Default()

	r.Static("/css", "./static/css")
	r.Static("/js", "./static/js")
	r.Static("/vendor", "./static/vendor")
	r.Static("/assets", "./assets")
	r.StaticFile("/favicon.ico", "./assets/favicon.ico")

	r.LoadHTMLGlob("templates/**/*")
	// Initialize the routes
	controller.InitRoutes(r)

	r.Run(":8090")
}

// https://dev.to/mizutani/how-to-build-web-app-with-go-gin-gonic-vue-3987 참고
