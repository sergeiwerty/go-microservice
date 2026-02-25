package auth

import (
	"context"
	"errors"
	"go-microservice/internal/models"

	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Login(ctx context.Context, email, password string) (bool, error)
	Register(ctx context.Context, email, plainPassword string) (*models.User, error)
}

type Service struct {
	Repo *UserRepository
}

func NewAuthService(repo *UserRepository) *Service {
	return &Service{Repo: repo}
}

func (s Service) Login(ctx context.Context, username, password string) (bool, error) {
	_, exists := s.Repo.Find(username, password)
	if exists {
		return true, nil
	}

	return false, nil
}

var ErrUserAlreadyExists = errors.New("user already exists")

func (s Service) Register(ctx context.Context, email, plainPassword string) (*models.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)

	newUser := models.User{
		Email:    email,
		Password: string(hashedPassword),
	}

	u, err := s.Repo.Create(ctx, newUser)

	if err != nil {
		return nil, err
	}

	return &u, nil
}
