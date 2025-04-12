package main

import (
	"log"

	"github.com/KarolinaLop/dp/web"
	"github.com/gin-gonic/gin"
)

func main() {
	s := web.SetupServer()

	log.Println("Server is running on http://localhost:" + web.PORT)
	log.Fatal(s.ListenAndServe())

	r := gin.Default()

	// Routes for scanning
	r.POST("/api/scan/start", web.StartScanHandler)
	r.POST("/api/scan/stop", web.StopScanHandler)

	// Start the server
	r.Run(":8080")
}
