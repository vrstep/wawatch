package controller_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/vrstep/wawatch-backend/config" // Direct DB access for setup/verification
	"github.com/vrstep/wawatch-backend/models"
	"golang.org/x/crypto/bcrypt" // For password hashing in setup
)

// testRouter and testDB are assumed to be initialized by TestMain

func TestSignup_Success(t *testing.T) {
	clearUserRelatedTables() // Use your helper from main_test.go or define locally

	signupPayload := gin.H{
		"username": "testsignup",
		"password": "password123",
		"email":    "signup@example.com",
	}
	rr := performRequest("POST", "/api/v1/auth/signup", signupPayload, testRouter)

	assert.Equal(t, http.StatusCreated, rr.Code)
	var response map[string]interface{}
	json.Unmarshal(rr.Body.Bytes(), &response)
	assert.Equal(t, "User created successfully", response["message"])

	// Verify user in DB
	var user models.User
	err := config.DB.Where("username = ?", "testsignup").First(&user).Error
	assert.NoError(t, err)
	assert.Equal(t, "signup@example.com", user.Email)
}

func TestSignup_DuplicateUsername(t *testing.T) {
	clearUserRelatedTables()
	// Create an initial user
	existingUser := models.User{Username: "existinguser", Password: "password", Email: "existing@example.com"}
	config.DB.Create(&existingUser)

	signupPayload := gin.H{"username": "existinguser", "password": "password123"}
	rr := performRequest("POST", "/api/v1/auth/signup", signupPayload, testRouter)

	assert.Equal(t, http.StatusConflict, rr.Code) // Or whatever your controller returns for duplicates
}

// TestLogin_Success was already outlined, ensure it's here.

func TestChangePassword_Success(t *testing.T) {
	user, token := createAndLoginTestUser(config.DB, "changepwuser", "oldpassword")

	changePwPayload := gin.H{
		"current_password": "oldpassword",
		"new_password":     "newpassword123",
	}
	rr := performAuthRequest("PUT", "/api/v1/me/profile/password", changePwPayload, token, testRouter)

	assert.Equal(t, http.StatusOK, rr.Code)
	var response map[string]string
	json.Unmarshal(rr.Body.Bytes(), &response)
	assert.Equal(t, "Password changed successfully", response["message"])

	// Verify new password in DB (fetch user and try to compare hash)
	var updatedUser models.User
	config.DB.First(&updatedUser, user.ID)
	err := bcrypt.CompareHashAndPassword([]byte(updatedUser.Password), []byte("newpassword123"))
	assert.NoError(t, err, "New password should match")
}

func TestChangePassword_IncorrectCurrent(t *testing.T) {
	_, token := createAndLoginTestUser(config.DB, "changepwuser2", "oldpassword")

	changePwPayload := gin.H{
		"current_password": "wrongoldpassword",
		"new_password":     "newpassword123",
	}
	rr := performAuthRequest("PUT", "/api/v1/me/profile/password", changePwPayload, token, testRouter)
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestChangeUsername_Success(t *testing.T) {
	user, token := createAndLoginTestUser(config.DB, "changeuname", "password123")

	payload := gin.H{
		"new_username":     "newuname",
		"current_password": "password123",
	}
	rr := performAuthRequest("PUT", "/api/v1/me/profile/username", payload, token, testRouter)

	assert.Equal(t, http.StatusOK, rr.Code)
	var dbUser models.User
	config.DB.First(&dbUser, user.ID)
	assert.Equal(t, "newuname", dbUser.Username)
}

func TestChangeEmail_Success(t *testing.T) {
	user, token := createAndLoginTestUser(config.DB, "changeemailuser", "password123")

	payload := gin.H{
		"new_email":        "new@example.com",
		"current_password": "password123",
	}
	rr := performAuthRequest("PUT", "/api/v1/me/profile/email", payload, token, testRouter)

	assert.Equal(t, http.StatusOK, rr.Code)
	var dbUser models.User
	config.DB.First(&dbUser, user.ID)
	assert.Equal(t, "new@example.com", dbUser.Email)
}

// Add tests for:
// - ChangeUsername/Email with incorrect password
// - ChangeUsername/Email with new username/email already taken
// - Input validation failures (e.g., short new password, invalid email format)
