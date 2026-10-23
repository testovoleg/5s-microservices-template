package utils

import (
	"strings"
	"unicode"
)

// NormalizeName collapses whitespace and uppercases for catalog / alias lookup.
func NormalizeName(name string) string {
	fields := strings.Fields(strings.TrimSpace(name))
	if len(fields) == 0 {
		return ""
	}
	return strings.ToUpper(strings.Join(fields, " "))
}

// ContainsCyrillic reports whether s has any Cyrillic letter.
func ContainsCyrillic(s string) bool {
	for _, r := range s {
		if unicode.In(r, unicode.Cyrillic) {
			return true
		}
	}
	return false
}

// NamesMatch reports whether two names are the same after NormalizeName.
func NamesMatch(a, b string) bool {
	na := NormalizeName(a)
	nb := NormalizeName(b)
	if na == "" || nb == "" {
		return false
	}
	return na == nb
}

// PickUniqueNameMatch returns the only candidate whose name matches want.
func PickUniqueNameMatch[T any](want string, items []T, nameFn func(T) string) (T, bool) {
	var zero T
	var found T
	n := 0
	for _, item := range items {
		if NamesMatch(want, nameFn(item)) {
			found = item
			n++
		}
	}
	if n == 1 {
		return found, true
	}
	return zero, false
}
