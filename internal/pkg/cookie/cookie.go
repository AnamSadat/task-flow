package cookie

import (
	"net/http"
	"time"
)

type CookieManager struct {
	CookieDomain string
	CookieSecure bool
	RefreshTTL   time.Duration
}

func NewCookie(domain string, secure bool, refreshTTL time.Duration) *CookieManager {
	return &CookieManager{
		CookieDomain: domain,
		CookieSecure: secure,
		RefreshTTL:   refreshTTL,
	}
}

func (m *CookieManager) SetRefreshCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    token,
		Path:     "/auth",
		Domain:   m.CookieDomain,
		HttpOnly: true,
		Secure:   m.CookieSecure,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(m.RefreshTTL),
	})
}

func (m *CookieManager) ClearRefreshCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/auth",
		Domain:   m.CookieDomain,
		HttpOnly: true,
		Secure:   m.CookieSecure,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})
}
