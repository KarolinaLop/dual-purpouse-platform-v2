package web

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/KarolinaLop/dp/data"
	"github.com/KarolinaLop/dp/models"
	"github.com/gin-gonic/gin"
)

// ShowDashboard renders the dashboard page.
func ShowDashboard(c *gin.Context) {
	// call some data package func that loads all scans for the current user from the db
	user, ok := c.Value("user").(models.User) // type assertion and interfaces
	if !ok {
		err := errors.New("failed to find the user in this context")
		c.Error(err)
		return
	}

	scans, err := data.GetAllNmapScans(data.DB, user.ID)
	if err != nil {
		err = fmt.Errorf("failed to retrieve scans: %w", err)
		c.Error(err)
		return
	}

	c.HTML(http.StatusOK, "dashboard.html", gin.H{
		"userName": c.Value("user").(models.User).Name,
		"scans":    scans,
	})

}
