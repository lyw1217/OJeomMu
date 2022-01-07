package main

import (
	"log"
	"ojeommu/config"
	"ojeommu/controller"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/unrolled/secure"
)

func ServeStaticFiles(r *gin.Engine) {
	r.Static("/css", "./static/css")
	r.Static("/js", "./static/js")
	r.Static("/vendor", "./static/vendor")
	r.Static("/assets", "./assets")
	r.StaticFile("/favicon.ico", "./assets/favicon.ico")
	r.LoadHTMLGlob("templates/**/*")
}

func main() {
	secureFunc := func() gin.HandlerFunc {
		return func(c *gin.Context) {
			secureMiddleware := secure.New(secure.Options{
				SSLRedirect: true,
				SSLHost:     "mumeog.site:443",
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
	routeHttps := gin.Default()
	ServeStaticFiles(routeHttps)
	// Initialize the routes
	controller.InitRoutes(routeHttps)

	// HTTP
	go routeHttp.Run(":80")
	// HTTPS
	routeHttps.RunTLS(":8443", config.ServerCrt, config.ServerKey)

	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscanll.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

}
