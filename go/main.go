package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"ojeommu/controller"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/gin-gonic/autotls"
	"github.com/gin-gonic/gin"
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
	r := gin.Default()
	ServeStaticFiles(r)
	// Initialize the routes
	controller.InitRoutes(r)

	/* https://github.com/gin-gonic/gin#graceful-shutdown-or-restart */
	srv := &http.Server{
		Addr:    ":8090",
		Handler: r,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Printf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so don't need to add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
