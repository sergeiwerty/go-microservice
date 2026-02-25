package auth

import (
	"context"
	"errors"
	"go-microservice/internal/models"
	"sync"

	"golang.org/x/crypto/bcrypt"
)

type UserRepository struct {
	mu     sync.Mutex
	data   map[string]models.User
	nextID int
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		data:   make(map[string]models.User),
		nextID: 1,
	}
}

func (ur *UserRepository) Create(ctx context.Context, user models.User) (models.User, error) {
	ur.mu.Lock()
	defer ur.mu.Unlock()

	if _, exists := ur.data[user.Email]; exists {
		return models.User{}, errors.New("user already exists")
	}

	user.ID = ur.nextID
	ur.data[user.Email] = user
	ur.nextID++

	return user, nil
}

func (ur *UserRepository) Save(u models.User) {
	ur.mu.Lock()
	defer ur.mu.Unlock()
	ur.data[u.Email] = u
}

func (ur *UserRepository) Find(username string, password string) (models.User, bool) {
	ur.mu.Lock()
	defer ur.mu.Unlock()
	u, p_exists := ur.data[username]

	if !p_exists {
		return models.User{}, false
	}

	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return models.User{}, false
	}

	return u, true
}

func (ur *UserRepository) ExistsByEmail(ctx context.Context, email string) (models.User, bool) {
	ur.mu.Lock()
	defer ur.mu.Unlock()
	user, exists := ur.data[email]

	return user, exists
}
