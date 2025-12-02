package database

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"app/database/postgres/db"
	"app/sessions"
	"app/storage/kv"
)

type PostgresSessionStore struct {
	DB      *pgxpool.Pool
	Queries *db.Queries
}

func NewPostgresSessionStore(p *Postgres) *PostgresSessionStore {
	return &PostgresSessionStore{
		DB:      p.Pool,
		Queries: p.Queries,
	}
}

func (p *PostgresSessionStore) Set(ctx context.Context, sessionID string, session sessions.Session) error {
	return p.Queries.UpsertSession(ctx, db.UpsertSessionParams{
		SessionID:        sessionID,
		SessionUser:      session.SessionUser,
		SessionCsrfToken: session.CSRFToken,
	})
}

func (p *PostgresSessionStore) Get(ctx context.Context, sessionID string) (sessions.Session, error) {
	empty := sessions.Session{}

	session, err := p.Queries.SelectSession(ctx, sessionID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return empty, kv.ErrNotFound
		}
		return empty, err
	}

	return sessions.Session{
		SessionUser: session.SessionUser,
		CSRFToken:   session.SessionCsrfToken,
		CreatedAt:   session.SessionCreatedAt.Time,
		LastUsedAt:  time.Now(),
	}, nil
}

func (p *PostgresSessionStore) Delete(ctx context.Context, sessionID string) error {
	return p.Queries.DeleteSession(ctx, sessionID)
}

func (p *PostgresSessionStore) GetElems(ctx context.Context) (map[string]sessions.Session, error) {
	dbSessions, err := p.Queries.GetSessions(ctx)
	if err != nil {
		return map[string]sessions.Session{}, err
	}

	ss := map[string]sessions.Session{}

	for _, session := range dbSessions {
		ss[session.SessionID] = sessions.Session{
			SessionUser: session.SessionUser,
			CSRFToken:   session.SessionCsrfToken,
			CreatedAt:   session.SessionCreatedAt.Time,
			LastUsedAt:  session.SessionLastUsedAt.Time,
		}
	}

	return ss, nil
}
