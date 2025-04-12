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

	// Set up middleware
	r.Use(gin.Logger(), ErrorHandler(), CacheControl())
	r.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusNotFound, "error.html", gin.H{
			"errorCode":    http.StatusNotFound,
			"errorMessage": "That page does not exist.",
		})
	})

	r.SetTrustedProxies(nil)
	loadHTMLTemplates(r)
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
		// Recover from panic and log the error
		defer func() {
			if err := recover(); err != nil {
				// Log the error (you can also log to a file or external service)
				log.Printf("Error occurred: %v", err)
				// set the error on the context
				c.Error(err.(error))

				renderErrorPage(c, http.StatusInternalServerError, "An internal server error occurred.")
			}
		}()

		// Continue with the next handler
		c.Next()

		// Log errors that ocurred in the handler
		for _, ginErr := range c.Errors {
			log.Println(ginErr)
		}

		// render error page
		if len(c.Errors) > 0 {
			renderErrorPage(c, http.StatusInternalServerError, "An error occurred while processing your request.")
		}
	}
}

func renderErrorPage(c *gin.Context, statusCode int, message string) {
	c.HTML(statusCode, "error.html", gin.H{
		"errorCode":    statusCode,
		"errorMessage": message,
	})
}

// CacheControl is a middleware that sets the Cache-Control header to prevent caching.
func CacheControl() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Cache-Control", "no-store, no-cache, must-revalidate, proxy-revalidate")
		c.Next()
	}
}

func serveStaticAssets(r *gin.Engine) {
	// serves assets from the embedded file system
	r.StaticFS("/static", http.FS(frontend.TemplateFS))
}

// loadHTMLTemplates loads the HTML templates from the frontend/templates directory
func loadHTMLTemplates(r *gin.Engine) {
	// Load templates from the embedded file system
	tmpl := template.Must(template.New("").ParseFS(frontend.TemplateFS, "templates/**/*"))

	// Set Gin to use the embedded templates
	r.SetHTMLTemplate(tmpl)
}
