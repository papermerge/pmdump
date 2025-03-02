package utils

import (
	"strings"

	"github.com/google/uuid"
)

/* Returns UUID as string without hyphens */
func UUID2STR(input uuid.UUID) string {
	str := input.String()
	ret := strings.ReplaceAll(str, "-", "")
	return ret
}

/* Returns UUID as string without hyphens */
func AnyUUID2STR(input any) string {
	uid := input.(uuid.UUID)
	ret := UUID2STR(uid)

	return ret
}
