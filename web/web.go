package web

import (
	"encoding/hex"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/KarolinaLop/dp/data"
	"github.com/KarolinaLop/dp/frontend"
	"github.com/gin-contrib/sessions"

	"github.com/gin-gonic/gin"
)

const (
	// PORT is the port the server will listen on.
	PORT = "8080"
)

// SetupServer creates a server, and sets up routes, middleware and assets.
func SetupServer() *http.Server {
	log.Println("Server starting...")
	ginMode := os.Getenv("GIN_MODE")
	gin.SetMode(ginMode)
	log.SetOutput(gin.DefaultWriter)
	r := gin.New()

	// TODO: make this configurable by loading the keys from environment variables
	encKeyHex := "99268541414541b9b9982c4b7a3de7c59b25b6f9dee0f9308c988732bc54e9f6"
	encKey, err := hex.DecodeString(encKeyHex)
	if err != nil {
		log.Fatalf("Failed to decode encKey: %v", err)
	}
	authKeyHex := "8a88674ad14dc1f0e95b4699cec94751e1f2762ee1e92dc95d82a430e03e52cd"
	authKey, err := hex.DecodeString(authKeyHex)
	if err != nil {
		log.Fatalf("Failed to decode authKey: %v", err)
	}
	// Set up session store
	store := NewSQLiteStore(data.DB, authKey, encKey)
	store.Options(sessions.Options{
		// Secure:   true,
		Path:     "/",
		MaxAge:   3600, // 1 hour
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})

	// Set up middleware
	r.Use(gin.Logger(), ErrorHandlerMiddleware(), CacheControlMiddleware(), sessions.Sessions("dp-session", store))

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
		Addr:    "0.0.0.0:" + PORT,
		Handler: r,
	}
}

func registerRoutes(r *gin.Engine) {

	authenticatedRoutes := r.Group("/", AuthenticationMiddleware)
	authenticatedRoutes.GET("/scans", ShowScansListHandler)
	authenticatedRoutes.DELETE("/logout", LogoutUserHandler)
	authenticatedRoutes.POST("/scans", StartScanHandler)
	authenticatedRoutes.GET("/scans/:id/show", ShowScanDetailsHandler)
	authenticatedRoutes.DELETE("/scans/:id", DeleteScanHandler)
	authenticatedRoutes.GET("/scans/:id/status", CheckScanStatusHandler)

	r.GET("/", HomeHandler)
	r.GET("/register", ShowRegistrationFormHandler)
	r.POST("/register", RegisterUserHandler)
	r.GET("/login", ShowLoginFormHandler)
	r.POST("/login", LoginUserHandler)
}

// ErrorHandlerMiddleware is an error handling Middleware.
func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if gin.Mode() == gin.ReleaseMode {
			// Recover from panic and log the error
			defer func() {
				if err := recover(); err != nil {
					// Log the error
					log.Printf("Error occurred: %v", err)
					// Set the error on the context
					c.Error(fmt.Errorf("%v", err))

					renderErrorPage(c, http.StatusInternalServerError, "An internal server error occurred.")
				}
			}()
		}
		// Continue with the next handler
		c.Next()

		// Log errors that ocurred in the handler
		for _, ginErr := range c.Errors {
			log.Println(ginErr)
		}

		// Render error page
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

// CacheControlMiddleware is a middleware that sets the Cache-Control header to prevent caching.
func CacheControlMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Cache-Control", "no-store, no-cache, must-revalidate, proxy-revalidate")
		c.Next()
	}
}

func serveStaticAssets(r *gin.Engine) {
	// Serves assets from the embedded file system
	r.StaticFS("/static", http.FS(frontend.TemplateFS))
}

// loadHTMLTemplates loads the HTML templates from the frontend/templates directory
func loadHTMLTemplates(r *gin.Engine) {
	// Load templates from the embedded file system
	tmpl := template.Must(template.New("").ParseFS(frontend.TemplateFS, "templates/**/*"))

	// Set Gin to use the embedded templates
	r.SetHTMLTemplate(tmpl)
}
