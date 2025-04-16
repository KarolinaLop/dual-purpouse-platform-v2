package web

import (
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"

	"github.com/KarolinaLop/dp/data"
	csessions "github.com/gin-contrib/sessions"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

// SQLiteStore is a session store that uses SQLite as the backend.
// It implements the sessions.Store interface from the gin-contrib/sessions package.
type SQLiteStore struct {
	db      *sql.DB
	authKey []byte // Key used to authenticate the cookie value using HMAC
	encKey  []byte // Key used to encrypt the cookie value
	options *sessions.Options
}

// NewSQLiteStore creates a new SQLiteStore with the given database connection and codecs.
func NewSQLiteStore(db *sql.DB, authKey []byte, encKey []byte) *SQLiteStore {
	return &SQLiteStore{
		db:      db,
		authKey: authKey,
		encKey:  encKey,
	}
}

// Get retrieves a session for the request cookie.
func (store SQLiteStore) Get(r *http.Request, name string) (*sessions.Session, error) {
	cookie, err := r.Cookie(name)
	if err != nil {
		if err == http.ErrNoCookie {
			// If no cookie is found, return a new empty session
			return store.New(r, name)
		}
		return nil, fmt.Errorf("Get: could not get request cookie: %w", err)
	}

	session := sessions.NewSession(store, name)
	// Decode the session ID from the cookie
	var c = securecookie.New(store.authKey, store.encKey)
	err = c.Decode(name, cookie.Value, &session.ID)
	if err != nil {
		return nil, fmt.Errorf("Get: could not decode cookie: %w", err)
	}

	store.applySessionOptions(session)

	user, err := data.GetSessionUser(store.db, session.ID)
	if err != nil {
		return store.New(r, name)
	}

	session.Values = map[interface{}]interface{}{
		"user_id":   user.ID,
		"user_name": user.Name,
	}

	return session, nil
}

// New creates a new session with a unique ID and default values.
func (store SQLiteStore) New(r *http.Request, name string) (*sessions.Session, error) {
	sessionID := generateSessionID()
	s := sessions.NewSession(store, name)
	s.ID = sessionID
	s.Values = map[interface{}]interface{}{
		"session_id": sessionID,
	}

	store.applySessionOptions(s)

	return s, nil
}

func (store SQLiteStore) applySessionOptions(s *sessions.Session) {
	s.Options = &sessions.Options{
		Path:     store.options.Path,
		Secure:   store.options.Secure,
		MaxAge:   store.options.MaxAge,
		HttpOnly: store.options.HttpOnly,
		SameSite: store.options.SameSite,
	}
}

func generateSessionID() string {
	b := securecookie.GenerateRandomKey(32)
	return base64.URLEncoding.EncodeToString(b)
}

func (store SQLiteStore) valid(s *sessions.Session) bool {
	if s == nil || s.ID == "" || s.Name() == "" {
		return false
	}
	if s.Options == nil {
		return false
	}
	if s.Options.MaxAge == 0 || !s.Options.HttpOnly {
		return false
	}
	return true
}

// Save saves the session data to the database and sets a secure cookie in the response.
func (store SQLiteStore) Save(r *http.Request, w http.ResponseWriter, s *sessions.Session) error {
	if !store.valid(s) {
		return errors.New("Save: invalid session data")
	}

	var c = securecookie.New(store.authKey, store.encKey)
	encoded, err := c.Encode(s.Name(), s.ID)
	if err != nil {
		return fmt.Errorf("Save: could not encode session ID: %w", err)
	}

	if err := data.CreateSession(data.DB, s.ID, s.Values["user_id"].(int)); err != nil {
		return fmt.Errorf("Save: could not save session to database: %w", err)
	}

	cookie := sessions.NewCookie(s.Name(), encoded, s.Options)
	http.SetCookie(w, cookie)

	return nil
}

// Options sets the default options for the session store.
func (store *SQLiteStore) Options(options csessions.Options) {
	store.options = &sessions.Options{
		Path:     options.Path,
		Domain:   options.Domain,
		MaxAge:   options.MaxAge,
		Secure:   options.Secure,
		HttpOnly: options.HttpOnly,
		SameSite: http.SameSite(options.SameSite),
	}
}
