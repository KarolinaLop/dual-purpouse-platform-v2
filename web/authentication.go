package web

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/KarolinaLop/dp/data"
	"github.com/KarolinaLop/dp/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// Authentication is a middleware that checks if the user is authenticated and redirects to the login page if not.
func Authentication(c *gin.Context) {
	if !isAuthenticated(c) {
		c.Redirect(http.StatusFound, "/login")
		c.Abort()
		return
	}
	c.Next()
}

func isAuthenticated(c *gin.Context) bool {
	// Get the session from the request
	session := sessions.Default(c)

	// Check if the user is logged in
	userID := session.Get("user_id")
	if userID == nil {
		return false
	}
	user, err := data.GetUserByID(data.DB, userID.(int))
	if err != nil {
		return false
	}

	// Store the user in the context
	c.Set("user", user)
	// Refresh the session expiration time, so the user stays logged in
	session.Set("time", time.Now().Format("2006-01-02 15:04:05"))
	if err := session.Save(); err != nil {
		log.Printf("Failed to save session: %v", err)
		return false
	}

	return true
}

// ShowLoginForm renders the login form.
func ShowLoginForm(c *gin.Context) {
	// Render the login form
	c.HTML(http.StatusOK, "login.html", gin.H{
		"title": "Login",
	})
}

type loginForm struct {
	Email    string `form:"email" binding:"required"`
	Password string `form:"password" binding:"required"`
}

// LoginUser handles the login form submission.
func LoginUser(c *gin.Context) {
	var form loginForm
	if err := c.ShouldBind(&form); err != nil {
		c.Error(err)
		return
	}

	// Get the user by email
	user, err := data.GetUserByEmail(data.DB, form.Email)
	if err != nil {
		if errors.Is(err, data.ErrUserNotFound) {
			// TODO: pass an error message to the template
			c.HTML(http.StatusBadRequest, "login.html", gin.H{
				"title":   "Login",
				"message": "User not found",
			})

			c.Abort()
			return
		}
		c.Error(err)
		c.Abort()
		return
	}

	// Compare the password with the hashed password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(form.Password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			// TODO: pass an error message to the template
			c.HTML(http.StatusBadRequest, "login.html", gin.H{
				"title":   "Login",
				"message": "Invalid password",
			})

			c.Abort()
			return
		}
		c.Error(err)
		return
	}

	// Create a session
	if err = createUserSession(c, user); err != nil {
		c.Error(err)
		return
	}

	// Redirect to the main page
	c.Redirect(http.StatusFound, "/scans")
}

func createUserSession(c *gin.Context, user models.User) error {
	session := sessions.Default(c)
	session.Set("user_id", user.ID)
	session.Set("user_name", user.Name)
	return session.Save()
}

// LogoutUser handles the logout action.
// It clears the session and redirects to the login page.
func LogoutUser(c *gin.Context) {

	// Get the session
	session := sessions.Default(c)
	if err := data.DeleteSessions(data.DB, session.Get("user_id").(int)); err != nil {
		c.Error(err)
		return
	}

	// Delete the cookie
	session.Options(sessions.Options{
		MaxAge:   -1, // Mark the cookie as expired
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
	session.Save()

	// Render the login form
	c.HTML(http.StatusOK, "login.html", gin.H{
		"title": "Login",
	})
}
