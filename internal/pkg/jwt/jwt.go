package jwt

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"time"
)

type JWT struct {
	Secret []byte
}

func New(secret []byte) *JWT {
	return &JWT{Secret: secret}
}

func b64url(b []byte) string {
	return base64.RawURLEncoding.EncodeToString(b)
}

func b64urldecode(s string) ([]byte, error) {
	return base64.RawURLEncoding.DecodeString(s)
}

func signHMACSHA256(msg, secret []byte) []byte {
	m := hmac.New(sha256.New, secret)
	m.Write(msg)
	return m.Sum(nil)
}

func (j *JWT) Sign(sub string, ttl time.Duration) (string, error) {
	header := map[string]any{"alg": "HS256", "typ": "JWT"}
	now := time.Now().UTC()
	payload := map[string]any{
		"sub": sub,
		"iat": now.Unix(),
		"exp": now.Add(ttl).Unix(),
	}

	hb, err := json.Marshal(header)
	if err != nil {
		return "", err
	}

	pb, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	h64 := b64url(hb)
	p64 := b64url(pb)
	msg := h64 + "." + p64

	sig := signHMACSHA256([]byte(msg), j.Secret)
	return msg + "." + b64url(sig), nil
}

func (j *JWT) Verify(token string) (sub string, err error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return "", errors.New("invalid token")
	}

	msg := parts[0] + "." + parts[1]
	want := signHMACSHA256([]byte(msg), j.Secret)
	got, err := b64urldecode(parts[2])
	if err != nil {
		return "", errors.New("invalid token")
	}

	if !hmac.Equal(got, want) {
		return "", errors.New("invalid signature")
	}

	payloadBytes, err := b64urldecode(parts[1])
	if err != nil {
		return "", errors.New("invalid token")
	}

	var payload struct {
		Sub string `json:"sub"`
		Exp int64  `json:"exp"`
	}

	if err := json.Unmarshal(payloadBytes, &payload); err != nil {
		return "", errors.New("invalid token")
	}

	if payload.Sub == "" {
		return "", errors.New("missing sub")
	}

	if time.Now().UTC().Unix() >= payload.Exp {
		return "", errors.New("expired")
	}

	return payload.Sub, nil
}
