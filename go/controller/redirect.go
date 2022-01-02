package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RedirectMaeKyung : GET, "/maekyung"
func RedirectMaeKyung(c *gin.Context) {
	c.Redirect(
		http.StatusMovedPermanently,
		MkMSGUrl,
	)
}

// RedirectHanKyung : GET, "/hankyung"
func RedirectHanKyung(c *gin.Context) {
	c.Redirect(
		http.StatusMovedPermanently,
		HkIssueTodayUrl,
	)
}
