package controller

import (
	"log"
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
		list, cond, err := SearchKeyword(jsonData)
		if err != nil {
			log.Println(err)
			return
		}
		matched_place := GetCondPlace(list, cond)
		log.Println("matched_place =", matched_place)
		if matched_place == nil {
			log.Println("Error, failed GetCondPlace()")
		} else {
			// 데이터 전송
			c.JSON(200, gin.H{
				"ID":           matched_place.Id,
				"NAME":         matched_place.PlaceName,
				"PHONE":        matched_place.Phone,
				"ADDRESS":      matched_place.AddressName,
				"ROAD_ADDRESS": matched_place.RoadAddressName,
				"X":            matched_place.X,
				"Y":            matched_place.Y,
				"URL":          matched_place.PlaceUrl,
				"DISTANCE":     matched_place.Distance,
			})
		}
	} else {
		// handle error
		log.Println("Error, failed BindJSON()")
		return
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
