package service

import (
	"vblog-core/model"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// AuthService handles user authentication.
type AuthService struct {
	DB *gorm.DB
}

// NewAuthService creates a new AuthService.
func NewAuthService(db *gorm.DB) *AuthService {
	return &AuthService{DB: db}
}

// HashPassword hashes a plaintext password using bcrypt.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPassword compares a bcrypt hash with a plaintext password.
func CheckPassword(hash, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

// Login authenticates a user by username and password.
func (s *AuthService) Login(username, password string) (*model.User, error) {
	var user model.User
	if err := s.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	if !CheckPassword(user.Password, password) {
		return nil, bcrypt.ErrMismatchedHashAndPassword
	}
	return &user, nil
}
