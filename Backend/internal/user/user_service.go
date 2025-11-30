package user

import (
	"context"
	"errors"
	"os"
	"github.com/golang-jwt/jwt/v5"
	"time"

	"github.com/Luemax58/be-fe-project/pkg/models"
	"golang.org/x/crypto/bcrypt"
)

type IUserService interface {
    Login(ctx context.Context, username, password string) (*models.User, string, error)
    GetUserByID(ctx context.Context, userID uint) (*models.User, error)
}

type userService struct {
    repo IUserRepository
}

func NewUserService(repo IUserRepository) IUserService {
    return &userService{repo: repo}
}

func (s *userService) Login(ctx context.Context, username, password string) (*models.User, string, error) {
    user, err := s.repo.GetUserByUsername(ctx, username)
    if err != nil {
        return nil, "", err
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
        return nil, "", errors.New("invalid username or password")
    }

    token, err := generateJWT(user)
    if err != nil {
        return nil, "", err
    }

    return user, token, nil
}

func (s *userService) GetUserByID(ctx context.Context, userID uint) (*models.User, error) {
    return s.repo.GetUserByID(ctx, userID)
}

func generateJWT(user *models.User) (string, error) {
    secretKey := os.Getenv("JWT_SECRET_KEY")
    if secretKey == "" {
        secretKey = "my-secret-key-1234"
    }

    claims := jwt.MapClaims{
        "user_id":  user.UserID,
        "username": user.Username,
        "role":     user.Role,
        "exp":      time.Now().Add(24 * time.Hour).Unix(),
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(secretKey))
}
