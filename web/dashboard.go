package web

import (
	"net/http"

	"github.com/KarolinaLop/dp/models"
	"github.com/gin-gonic/gin"
)

// ShowDashboard renders the dashboard page.
func ShowDashboard(c *gin.Context) {
	c.HTML(http.StatusOK, "dashboard.html", gin.H{
		"userName": c.Value("user").(models.User).Name,
	})
}
