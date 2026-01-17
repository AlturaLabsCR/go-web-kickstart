package config

type SessionData struct {
	// Static session data
	UserID string
	Agent  string

	// Dynamic session data
	LastUsed int64
}
