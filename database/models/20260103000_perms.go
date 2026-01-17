package models

import (
	"slices"
)

func HasPermission(perms []string, target string) bool {
	return slices.Contains(perms, target)
}
