package services

import (
	"errors"

	"backend/models"
	"backend/repository"
	"backend/utils"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrPasswordMismatch   = errors.New("รหัสผ่านและยืนยันรหัสผ่านไม่ตรงกัน")
	ErrUserAlreadyExists  = errors.New("ชื่อผู้ใช้นี้ถูกใช้งานไปแล้ว")
	ErrInvalidCredentials = errors.New("ชื่อผู้ใช้งานหรือรหัสผ่านไม่ถูกต้อง")
	ErrInternalServer     = errors.New("เกิดข้อผิดพลาดภายในเซิร์ฟเวอร์")
)

type AuthService interface {
	Register(req *models.RegisterRequest) error
	Login(req *models.LoginRequest) (*models.LoginResponse, error)
	GetProfile(username string) (*models.ProfileResponse, error)
}

type authService struct {
	userRepo repository.UserRepository
}

func NewAuthService(repo repository.UserRepository) AuthService {
	return &authService{userRepo: repo}
}

func (s *authService) Register(req *models.RegisterRequest) error {
	if req.Password != req.ConfirmPassword {
		return ErrPasswordMismatch
	}

	existingUser, err := s.userRepo.FindByUsername(req.Username)
	if err != nil {
		return ErrInternalServer
	}
	if existingUser != nil {
		return ErrUserAlreadyExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return ErrInternalServer
	}

	newUser := &models.User{
		Username: req.Username,
		Password: string(hashedPassword),
	}

	if err := s.userRepo.Create(newUser); err != nil {
		return ErrInternalServer
	}

	return nil
}

func (s *authService) Login(req *models.LoginRequest) (*models.LoginResponse, error) {
	user, err := s.userRepo.FindByUsername(req.Username)
	if err != nil || user == nil {
		return nil, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	token, err := utils.GenerateToken(user.ID, user.Username)
	if err != nil {
		return nil, ErrInternalServer
	}

	return &models.LoginResponse{
		Message:  "เข้าสู่ระบบสำเร็จ",
		Token:    token,
		Username: user.Username,
	}, nil
}

func (s *authService) GetProfile(username string) (*models.ProfileResponse, error) {
	return &models.ProfileResponse{
		Username: username,
		Status:   "verified",
		Message:  "Authenticated successfully",
	}, nil
}
