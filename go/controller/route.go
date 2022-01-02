package controller

import (
	"fmt"
	"net/http"
	"ojeommu/config"

	"github.com/gin-gonic/gin"
)

const TITLE_NAME = "오점무"

// NotFoundPage : NoRoute
func NotFoundPage(c *gin.Context) {
	c.HTML(
		http.StatusNotFound,
		"views/404.html",
		gin.H{"title": TITLE_NAME},
	)
}

// HomePage : GET, "/"
// https://startbootstrap.com/template/simple-sidebar
func HomePage(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"views/index.html",
		gin.H{
			"title": TITLE_NAME,
			"key":   config.Keys.Kakao.JS,
		},
	)
}

func SearchHandler(c *gin.Context) {
	var jsonData SearchCond_t
	if c.BindJSON(&jsonData) == nil {
		fmt.Println(SearchKeyword(jsonData))
	} else {
		// handle error
		fmt.Println("ERROR")
	}
}

func InitRoutes(r *gin.Engine) {

	r.NoRoute(NotFoundPage)

	r.GET("/", HomePage)
	r.GET("/index.html", HomePage)

	r.POST("/sendToGo", SearchHandler)

	/* Redirect, for scraping-news-go */
	r.GET("/maekyung", RedirectMaeKyung)
	r.GET("/hankyung", RedirectHanKyung)
}
