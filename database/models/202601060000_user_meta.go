package models

type UserMeta struct {
	ID      string
	Created int64
	Name    string
	Perms   []string
}
