package database

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"app/database/sqlite/db"
	"app/sessions"
	"app/storage/kv"
)

type SqliteSessionStore struct {
	DB      *sql.DB
	Queries *db.Queries
}

func NewSqliteSessionStore(s *Sqlite) *SqliteSessionStore {
	return &SqliteSessionStore{
		DB:      s.DB,
		Queries: s.Queries,
	}
}

func (s *SqliteSessionStore) Set(ctx context.Context, sessionID string, session *sessions.Session) error {
	return s.Queries.UpsertSession(ctx, db.UpsertSessionParams{
		SessionID:        sessionID,
		SessionUser:      session.SessionUser,
		SessionCsrfToken: session.CSRFToken,
	})
}

func (s *SqliteSessionStore) Get(ctx context.Context, sessionID string) (*sessions.Session, error) {
	session, err := s.Queries.SelectSession(ctx, sessionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, kv.ErrNotFound
		}
		return nil, err
	}

	return &sessions.Session{
		SessionUser: session.SessionUser,
		CSRFToken:   session.SessionCsrfToken,
		CreatedAt:   time.Unix(session.SessionCreatedAt, 0),
		LastUsedAt:  time.Now(),
	}, nil
}

func (s *SqliteSessionStore) Delete(ctx context.Context, sessionID string) error {
	return s.Queries.DeleteSession(ctx, sessionID)
}

func (s *SqliteSessionStore) GetElems(ctx context.Context) (map[string]sessions.Session, error) {
	dbSessions, err := s.Queries.GetSessions(ctx)
	if err != nil {
		return map[string]sessions.Session{}, err
	}

	ss := map[string]sessions.Session{}

	for _, session := range dbSessions {
		ss[session.SessionID] = sessions.Session{
			SessionUser: session.SessionUser,
			CSRFToken:   session.SessionCsrfToken,
			CreatedAt:   time.Unix(session.SessionCreatedAt, 0),
			LastUsedAt:  time.Unix(session.SessionLastUsedAt, 0),
		}
	}

	return ss, nil
}
