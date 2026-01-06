package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"time"

	"task-flow/internal/handler/users"

	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	UserRepo users.Repo
	RefreshRepo RefreshRepo
	
	JWT JWT
	AccessTTL time.Duration
	RefreshTTL time.Duration

	CookieDomain string
	CookieSecure int
}

func newRefresh(ttl time.Duration) (plain string, hash []byte, expUnix int64, err error ) {
	b := make([]byte, 32)

	if _ , err := rand.Read(b); err != nil {
		return "", nil, 0, err
	}

	plain = base64.RawURLEncoding.EncodeToString(b)
	hash = hashRefresh(plain)
	expUnix = time.Now().Add(ttl).Unix()

	return plain, hash, expUnix, nil
}

func hashRefresh(plain string) []byte {
	sum := sha256.Sum256([]byte(plain))

	return sum[:]
}

func (s *Service) Login(ctx context.Context, email, password string) (access, refresh string, err error) {
	u, ok, err := s.UserRepo.FindByEmail(ctx, email)

	if err != nil {
		return "", "", err
	}

	if !ok {
		return "", "", errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword(u.PassHas, []byte(password)); err != nil {
		return "", "", errors.New("invalid credentials")
	}

	access, err = s.JWT.SignHS256(u.ID, s.AccessTTL)

	if err != nil {
		return "", "", err
	} 

	refresh, hash, expUnix, err := newRefresh(s.RefreshTTL)

	if err != nil {
		return "", "", err
	}

	if err := s.RefreshRepo.Insert(ctx, u.ID, hash, expUnix); err != nil {
		return "", "", nil
	}

	return access, refresh, nil
}

func (s *Service) Refresh(ctx context.Context, refreshPlain string) (newAccess, refreshToken string, err error) {
	hash := hashRefresh(refreshPlain)

	userID, ok, err := s.RefreshRepo.FindUserIDByHash(ctx, hash)
	
	if err != nil {
		return "", "", err
	}

	if !ok {
		return "", "", errors.New("invalid refresh")
	}

	// rotation: revoke old -> issue new
	if err := s.RefreshRepo.RevokeByHash(ctx, hash); err != nil {
		return "", "", err
	}

	newAccess, err = s.JWT.SignHS256(userID, s.AccessTTL)
	
	if err != nil {
		return "", "", err
	}

	refreshToken, newHash, expUnix, err := newRefresh(s.RefreshTTL)

	if err != nil {
		return "", "", err
	}

	if err := s.RefreshRepo.Insert(ctx, userID, newHash, expUnix); err != nil {
		return "", "", err
	}

	return newAccess, refreshToken, nil
}

func (s *Service) Logout(ctx context.Context, refreshPlain string) error {
	if refreshPlain == "" {
		return nil
	}
	return s.RefreshRepo.RevokeByHash(ctx, hashRefresh(refreshPlain))
}
