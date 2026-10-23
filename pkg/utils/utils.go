package utils

import (
	"fmt"
	"slices"
	"strings"
	"time"
	"unicode"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func GenerateUuid() string {
	return uuid.New().String()
}

func CurrentTime() time.Time {
	return time.Now()
}

func CurrentTimestamppb() *timestamppb.Timestamp {
	return timestamppb.New(CurrentTime())
}

// AppendUnique добавляет элемент в слайс, только если его там еще нет.
// type string, int, GUID и т.д. (благодаря comparable)
func AppendUnique[T comparable](slice []T, item T) []T {
	for _, v := range slice {
		if v == item {
			return slice
		}
	}
	return append(slice, item)
}

// ToRedisKey transliterates Cyrillic characters to Latin, replaces spaces
// with '-', and removes all other characters that are not letters or digits.
func ToRedisKey(s string) string {
	var b strings.Builder
	for _, r := range Transliterate(s) {
		switch {
		case r == ' ':
			b.WriteRune('-')
		case unicode.IsLetter(r) || unicode.IsDigit(r):
			b.WriteRune(r)
		}
	}
	return b.String()
}

func RemoveStringDublicates(in []string) []string {
	if len(in) < 2 {
		return in
	}

	slices.Sort(in)
	return slices.Compact(in)
}

// NormalizeStringList trims whitespace, drops empty values, and removes duplicates
// (result is sorted, see RemoveStringDublicates).
func NormalizeStringList(in []string) []string {
	out := make([]string, 0, len(in))
	for _, v := range in {
		if v = strings.TrimSpace(v); v != "" {
			out = append(out, v)
		}
	}
	return RemoveStringDublicates(out)
}

// FlattenUniqueStrings concatenates all sets into a single deduplicated, sorted slice.
func FlattenUniqueStrings(sets [][]string) []string {
	var out []string
	for _, set := range sets {
		out = append(out, set...)
	}
	return RemoveStringDublicates(out)
}

// CoalescePatchString returns *next unless it is nil or equal to a placeholder
// (e.g. OpenAPI default "string").
func CoalescePatchString(current string, next *string, placeholders ...string) string {
	if next == nil {
		return current
	}
	for _, p := range placeholders {
		if *next == p {
			return current
		}
	}
	return *next
}

// FormatDurationHuman renders a duration as a short string: "45s", "12min", "1h 20min".
func FormatDurationHuman(d time.Duration) string {
	if d < 0 {
		d = 0
	}
	switch {
	case d < time.Minute:
		return fmt.Sprintf("%ds", int64(d.Seconds()))
	case d < time.Hour:
		return fmt.Sprintf("%dmin", int64(d.Minutes()))
	default:
		hours := int64(d.Hours())
		minutes := int64(d.Minutes()) % 60
		if minutes == 0 {
			return fmt.Sprintf("%dh", hours)
		}
		return fmt.Sprintf("%dh %dmin", hours, minutes)
	}
}

func formatServiceName(input string) string {
	result := []string{"5S"}
	for _, w := range strings.Split(input, "-") {
		if w != "" {
			result = append(result, capitalize(w))
		}
	}
	return strings.Join(result, " ")
}

func capitalize(s string) string {
	if s == "" {
		return s
	}
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	for i := 1; i < len(runes); i++ {
		runes[i] = unicode.ToLower(runes[i])
	}
	return string(runes)
}
