package database

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"app/database/postgres/db"
	"app/sessions"
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
	session, err := p.Queries.SelectSession(ctx, sessionID)
	if err != nil {
		return sessions.Session{}, err
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
