package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"backend/database"
	"backend/handlers"
	"backend/middleware"
	"backend/repository"
	"backend/services"

	"github.com/gin-gonic/gin"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	_ = os.Remove("test_app.db")

	db := database.InitDB("test_app.db")
	userRepo := repository.NewUserRepository(db)
	authService := services.NewAuthService(userRepo)
	authHandler := handlers.NewAuthHandler(authService)

	r := gin.Default()
	r.Use(middleware.CORSMiddleware("http://localhost:4200"))

	api := r.Group("/api")
	{
		api.POST("/register", authHandler.Register)
		api.POST("/login", authHandler.Login)

		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware())
		{
			protected.GET("/profile", authHandler.Profile)
		}
	}

	return r
}

func TestAuthWorkflow(t *testing.T) {
	router := setupTestRouter()
	defer os.Remove("test_app.db")

	t.Run("CORS Preflight", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodOptions, "/api/register", nil)
		req.Header.Set("Origin", "http://localhost:4200")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusNoContent {
			t.Fatalf("Expected status %d, got %d", http.StatusNoContent, w.Code)
		}
		if origin := w.Header().Get("Access-Control-Allow-Origin"); origin != "http://localhost:4200" {
			t.Fatalf("Expected Access-Control-Allow-Origin 'http://localhost:4200', got '%s'", origin)
		}
	})

	t.Run("Register Password Mismatch", func(t *testing.T) {
		body, _ := json.Marshal(map[string]string{
			"username":         "tester1",
			"password":         "pass1234",
			"confirm_password": "pass9999",
		})
		req, _ := http.NewRequest(http.MethodPost, "/api/register", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("Register Success", func(t *testing.T) {
		body, _ := json.Marshal(map[string]string{
			"username":         "tester1",
			"password":         "pass1234",
			"confirm_password": "pass1234",
		})
		req, _ := http.NewRequest(http.MethodPost, "/api/register", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Fatalf("Expected status %d, got %d: %s", http.StatusCreated, w.Code, w.Body.String())
		}
	})

	t.Run("Register Duplicate Username", func(t *testing.T) {
		body, _ := json.Marshal(map[string]string{
			"username":         "tester1",
			"password":         "pass1234",
			"confirm_password": "pass1234",
		})
		req, _ := http.NewRequest(http.MethodPost, "/api/register", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("Login Incorrect Password", func(t *testing.T) {
		body, _ := json.Marshal(map[string]string{
			"username": "tester1",
			"password": "wrongpassword",
		})
		req, _ := http.NewRequest(http.MethodPost, "/api/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Fatalf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
		}
	})

	var token string
	t.Run("Login Success", func(t *testing.T) {
		body, _ := json.Marshal(map[string]string{
			"username": "tester1",
			"password": "pass1234",
		})
		req, _ := http.NewRequest(http.MethodPost, "/api/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status %d, got %d", http.StatusOK, w.Code)
		}

		var res map[string]interface{}
		_ = json.Unmarshal(w.Body.Bytes(), &res)
		tok, ok := res["token"].(string)
		if !ok || tok == "" {
			t.Fatalf("Token not found in response: %v", res)
		}
		if res["username"] != "tester1" {
			t.Fatalf("Expected username tester1, got %v", res["username"])
		}
		token = tok
	})

	t.Run("Profile Missing Token", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/api/profile", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Fatalf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
		}
	})

	t.Run("Profile Valid Token", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/api/profile", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status %d, got %d: %s", http.StatusOK, w.Code, w.Body.String())
		}

		var res map[string]interface{}
		_ = json.Unmarshal(w.Body.Bytes(), &res)
		if res["username"] != "tester1" {
			t.Fatalf("Expected username tester1, got %v", res["username"])
		}
		if res["status"] != "verified" {
			t.Fatalf("Expected status verified, got %v", res["status"])
		}
	})
}
