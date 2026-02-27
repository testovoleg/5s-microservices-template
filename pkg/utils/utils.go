package utils

import (
	"context"
	"slices"
	"strings"
	"time"
	"unicode"

	"github.com/google/uuid"
)

func GenerateUuid() string {
	return uuid.New().String()
}

func Delay(ctx context.Context, d time.Duration) error {
	select {
	case <-time.After(d):
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}

	//how to use
	// if err := utils.Delay(ctx, 6*time.Second); err != nil {
	// 	c.log.Info(err)
	// }
}

func RemoveStringDublicates(in []string) []string {
	if len(in) < 2 {
		return in
	}

	slices.Sort(in)
	return slices.Compact(in)
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
