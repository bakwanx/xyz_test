package usecase

import (
	"errors"
	"time"
	"xyz_backend/config"
	"xyz_backend/src/models"
	"xyz_backend/src/repository"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthUsecase struct {
	adminRepo repository.AdminRepository
	tokenRepo repository.TokenBlacklistRepository
	cfg       *config.Config
}

func NewAuthUsecase(ar repository.AdminRepository, tr repository.TokenBlacklistRepository, cfg *config.Config) *AuthUsecase {
	return &AuthUsecase{adminRepo: ar, tokenRepo: tr, cfg: cfg}
}

func (u *AuthUsecase) RegisterAdmin(name, email, password string) (*models.Admin, error) {

	checkRegisteredUser, err := u.adminRepo.FindByEmail(email)
	if err != nil {
		return nil, err
	}

	if checkRegisteredUser != nil {
		return nil, errors.New("User has been registered")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	a := &models.Admin{Name: name, Email: email, PasswordHash: string(hash), Role: "admin"}
	if err := u.adminRepo.Create(a); err != nil {
		return nil, err
	}
	return a, nil
}

func (u *AuthUsecase) Login(email, password string) (string, error) {
	admin, err := u.adminRepo.FindByEmail(email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}
	if admin == nil {
		return "", errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(admin.PasswordHash), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	now := time.Now().UTC()
	exp := now.Add(time.Duration(u.cfg.JWT.ExpiryMinutes) * time.Minute)
	jti := uuid.New().String()

	claims := jwt.MapClaims{
		"sub":   admin.ID,
		"email": admin.Email,
		"role":  admin.Role,
		"exp":   exp.Unix(),
		"iat":   now.Unix(),
		"jti":   jti,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(u.cfg.JWT.Secret))
	if err != nil {
		return "", err
	}
	return signed, nil
}

func (u *AuthUsecase) Logout(jti string, exp time.Time) error {
	return u.tokenRepo.Blacklist(jti, exp)
}
