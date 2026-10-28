package handlers

import (
	"errors"
	"net/http"

	"backend/models"
	"backend/services"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService services.AuthService
}

func NewAuthHandler(service services.AuthService) *AuthHandler {
	return &AuthHandler{authService: service}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "กรุณากรอกข้อมูลให้ครบถ้วน: " + err.Error()})
		return
	}

	if err := h.authService.Register(&req); err != nil {
		switch {
		case errors.Is(err, services.ErrPasswordMismatch), errors.Is(err, services.ErrUserAlreadyExists):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "เกิดข้อผิดพลาดในการลงทะเบียน"})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "กรุณากรอกข้อมูลให้ครบถ้วน: " + err.Error()})
		return
	}

	res, err := h.authService.Login(&req)
	if err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "เกิดข้อผิดพลาดในการเข้าสู่ระบบ"})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *AuthHandler) Profile(c *gin.Context) {
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ไม่พบข้อมูลยืนยันตัวตน"})
		return
	}

	res, err := h.authService.GetProfile(username.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ไม่สามารถดึงข้อมูลโปรไฟล์ได้"})
		return
	}

	c.JSON(http.StatusOK, res)
}
