package web

import (
	"errors"
	"net/http"
	"time"

	"github.com/KarolinaLop/dp/data"
	"github.com/KarolinaLop/dp/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// ShowRegistrationFormHandler handles the GET request for the registration form.
func ShowRegistrationFormHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "registration.html", gin.H{})
}

type registrationForm struct {
	Username string `form:"username" binding:"required"`
	Email    string `form:"email" binding:"required,email"`
	Password string `form:"password" binding:"required"`
}

// RegisterUserHandler handles the POST request for the registration form.
func RegisterUserHandler(c *gin.Context) {
	var form registrationForm

	if err := c.ShouldBind(&form); err != nil {
		c.Error(err)
		return
	}

	// Validate the form data
	if form.Username == "" || form.Email == "" || form.Password == "" {
		c.Error(errors.New("all fields are required"))
		return
	}

	exists, err := data.UserExists(data.DB, form.Email)
	if err != nil {
		c.Error(err)
		return
	}

	if exists {
		c.Error(errors.New("user already exists"))
		return
	}

	user := models.User{
		Name:      form.Username,
		Email:     form.Email,
		Password:  form.Password,
		CreatedAt: time.Now().Format("2006-01-02 15:04:05"),
	}
	// Hash the password before storing it
	user.PasswordHash, err = hashPassword(form.Password)
	if err != nil {
		c.Error(err)
		return
	}

	user, err = data.CreateUser(data.DB, user)
	if err != nil {
		c.Error(err)
		return
	}

	// Create a session
	if err = createUserSession(c, user); err != nil {
		c.Error(err)
		return
	}

	c.Redirect(http.StatusSeeOther, "/dashboard")
}

// hashPassword hashes the password using bcrypt.
func hashPassword(password string) (string, error) {
	// Hash the password using bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}
