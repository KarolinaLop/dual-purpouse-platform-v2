package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// StartScanHandler starts an nmap scan and saves the results to the database.
func StartScanHandler(c *gin.Context) {
	// TODO: start scan
	// TODO: save to DB

	c.JSON(http.StatusOK, nil)
}
