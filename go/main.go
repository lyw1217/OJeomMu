package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HomePage : GET, "/"
func HomePage(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"templates/index.html",
		gin.H{},
	)
}

func main() {
	r := gin.Default()

	r.Static("/css", "./static/css")
	r.Static("/js", "./static/js")
	r.StaticFile("/favicon.ico", "./static/favicon.ico")

	r.LoadHTMLGlob("templates/*")
	//r.LoadHTMLFiles("templates/index.html")

	r.GET("/", HomePage)
	r.Run(":8090")
}

// https://dev.to/mizutani/how-to-build-web-app-with-go-gin-gonic-vue-3987 참고
