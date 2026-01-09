package auth

import (
	"context"
	"errors"
	"testing"
	"time"

	"task-flow/internal/model"
	"task-flow/internal/pkg/jwt"

	"golang.org/x/crypto/bcrypt"
)

// ============================================
// MOCK REPOSITORIES
// ============================================

type mockUserRepo struct {
	users     map[string]model.User
	createErr error
	findErr   error
}

func newMockUserRepo() *mockUserRepo {
	return &mockUserRepo{
		users: make(map[string]model.User),
	}
}

func (m *mockUserRepo) Create(ctx context.Context, user model.User) error {
	if m.createErr != nil {
		return m.createErr
	}
	m.users[user.Email] = user
	return nil
}

func (m *mockUserRepo) FindByID(ctx context.Context, id string) (model.User, bool, error) {
	if m.findErr != nil {
		return model.User{}, false, m.findErr
	}
	for _, u := range m.users {
		if u.ID == id {
			return u, true, nil
		}
	}
	return model.User{}, false, nil
}

func (m *mockUserRepo) FindByEmail(ctx context.Context, email string) (model.User, bool, error) {
	if m.findErr != nil {
		return model.User{}, false, m.findErr
	}
	user, found := m.users[email]
	return user, found, nil
}

type mockRefreshRepo struct {
	tokens    map[string]string // hash -> userID
	insertErr error
	findErr   error
	revokeErr error
}

func newMockRefreshRepo() *mockRefreshRepo {
	return &mockRefreshRepo{
		tokens: make(map[string]string),
	}
}

func (m *mockRefreshRepo) Insert(ctx context.Context, userID string, tokenHash []byte, expUnix int64) error {
	if m.insertErr != nil {
		return m.insertErr
	}
	m.tokens[string(tokenHash)] = userID
	return nil
}

func (m *mockRefreshRepo) FindUserIDByHash(ctx context.Context, tokenHash []byte) (string, bool, error) {
	if m.findErr != nil {
		return "", false, m.findErr
	}
	userID, found := m.tokens[string(tokenHash)]
	return userID, found, nil
}

func (m *mockRefreshRepo) RevokeByHash(ctx context.Context, tokenHash []byte) error {
	if m.revokeErr != nil {
		return m.revokeErr
	}
	delete(m.tokens, string(tokenHash))
	return nil
}

// ============================================
// HELPER
// ============================================

func newTestService(userRepo *mockUserRepo, refreshRepo *mockRefreshRepo) *Service {
	return NewService(
		userRepo,
		refreshRepo,
		jwt.New([]byte("test-secret")),
		15*time.Minute,
		7*24*time.Hour,
	)
}

func hashPassword(password string) []byte {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return hash
}

// ============================================
// TEST LOGIN
// ============================================

func TestLogin_Success(t *testing.T) {
	userRepo := newMockUserRepo()
	userRepo.users["test@mail.com"] = model.User{
		ID:       "user-1",
		Email:    "test@mail.com",
		PassHash: hashPassword("password123"),
	}

	svc := newTestService(userRepo, newMockRefreshRepo())

	access, refresh, err := svc.Login(context.Background(), "test@mail.com", "password123")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if access == "" {
		t.Error("expected access token, got empty")
	}
	if refresh == "" {
		t.Error("expected refresh token, got empty")
	}
}

func TestLogin_WrongEmail(t *testing.T) {
	userRepo := newMockUserRepo()
	svc := newTestService(userRepo, newMockRefreshRepo())

	_, _, err := svc.Login(context.Background(), "notfound@mail.com", "password123")

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "invalid credentials" {
		t.Errorf("expected 'invalid credentials', got '%s'", err.Error())
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	userRepo := newMockUserRepo()
	userRepo.users["test@mail.com"] = model.User{
		ID:       "user-1",
		Email:    "test@mail.com",
		PassHash: hashPassword("correctpassword"),
	}

	svc := newTestService(userRepo, newMockRefreshRepo())

	_, _, err := svc.Login(context.Background(), "test@mail.com", "wrongpassword")

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "invalid credentials" {
		t.Errorf("expected 'invalid credentials', got '%s'", err.Error())
	}
}

func TestLogin_DatabaseError(t *testing.T) {
	userRepo := newMockUserRepo()
	userRepo.findErr = errors.New("database connection failed")

	svc := newTestService(userRepo, newMockRefreshRepo())

	_, _, err := svc.Login(context.Background(), "test@mail.com", "password123")

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "database connection failed" {
		t.Errorf("expected 'database connection failed', got '%s'", err.Error())
	}
}

// ============================================
// TEST REGISTER
// ============================================

func TestRegister_Success(t *testing.T) {
	userRepo := newMockUserRepo()
	svc := newTestService(userRepo, newMockRefreshRepo())

	err := svc.Register(context.Background(), "new@mail.com", "password123")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Verify user was created
	user, found := userRepo.users["new@mail.com"]
	if !found {
		t.Fatal("expected user to be created")
	}
	if user.Email != "new@mail.com" {
		t.Errorf("expected email 'new@mail.com', got '%s'", user.Email)
	}
}

func TestRegister_EmailAlreadyExists(t *testing.T) {
	userRepo := newMockUserRepo()
	userRepo.users["existing@mail.com"] = model.User{
		ID:    "user-1",
		Email: "existing@mail.com",
	}

	svc := newTestService(userRepo, newMockRefreshRepo())

	err := svc.Register(context.Background(), "existing@mail.com", "password123")

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "email already registered" {
		t.Errorf("expected 'email already registered', got '%s'", err.Error())
	}
}

func TestRegister_DatabaseError(t *testing.T) {
	userRepo := newMockUserRepo()
	userRepo.createErr = errors.New("database error")

	svc := newTestService(userRepo, newMockRefreshRepo())

	err := svc.Register(context.Background(), "new@mail.com", "password123")

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// ============================================
// TEST REFRESH
// ============================================

func TestRefresh_Success(t *testing.T) {
	userRepo := newMockUserRepo()
	refreshRepo := newMockRefreshRepo()

	// Simulate existing refresh token
	tokenPlain := "test-refresh-token"
	tokenHash := hashToken(tokenPlain)
	refreshRepo.tokens[string(tokenHash)] = "user-1"

	svc := newTestService(userRepo, refreshRepo)

	newAccess, newRefresh, err := svc.Refresh(context.Background(), tokenPlain)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if newAccess == "" {
		t.Error("expected new access token, got empty")
	}
	if newRefresh == "" {
		t.Error("expected new refresh token, got empty")
	}

	// Old token should be revoked
	if _, found := refreshRepo.tokens[string(tokenHash)]; found {
		t.Error("expected old token to be revoked")
	}
}

func TestRefresh_InvalidToken(t *testing.T) {
	svc := newTestService(newMockUserRepo(), newMockRefreshRepo())

	_, _, err := svc.Refresh(context.Background(), "invalid-token")

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "invalid refresh token" {
		t.Errorf("expected 'invalid refresh token', got '%s'", err.Error())
	}
}

// ============================================
// TEST LOGOUT
// ============================================

func TestLogout_Success(t *testing.T) {
	refreshRepo := newMockRefreshRepo()
	tokenPlain := "test-refresh-token"
	tokenHash := hashToken(tokenPlain)
	refreshRepo.tokens[string(tokenHash)] = "user-1"

	svc := newTestService(newMockUserRepo(), refreshRepo)

	err := svc.Logout(context.Background(), tokenPlain)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Token should be revoked
	if _, found := refreshRepo.tokens[string(tokenHash)]; found {
		t.Error("expected token to be revoked")
	}
}

func TestLogout_EmptyToken(t *testing.T) {
	svc := newTestService(newMockUserRepo(), newMockRefreshRepo())

	err := svc.Logout(context.Background(), "")
	if err != nil {
		t.Fatalf("expected no error for empty token, got %v", err)
	}
}
