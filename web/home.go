package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HomeHandler renders the home page.
func HomeHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "home.html", gin.H{})
}
