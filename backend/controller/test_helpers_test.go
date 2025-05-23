// backend/controller/main_test.go or a new backend/controller/test_helpers_test.go
package controller_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/vrstep/wawatch-backend/models" // For models.User
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Use the JWT secret from your actual application logic for consistency in tests
// It's hardcoded as "secret" in your middleware.
var testJWTSecret = []byte("secret")

func generateTestToken(userID uint, username string, role string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"unm": username,
		"rol": role,
		"exp": time.Now().Add(time.Hour * 24).Unix(), // Token valid for 24 hours
		"iat": time.Now().Unix(),
	})
	tokenString, err := token.SignedString(testJWTSecret)
	if err != nil {
		panic(fmt.Sprintf("Failed to sign test token: %v", err))
	}
	return tokenString
}

func performRequest(method, path string, body interface{}, router *gin.Engine) *httptest.ResponseRecorder {
	var reqBody *bytes.Buffer
	if body != nil {
		b, _ := json.Marshal(body)
		reqBody = bytes.NewBuffer(b)
	} else {
		reqBody = bytes.NewBuffer(nil)
	}

	req, _ := http.NewRequest(method, path, reqBody)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	return rr
}

func performAuthRequest(method, path string, body interface{}, token string, router *gin.Engine) *httptest.ResponseRecorder {
	var reqBody *bytes.Buffer
	if body != nil {
		b, _ := json.Marshal(body)
		reqBody = bytes.NewBuffer(b)
	} else {
		reqBody = bytes.NewBuffer(nil)
	}

	req, _ := http.NewRequest(method, path, reqBody)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	// For Gin's c.Cookie("Auth"), we need to simulate the cookie.
	// Alternatively, if your RequireAuth also checks Authorization Bearer, that's easier.
	// Your RequireAuth specifically looks for a cookie.
	if token != "" {
		req.AddCookie(&http.Cookie{Name: "Auth", Value: token, Path: "/"})
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	return rr
}

// Helper to create a user and get a token for tests
func createAndLoginTestUser(db *gorm.DB, username, password string) (models.User, string) {
	clearUserRelatedTables() // Clear before creating to avoid conflicts if called multiple times
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user := models.User{Username: username, Password: string(hashedPassword), Email: username + "@example.com"}
	db.Create(&user)
	token := generateTestToken(user.ID, user.Username, user.Role)
	return user, token
}
