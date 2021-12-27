package controller

import (
	"net/http"
	"ojeommu/config"

	"github.com/gin-gonic/gin"
)

// NotFoundPage : NoRoute
func notFoundPage(c *gin.Context) {
	c.HTML(
		http.StatusNotFound,
		"views/404.html",
		gin.H{},
	)
}

// HomePage : GET, "/"
func homePage(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"views/index.html",
		gin.H{
			"key": config.Keys.Kakao.JS,
		},
	)
}

// HomePage : GET, "/"
// https://startbootstrap.com/template/simple-sidebar
func homePage2(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"views/index2.html",
		gin.H{
			"key": config.Keys.Kakao.JS,
		},
	)
}

func InitRoutes(r *gin.Engine) {

	r.NoRoute(notFoundPage)

	r.GET("/", homePage)
	r.GET("/index.html", homePage)

	r.GET("/index2.html", homePage2)

	/* Redirect, for scraping-news-go */
	r.GET("/maekyung", redirectMaeKyung)
	r.GET("/hankyung", redirectHanKyung)
}
