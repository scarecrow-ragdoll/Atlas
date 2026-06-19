package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"

	"monorepo-template/libs/go/logger"
)

type User struct {
	ID        string
	Email     string
	Name      string
	CreatedAt string
	UpdatedAt string
}

type CreateUserInput struct {
	Email    string
	Name     string
	Password string
}

type UpdateUserInput struct {
	Name  *string
	Email *string
}

var (
	ErrDuplicateEmail = errors.New("duplicate email")
	ErrNotFound       = errors.New("not found")
)

func IsDuplicateEmail(err error) bool {
	if err == nil {
		return false
	}
	return errors.Is(err, ErrDuplicateEmail) || strings.Contains(err.Error(), "duplicate email")
}

func IsNotFound(err error) bool {
	return errors.Is(err, ErrNotFound)
}

type UserRepository interface {
	GetByID(ctx context.Context, id string) (*User, error)
	List(ctx context.Context, first *int, after *string) ([]*User, int, error)
	Create(ctx context.Context, input CreateUserInput) (*User, error)
	Update(ctx context.Context, id string, input UpdateUserInput) (*User, error)
	Delete(ctx context.Context, id string) error
}

type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetByID(ctx context.Context, id string) (*User, error) {
	const op = "UserService.GetByID"
	log := logger.FromContext(ctx).With(zap.String("op", op))
	log.Debug("getting user", zap.String("user_id", id))
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return user, nil
}

func (s *UserService) List(ctx context.Context, first *int, after *string) ([]*User, int, error) {
	const op = "UserService.List"
	log := logger.FromContext(ctx).With(zap.String("op", op))
	log.Debug("listing users")
	users, total, err := s.repo.List(ctx, first, after)
	if err != nil {
		return nil, 0, fmt.Errorf("%s: %w", op, err)
	}
	return users, total, nil
}

func (s *UserService) Create(ctx context.Context, input CreateUserInput) (*User, error) {
	const op = "UserService.Create"
	log := logger.FromContext(ctx).With(zap.String("op", op))
	log.Debug("creating user", zap.String("email", input.Email))
	hashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("%s: hash password: %w", op, err)
	}
	input.Password = string(hashed)
	user, err := s.repo.Create(ctx, input)
	if err != nil {
		if IsDuplicateEmail(err) {
			return nil, fmt.Errorf("%s: %w", op, ErrDuplicateEmail)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return user, nil
}

func (s *UserService) Update(ctx context.Context, id string, input UpdateUserInput) (*User, error) {
	const op = "UserService.Update"
	log := logger.FromContext(ctx).With(zap.String("op", op))
	log.Debug("updating user", zap.String("user_id", id))
	user, err := s.repo.Update(ctx, id, input)
	if err != nil {
		if IsDuplicateEmail(err) {
			return nil, fmt.Errorf("%s: %w", op, ErrDuplicateEmail)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return user, nil
}

func (s *UserService) Delete(ctx context.Context, id string) error {
	const op = "UserService.Delete"
	log := logger.FromContext(ctx).With(zap.String("op", op))
	log.Debug("deleting user", zap.String("user_id", id))
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
