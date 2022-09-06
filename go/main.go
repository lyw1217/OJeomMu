package main

import (
	"ojeommu/config"
	"ojeommu/controller"
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
	"github.com/unrolled/secure"
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
	secureFunc := func() gin.HandlerFunc {
		return func(c *gin.Context) {
			secureMiddleware := secure.New(secure.Options{
				SSLRedirect: true,
				SSLHost:     "mumeog.site",
			})
			err := secureMiddleware.Process(c.Writer, c.Request)

			// If there was an error, do not continue.
			if err != nil {
				return
			}

			c.Next()
		}
	}()

	routeHttp := gin.Default()
	routeHttp.Use(secureFunc)
	ServeStaticFiles(routeHttp)
	// Initialize the routes
	controller.InitRoutes(routeHttp)

	// HTTP
	go routeHttp.Run(":80")
	// HTTPS
	routeHttp.RunTLS(":8443", config.ServerCrt, config.ServerKey)

	quit := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	// kill (no param) default send syscanll.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-quit
		log.Println("RECEIVE SIG : ", sig)
		done <- true
	}()
	<-done
	log.Println("Shutdown Server ...")
}
