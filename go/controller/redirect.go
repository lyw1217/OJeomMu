package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RedirectMaeKyung : GET, "/maekyung"
func redirectMaeKyung(c *gin.Context) {
	c.Redirect(
		http.StatusMovedPermanently,
		MkMSGUrl,
	)
}

// RedirectHanKyung : GET, "/hankyung"
func redirectHanKyung(c *gin.Context) {
	c.Redirect(
		http.StatusMovedPermanently,
		HkIssueTodayUrl,
	)
}
