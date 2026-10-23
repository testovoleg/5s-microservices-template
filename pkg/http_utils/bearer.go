package httpUtils

import "strings"

// StripBearer removes a leading "Bearer " prefix (any case). Other values are returned as-is.
func StripBearer(token string) string {
	if len(token) >= 7 && strings.EqualFold(token[:7], "bearer ") {
		return token[7:]
	}
	return token
}
