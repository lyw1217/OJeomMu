package controller

import (
	b64 "encoding/base64"
	"fmt"
	"net/http"
	"ojeommu/config"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"

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

	red_url := c.Query("redirect")

	if len(red_url) > 0 {
		c.Redirect(
			http.StatusMovedPermanently,
			red_url,
		)
		return
	}

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
			c.JSON(http.StatusNotFound, gin.H{
				"status": http.StatusNotFound,
				"reason": "Not Found",
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

func checkAuth(key string) bool {
	sDec, _ := b64.StdEncoding.DecodeString(key)
	if strings.Compare(strings.Trim(string(sDec), " "), config.Keys.Newyo.Apikey) == 0 {
		return true
	} else {
		return false
	}
}

func SearchBotHandler(c *gin.Context) {

	/*
		auth := c.Query("auth")
		if !checkAuth(auth) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status": http.StatusUnauthorized,
				"reason": "Unauthorized API Key",
			})
			return
		}
	*/

	auth := c.Query("auth")
	if len(auth) > 0 {
		config.Keys.Kakao.Rest = auth
	}

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

		tmp, err := GetSearchKeyword(p, 500)
		if err != nil {
			log.Println("Error, Failed GetSearchKeyword()")
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": http.StatusInternalServerError,
				"reason": "Internal Server Error",
			})
			return
		}

		cat := c.Query("cat")
		if len(cat) > 0 {
			jsonData.Category = cat
		} else {
			jsonData.Category = "anything"
		}
		jsonData.Radius = "0.5"
		if len(tmp) > 0 {
			jsonData.X = tmp[0].X
			jsonData.Y = tmp[0].Y
			qry_result.PlaceName = tmp[0].PlaceName
		} else {
			log.Println("Not Found Keyword Place")
			c.JSON(http.StatusNotFound, gin.H{
				"status": http.StatusNotFound,
				"reason": fmt.Sprintf("Query(%s) Place Not Found", qry),
			})
			return
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
		//log.Println("matched_place =", matched_place)
		if matched_place == nil {
			c.JSON(http.StatusNotFound, gin.H{
				"status": http.StatusNotFound,
				"reason": fmt.Sprintf("Not Found Restaurant 500m around (query(%s), searchd_place(%s))", qry, qry_result.PlaceName),
			})
			return
		} else {
			// 현재 위치와 place간 거리 구하기
			d := GetDistance(jsonData.X, jsonData.Y, matched_place.X, matched_place.Y)

			c.JSON(http.StatusOK, gin.H{
				"hdr":     matched_place.PlaceName,
				"place":   qry_result.PlaceName,
				"d":       d,
				"lnk":     matched_place.PlaceUrl,
				"cat":     matched_place.CategoryName,
				"address": matched_place.AddressName,
				"road":    matched_place.RoadAddressName,
				"phone":   matched_place.Phone,
				"x":       matched_place.X,
				"y":       matched_place.Y,
			})
			return
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"reason": "Bad Request",
		})
		return
	}
}

func WtImgHandler(c *gin.Context) {
	auth := c.Query("auth")
	if !checkAuth(auth) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"reason": "Unauthorized API Key",
		})
		return
	}
	qry := c.Query("query")
	if len(qry) > 0 {
		img_path, _ := filepath.Abs(fmt.Sprintf("./assets/img/wt/wt%s.png", qry))
		fmt.Println(img_path)
		if _, err := os.Stat(img_path); os.IsNotExist(err) {
			c.JSON(http.StatusNotFound, gin.H{
				"status": http.StatusNotFound,
				"reason": "Not Found",
			})
			return
		}

		c.File(img_path)
	}
}

func InitRoutes(r *gin.Engine) {

	r.NoRoute(NotFoundPage)

	r.GET("/", HomePage)
	r.GET("/index.html", HomePage)
	r.GET("/info.html", InfoPage)
	r.GET("/test.html", TestPage)

	r.POST("/sendToGo", SearchHandler)
	r.GET("/api", SearchBotHandler)
	r.GET("/weather", WtImgHandler)
}
