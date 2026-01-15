package models

import "time"

const (
	UserCacheScopePrefix = "user:"
)

type User struct {
	ID        string
	CreatedAt time.Time
}
