package controller

import (
	"fmt"
	"log"
	"net/http"
	"ojeommu/config"

	"github.com/gin-gonic/gin"
)

const TITLE_NAME = "오늘 무먹?"

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

// InfoPage : GET, "/"
func InfoPage(c *gin.Context) {

	c.HTML(
		http.StatusOK,
		"views/info.html",
		gin.H{
			"title": TITLE_NAME,
		},
	)
}

// TestPage : GET, "/"
func TestPage(c *gin.Context) {

	c.HTML(
		http.StatusOK,
		"views/test.html",
		gin.H{
			"title": TITLE_NAME,
		},
	)
}

func SearchHandler(c *gin.Context) {
	var jsonData SearchCond_t
	var qry_result KeywordDocuments_t
	if c.BindJSON(&jsonData) == nil {
		qry := c.Query("query")
		if len(qry) > 0 {

			var p = KeywordParam_t{
				Query: qry,
				Page:  1,
				Size:  15,
				Sort:  "accuracy",
			}
			fmt.Println(p)
			tmp, err := GetSearchKeyword(p, 500)
			if err != nil {
				log.Println("Error, Failed GetSearchKeyword()")
				c.JSON(http.StatusInternalServerError, gin.H{
					"status": http.StatusInternalServerError,
					"reason": "Internal Server Error",
				})
				return
			}
			jsonData.Category = "FD6"
			jsonData.Radius = "0.5"
			if len(tmp) > 0 {
				jsonData.X = tmp[0].X
				jsonData.Y = tmp[0].Y
				qry_result.PlaceName = tmp[0].PlaceName
			} else {
				c.JSON(http.StatusNotFound, gin.H{
					"status": http.StatusNotFound,
					"reason": "Not Found",
				})
				return
			}
		}
		matched_place, total_nums, err := RectSearch(jsonData)

		if err != nil {
			log.Println("Error, failed RectSearch()")
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": http.StatusInternalServerError,
				"reason": "Internal Server Error",
			})
			return
		}
		log.Println("matched_place =", matched_place)
		if matched_place == nil {
			log.Println("Error, failed GetCondPlace()")
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": http.StatusInternalServerError,
				"reason": "Internal Server Error",
			})
		} else {
			// 현재 위치와 place간 거리 구하기
			d := GetDistance(jsonData.X, jsonData.Y, matched_place.X, matched_place.Y)

			// 데이터 전송
			c.JSON(200, gin.H{
				"ID":             matched_place.Id,
				"NAME":           matched_place.PlaceName,
				"CAT_NAME":       matched_place.CategoryName,
				"PHONE":          matched_place.Phone,
				"ADDRESS":        matched_place.AddressName,
				"ROAD_ADDRESS":   matched_place.RoadAddressName,
				"X":              matched_place.X,
				"Y":              matched_place.Y,
				"URL":            matched_place.PlaceUrl,
				"DISTANCE":       d,
				"TOTAL_NUMS":     total_nums,
				"QRY_PLACE_NAME": qry_result.PlaceName,
			})
		}

	} else {
		// handle error
		log.Println("Error, failed BindJSON()")
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"reason": "Internal Server Error",
		})
		return
	}
}

func SearchBotHandler(c *gin.Context) {
	var jsonData SearchCond_t
	var qry_result KeywordDocuments_t

	qry := c.Query("query")
	if len(qry) > 0 {

		var p = KeywordParam_t{
			Query: qry,
			Page:  1,
			Size:  15,
			Sort:  "accuracy",
		}
		fmt.Println(p)
		tmp, err := GetSearchKeyword(p, 500)
		if err != nil {
			log.Println("Error, Failed GetSearchKeyword()")
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": http.StatusInternalServerError,
				"reason": "Internal Server Error",
			})
			return
		}
		jsonData.Category = "anything"
		jsonData.Radius = "0.5"
		if len(tmp) > 0 {
			jsonData.X = tmp[0].X
			jsonData.Y = tmp[0].Y
			qry_result.PlaceName = tmp[0].PlaceName
		} else {
			c.JSON(http.StatusNotFound, gin.H{
				"status": http.StatusNotFound,
				"reason": "Not Found",
			})
			return
		}
	}
	matched_place, _, err := RectSearch(jsonData)

	if err != nil {
		log.Println("Error, failed RectSearch()")
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"reason": "Internal Server Error",
		})
		return
	}
	log.Println("matched_place =", matched_place)
	if matched_place == nil {
		log.Println("Error, failed GetCondPlace()")
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"reason": "Internal Server Error",
		})
		return
	} else {
		// 현재 위치와 place간 거리 구하기
		d := GetDistance(jsonData.X, jsonData.Y, matched_place.X, matched_place.Y)

		c.String(http.StatusOK, "오늘 점심은 '%s' 어떠세요?\n------------------------------------------\n %s (으)로부터 거리는 %dm\n URL : %s", matched_place.PlaceName, qry_result.PlaceName, d, matched_place.PlaceUrl)
		return
	}
}

func InitRoutes(r *gin.Engine) {

	r.NoRoute(NotFoundPage)

	r.GET("/", HomePage)
	r.GET("/index.html", HomePage)
	r.GET("/info.html", InfoPage)
	r.GET("/test.html", TestPage)

	r.POST("/sendToGo", SearchHandler)
	r.POST("/ojeommu", SearchBotHandler)

	/* Redirect, for scraping-news-go */
	r.GET("/maekyung", RedirectMaeKyung)
	r.GET("/hankyung", RedirectHanKyung)
	r.GET("/quicknews", RedirectQuickNews)
}
