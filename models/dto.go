package models

type RegisterRequest struct {
	Username        string `json:"username" binding:"required"`
	Password        string `json:"password" binding:"required"`
	ConfirmPassword string `json:"confirm_password" binding:"required"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Message  string `json:"message"`
	Token    string `json:"token"`
	Username string `json:"username"`
}

type ProfileResponse struct {
	Username string `json:"username"`
	Status   string `json:"status"`
	Message  string `json:"message"`
}
