// FILE: apps/api/internal/service/admin_auth.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Own transport-neutral web-admin authentication service contracts and behavior.
//   SCOPE: Admin domain types, repository/session interfaces, bootstrap seed, login, current admin, create admin, and logout; excludes PostgreSQL, Redis, GraphQL, and HTTP cookie adapters.
//   DEPENDS: context, errors.
//   LINKS: M-API / V-M-API.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   AdminRoleAdmin - Initial admin role value.
//   Admin - Public admin identity returned to GraphQL and context consumers.
//   CreateAdminInput - Repository/service admin creation input.
//   AdminRepository - Persistence boundary for admin identities.
//   AdminSessionStore - Session boundary for Redis-backed sessions.
//   AdminAuthService - Coordinates bootstrap, login, current admin, admin creation, and logout.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.1 - Added transport-neutral admin auth service behavior.
// END_CHANGE_SUMMARY

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

const AdminRoleAdmin = "ADMIN"

type Admin struct {
	ID           string
	Email        string
	Name         string
	PasswordHash string
	Role         string
	IsActive     bool
	CreatedAt    string
	UpdatedAt    string
}

type CreateAdminInput struct {
	Email        string
	Name         string
	PasswordHash string
	Role         string
}

type AdminRepository interface {
	Count(ctx context.Context) (int, error)
	Create(ctx context.Context, input CreateAdminInput) (*Admin, error)
	GetByEmail(ctx context.Context, email string) (*Admin, error)
	GetByID(ctx context.Context, id string) (*Admin, error)
}

type AdminSessionStore interface {
	Create(ctx context.Context, adminID string) (string, error)
	Get(ctx context.Context, sessionID string) (string, error)
	Delete(ctx context.Context, sessionID string) error
}

var (
	ErrAdminDuplicateEmail = errors.New("admin duplicate email")
	ErrAdminNotFound       = errors.New("admin not found")
	ErrAdminAuth           = errors.New("admin authentication failed")
	ErrAdminValidation     = errors.New("admin validation failed")
)

type InitialAdminInput struct {
	Email    string
	Name     string
	Password string
}

type LoginAdminInput struct {
	Email    string
	Password string
}

type LoginAdminResult struct {
	Admin     *Admin
	SessionID string
}

type NewAdminInput struct {
	Email    string
	Name     string
	Password string
}

type AdminAuthService struct {
	repo     AdminRepository
	sessions AdminSessionStore
}

var bcryptGenerateFromPassword = bcrypt.GenerateFromPassword

func NewAdminAuthService(repo AdminRepository, sessions AdminSessionStore) *AdminAuthService {
	return &AdminAuthService{repo: repo, sessions: sessions}
}

// START_CONTRACT: SeedInitialAdmin
//
//	PURPOSE: Create the first admin from validated env-backed input only when no admins exist.
//	INPUTS: { ctx: context.Context - startup context, input: InitialAdminInput - env-backed first admin fields }
//	OUTPUTS: { bool - true when created, error - count, validation, hashing, or repository failure }
//	SIDE_EFFECTS: Reads admin count and may insert one admin identity.
//	LINKS: M-API / V-M-API.
//
// END_CONTRACT: SeedInitialAdmin
func (s *AdminAuthService) SeedInitialAdmin(ctx context.Context, input InitialAdminInput) (bool, error) {
	count, err := s.repo.Count(ctx)
	if err != nil {
		return false, err
	}
	if count > 0 {
		return false, nil
	}
	if err := validateAdminInput(input.Email, input.Name, input.Password); err != nil {
		return false, err
	}
	hash, err := hashPassword(input.Password)
	if err != nil {
		return false, err
	}
	_, err = s.repo.Create(ctx, CreateAdminInput{
		Email:        normalizeEmail(input.Email),
		Name:         strings.TrimSpace(input.Name),
		PasswordHash: hash,
		Role:         AdminRoleAdmin,
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

// START_CONTRACT: Login
//
//	PURPOSE: Validate admin credentials and create a Redis-backed session.
//	INPUTS: { ctx: context.Context - request context, input: LoginAdminInput - email and password credentials }
//	OUTPUTS: { *LoginAdminResult - public admin plus opaque session id, error - auth or infrastructure failure }
//	SIDE_EFFECTS: Reads admin identity, verifies bcrypt hash, and creates a session.
//	LINKS: M-API / V-M-API.
//
// END_CONTRACT: Login
func (s *AdminAuthService) Login(ctx context.Context, input LoginAdminInput) (*LoginAdminResult, error) {
	log := logger.FromContext(ctx)
	log.Info("[AdminAuth][login][BLOCK_VERIFY_CREDENTIALS] login attempt")
	admin, err := s.repo.GetByEmail(ctx, normalizeEmail(input.Email))
	if err != nil {
		return nil, err
	}
	if admin == nil || !admin.IsActive {
		log.Warn("[AdminAuth][login][BLOCK_VERIFY_CREDENTIALS] credential check failed")
		return nil, ErrAdminAuth
	}
	if err := bcrypt.CompareHashAndPassword([]byte(admin.PasswordHash), []byte(input.Password)); err != nil {
		log.Warn("[AdminAuth][login][BLOCK_VERIFY_CREDENTIALS] credential check failed")
		return nil, ErrAdminAuth
	}
	sessionID, err := s.sessions.Create(ctx, admin.ID)
	if err != nil {
		return nil, err
	}
	log.Info("[AdminAuth][session][BLOCK_VALIDATE_SESSION] session created", zap.String("admin_id", admin.ID))
	return &LoginAdminResult{Admin: publicAdmin(admin), SessionID: sessionID}, nil
}

// START_CONTRACT: CurrentAdmin
//
//	PURPOSE: Resolve a session id to the current active admin.
//	INPUTS: { ctx: context.Context - request context, sessionID: string - opaque session id }
//	OUTPUTS: { *Admin - public active admin or nil, error - session or repository failure }
//	SIDE_EFFECTS: Reads session store and admin repository.
//	LINKS: M-API / V-M-API.
//
// END_CONTRACT: CurrentAdmin
func (s *AdminAuthService) CurrentAdmin(ctx context.Context, sessionID string) (*Admin, error) {
	log := logger.FromContext(ctx)
	log.Debug("[AdminAuth][session][BLOCK_VALIDATE_SESSION] session lookup")
	adminID, err := s.sessions.Get(ctx, sessionID)
	if err != nil {
		log.Error("[AdminAuth][session][BLOCK_VALIDATE_SESSION] session lookup failed", zap.Error(err))
		return nil, err
	}
	if adminID == "" {
		return nil, nil
	}
	admin, err := s.repo.GetByID(ctx, adminID)
	if err != nil {
		return nil, err
	}
	if admin == nil || !admin.IsActive {
		return nil, nil
	}
	return publicAdmin(admin), nil
}

// START_CONTRACT: CreateAdmin
//
//	PURPOSE: Create a later admin only when an active admin actor is present.
//	INPUTS: { ctx: context.Context - request context, actor: *Admin - authenticated admin actor, input: NewAdminInput - new admin fields }
//	OUTPUTS: { *Admin - created public admin, error - auth, validation, hashing, duplicate, or repository failure }
//	SIDE_EFFECTS: Inserts one admin identity through the repository.
//	LINKS: M-API / V-M-API.
//
// END_CONTRACT: CreateAdmin
func (s *AdminAuthService) CreateAdmin(ctx context.Context, actor *Admin, input NewAdminInput) (*Admin, error) {
	if actor == nil || !actor.IsActive {
		return nil, ErrAdminAuth
	}
	if err := validateAdminInput(input.Email, input.Name, input.Password); err != nil {
		return nil, err
	}
	hash, err := hashPassword(input.Password)
	if err != nil {
		return nil, err
	}
	admin, err := s.repo.Create(ctx, CreateAdminInput{
		Email:        normalizeEmail(input.Email),
		Name:         strings.TrimSpace(input.Name),
		PasswordHash: hash,
		Role:         AdminRoleAdmin,
	})
	if err != nil {
		return nil, err
	}
	return publicAdmin(admin), nil
}

// START_CONTRACT: Logout
//
//	PURPOSE: Revoke a session id idempotently.
//	INPUTS: { ctx: context.Context - request context, sessionID: string - opaque session id }
//	OUTPUTS: { error - session store deletion failure }
//	SIDE_EFFECTS: Deletes one session key through the session store.
//	LINKS: M-API / V-M-API.
//
// END_CONTRACT: Logout
func (s *AdminAuthService) Logout(ctx context.Context, sessionID string) error {
	logger.FromContext(ctx).Info("[AdminAuth][logout][BLOCK_REVOKE_SESSION] logout requested")
	return s.sessions.Delete(ctx, sessionID)
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func validateAdminInput(email string, name string, password string) error {
	if normalizeEmail(email) == "" {
		return fmt.Errorf("%w: email is required", ErrAdminValidation)
	}
	if strings.TrimSpace(name) == "" {
		return fmt.Errorf("%w: name is required", ErrAdminValidation)
	}
	if len(password) < 12 {
		return fmt.Errorf("%w: password must be at least 12 characters", ErrAdminValidation)
	}
	return nil
}

func hashPassword(password string) (string, error) {
	hashed, err := bcryptGenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("hash admin password: %w", err)
	}
	return string(hashed), nil
}

func publicAdmin(admin *Admin) *Admin {
	if admin == nil {
		return nil
	}
	copy := *admin
	copy.PasswordHash = ""
	return &copy
}
