package web

import (
	"net/http"
	"os/exec"

	"github.com/gin-gonic/gin"
)

// StartScanHandler starts an nmap scan and saves the results to the database. Edited a comment to test git commit.
func StartScanHandler(c *gin.Context) {
	// TODO: start scan
	exec.Cmd()
	// TODO: save to DB

	c.JSON(http.StatusOK, nil)
}
