package main

import (
	"ojeommu/config"
	"ojeommu/controller"
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"
	"github.com/unrolled/secure"

	"github.com/gin-gonic/gin"
)

func main() {

	config.Keys = config.LoadKeysConfig()
	config.SetupLogger()

	secureFunc := func() gin.HandlerFunc {
		return func(c *gin.Context) {
			secureMiddleware := secure.New(secure.Options{
				SSLRedirect: true,
				SSLHost:     "lyw1217.synology.me",
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
	controller.ServeStaticFiles(routeHttp)
	// Initialize the routes
	controller.InitRoutes(routeHttp)

	// HTTP
	go routeHttp.Run(":80")
	// HTTPS
	routeHttp.RunTLS(":443", config.ServerCrt, config.ServerKey)

	/*
		routeHttp := gin.Default()
		controller.ServeStaticFiles(routeHttp)
		// Initialize the routes
		controller.InitRoutes(routeHttp)

		log.Fatal(autotls.Run(routeHttp, "mumeog.site", "lunchtoday.site", "오점무.site"))
	*/
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
