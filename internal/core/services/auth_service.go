package services

import (
	"context"
	"errors"
	"time"

	"ApiSmart/internal/core/domain"
	"ApiSmart/internal/core/ports"
	"ApiSmart/pkg/auth"
	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	userRepo ports.UserRepository
}

func NewAuthService(userRepo ports.UserRepository) ports.AuthService {
	return &authService{
		userRepo: userRepo,
	}
}

func (s *authService) Register(ctx context.Context, req domain.RegisterRequest) (*domain.AuthResponse, error) {
	// Comprobar si el usuario ya existe
	existingUser, _ := s.userRepo.FindByEmail(ctx, req.Email)
	if existingUser != nil {
		return nil, errors.New("el correo electrónico ya está registrado")
	}

	// Hash de la contraseña
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Crear nuevo usuario
	now := time.Now()
	user := &domain.User{
		Username:  req.Username,
		Email:     req.Email,
		Password:  string(hashedPassword),
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Guardar en base de datos
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	// Generar token JWT
	token, err := auth.GenerateJWT(user.ID, user.Username)
	if err != nil {
		return nil, err
	}

	return &domain.AuthResponse{
		Token:    token,
		Username: user.Username,
		Email:    user.Email,
	}, nil
}

func (s *authService) Login(ctx context.Context, req domain.LoginRequest) (*domain.AuthResponse, error) {
	// Buscar usuario por email
	user, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.New("credenciales inválidas")
	}

	// Verificar contraseña
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, errors.New("credenciales inválidas")
	}

	// Generar token JWT
	token, err := auth.GenerateJWT(user.ID, user.Username)
	if err != nil {
		return nil, err
	}

	return &domain.AuthResponse{
		Token:    token,
		Username: user.Username,
		Email:    user.Email,
	}, nil
}

func (s *authService) ValidateToken(token string) (uint, error) {
	return auth.ValidateJWT(token)
}
