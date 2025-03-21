package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/tassyosilva/consultapix/internal/config"
	"github.com/tassyosilva/consultapix/internal/database/models"
	"github.com/tassyosilva/consultapix/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo *repository.UserRepository
	config   *config.Config
}

type JWTClaims struct {
	jwt.RegisteredClaims
	ID        int    `json:"id"`
	CPF       string `json:"cpf"`
	Nome      string `json:"name"`
	Email     string `json:"email"`
	Lotacao   string `json:"lotacao"`
	Matricula string `json:"matricula"`
	Admin     bool   `json:"admin"`
}

func NewAuthService(cfg *config.Config) *AuthService {
	return &AuthService{
		userRepo: repository.NewUserRepository(),
		config:   cfg,
	}
}

func (s *AuthService) Login(email, password string) (string, *models.Usuario, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return "", nil, errors.New("email ou senha inválidos")
	}

	// Verificar senha
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", nil, errors.New("email ou senha inválidos")
	}

	// Criar token JWT
	claims := JWTClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(s.config.TokenExpiryHours))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		ID:        user.ID,
		CPF:       user.CPF,
		Nome:      user.Nome,
		Email:     user.Email,
		Lotacao:   user.Lotacao,
		Matricula: user.Matricula,
		Admin:     user.Admin,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	
	// Assinar token
	tokenString, err := token.SignedString([]byte(s.config.JWTSecret))
	if err != nil {
		return "", nil, err
	}

	return tokenString, user, nil
}

func (s *AuthService) ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.config.JWTSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("token inválido")
}