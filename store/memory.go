package store

import (
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

func NewMemoryStore() *sessions.CookieStore {
	secret := securecookie.GenerateRandomKey(64)
	store := sessions.NewCookieStore(secret)
	return store
}
