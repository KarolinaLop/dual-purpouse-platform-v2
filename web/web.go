package web

import (
	"html/template"
	"log"
	"net/http"

	"github.com/KarolinaLop/dp/frontend"

	"github.com/gin-gonic/gin"
)

const (
	assetsDir    = "frontend/assets/"
	templatesDir = "frontend/templates/**/*.html"
	// PORT is the port the server will listen on.
	PORT = "8080"
)

// SetupServer creates a server, and sets up routes, middleware and assets.
func SetupServer() *http.Server {
	log.Println("Server starting...")
	gin.SetMode(gin.DebugMode)
	log.SetOutput(gin.DefaultWriter)
	r := gin.New()

	// Global middleware to prevent caching of all files
	r.Use(func(c *gin.Context) {
		// Set Cache-Control header for all responses
		c.Header("Cache-Control", "no-store, no-cache, must-revalidate, proxy-revalidate")
		c.Next()
	}, ErrorHandler())

	r.SetTrustedProxies(nil)
	loadHTMLTemplates(r)
	// r.LoadHTMLGlob(templatesDir)
	registerRoutes(r)
	serveStaticAssets(r)

	return &http.Server{
		Addr:    "127.0.0.1:" + PORT,
		Handler: r,
	}
}

func registerRoutes(r *gin.Engine) {
	r.GET("/", gin.HandlerFunc(HomeHandler))
}

// ErrorHandler is our error handling Middleware.
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Create a deferred function to catch any errors that occur
		// defer func() {
		// 	if err := recover(); err != nil {
		// 		// Log the error (you can also log to a file or external service)
		// 		log.Printf("Error occurred: %v", err)

		// 		// Respond with a 500 Internal Server Error and a message
		// 		c.JSON(http.StatusInternalServerError, gin.H{
		// 			"message": "Internal Server Error",
		// 		})
		// 	}
		// }()

		// Continue with the next handler
		c.Next()

		// Log errors that ocurred in the handler
		for _, ginErr := range c.Errors {
			log.Println(ginErr)
		}
	}
}
func serveStaticAssets(r *gin.Engine) {
	// TODO: Serve static files from the embedded file system
	// r.StaticFS("/static", http.FS(frontend.TemplateFS))

	r.Static("/static", assetsDir)
}

// loadHTMLTemplates loads the HTML templates from the frontend/templates directory
func loadHTMLTemplates(r *gin.Engine) {
	// Load templates from the embedded file system
	tmpl := template.Must(template.New("").ParseFS(frontend.TemplateFS, "templates/**/*"))

	// Set Gin to use the embedded templates
	r.SetHTMLTemplate(tmpl)
}
