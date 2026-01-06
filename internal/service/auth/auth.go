package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"time"

	"task-flow/internal/pkg/jwt"
	"task-flow/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	UserRepo         repository.UserRepo
	RefreshTokenRepo repository.RefreshTokenRepo

	JWT        *jwt.JWT
	AccessTTL  time.Duration
	RefreshTTL time.Duration
}

func NewService(
	userRepo repository.UserRepo,
	refreshRepo repository.RefreshTokenRepo,
	jwtInstance *jwt.JWT,
	accessTTL, refreshTTL time.Duration,
) *Service {
	return &Service{
		UserRepo:         userRepo,
		RefreshTokenRepo: refreshRepo,
		JWT:              jwtInstance,
		AccessTTL:        accessTTL,
		RefreshTTL:       refreshTTL,
	}
}

func newRefreshToken(ttl time.Duration) (plain string, hash []byte, expUnix int64, err error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", nil, 0, err
	}

	plain = base64.RawURLEncoding.EncodeToString(b)
	hash = hashToken(plain)
	expUnix = time.Now().Add(ttl).Unix()

	return plain, hash, expUnix, nil
}

func hashToken(plain string) []byte {
	sum := sha256.Sum256([]byte(plain))
	return sum[:]
}

func (s *Service) Login(ctx context.Context, email, password string) (access, refresh string, err error) {
	user, found, err := s.UserRepo.FindByEmail(ctx, email)
	if err != nil {
		return "", "", err
	}
	if !found {
		return "", "", errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		return "", "", errors.New("invalid credentials")
	}

	access, err = s.JWT.Sign(user.ID, s.AccessTTL)
	if err != nil {
		return "", "", err
	}

	refresh, hash, expUnix, err := newRefreshToken(s.RefreshTTL)
	if err != nil {
		return "", "", err
	}

	if err := s.RefreshTokenRepo.Insert(ctx, user.ID, hash, expUnix); err != nil {
		return "", "", err
	}

	return access, refresh, nil
}

func (s *Service) Refresh(ctx context.Context, refreshPlain string) (newAccess, newRefresh string, err error) {
	hash := hashToken(refreshPlain)

	userID, found, err := s.RefreshTokenRepo.FindUserIDByHash(ctx, hash)
	if err != nil {
		return "", "", err
	}
	if !found {
		return "", "", errors.New("invalid refresh token")
	}

	// Token rotation: revoke old, issue new
	if err := s.RefreshTokenRepo.RevokeByHash(ctx, hash); err != nil {
		return "", "", err
	}

	newAccess, err = s.JWT.Sign(userID, s.AccessTTL)
	if err != nil {
		return "", "", err
	}

	newRefresh, newHash, expUnix, err := newRefreshToken(s.RefreshTTL)
	if err != nil {
		return "", "", err
	}

	if err := s.RefreshTokenRepo.Insert(ctx, userID, newHash, expUnix); err != nil {
		return "", "", err
	}

	return newAccess, newRefresh, nil
}

func (s *Service) Logout(ctx context.Context, refreshPlain string) error {
	if refreshPlain == "" {
		return nil
	}
	return s.RefreshTokenRepo.RevokeByHash(ctx, hashToken(refreshPlain))
}
