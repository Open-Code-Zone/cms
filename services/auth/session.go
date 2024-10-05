package auth

import "github.com/gorilla/sessions"

const (
	SessionName = "session"
)

type SessionOptions struct {
	CookiesKey string
	MaxAge     int
	HttpOnly   bool // Should be true in production
	Secure     bool // Should be true in production
}

// cockie store is not used since it doesn't able to store cookie of larger size
func NewCookieStore(opts SessionOptions) *sessions.CookieStore {
	store := sessions.NewCookieStore([]byte(opts.CookiesKey))

	store.MaxAge(opts.MaxAge)
	store.Options.Path = "/"
	store.Options.HttpOnly = opts.HttpOnly
	store.Options.Secure = opts.Secure

	return store
}

func NewFileSystemStore(opts SessionOptions) *sessions.FilesystemStore {
	store := sessions.NewFilesystemStore("", []byte(opts.CookiesKey))
  store.MaxLength(8192)

	store.MaxAge(opts.MaxAge)
	store.Options.Path = "/"
	store.Options.HttpOnly = opts.HttpOnly
	store.Options.Secure = opts.Secure

	return store
}
