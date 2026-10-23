package utils

import (
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"strings"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func Attr(span trace.Span, attributes ...string) {
	setFilled := func(name, value string) {
		if name != "" && value != "" {
			span.SetAttributes(attribute.String(name, value))
		}
	}

	if len(attributes)%2 != 0 {
		for _, attr := range attributes {
			setFilled(fmt.Sprintf("key_%d", rand.IntN(1000)), attr)
		}
	} else {
		for i := 0; i < len(attributes)-1; i += 2 {
			setFilled(attributes[i], attributes[i+1])
		}
	}

	//how to use
	//utils.Attr(span, "api_uuid", command.ApiUuid)
}

// AttrJSON marshals v to JSON and stores it under key as a single span attribute.
// Secret-looking fields are replaced with "[redacted]".
func AttrJSON(span trace.Span, key string, v any) {
	if v == nil {
		return
	}
	b, err := json.Marshal(RedactSecrets(v))
	if err != nil {
		return
	}
	Attr(span, key, string(b))
}

var secretJSONKeys = map[string]struct{}{
	"token":         {},
	"access_token":  {},
	"accesstoken":   {},
	"password":      {},
	"client_secret": {},
	"clientsecret":  {},
	"authorization": {},
	"secret":        {},
	"refreshtoken":  {},
	"refresh_token": {},
}

// RedactSecrets returns a JSON-round-tripped copy of v with secret fields masked.
func RedactSecrets(v any) any {
	if v == nil {
		return nil
	}
	b, err := json.Marshal(v)
	if err != nil {
		return v
	}
	var decoded any
	if err := json.Unmarshal(b, &decoded); err != nil {
		return v
	}
	return redactWalk(decoded)
}

func redactWalk(v any) any {
	switch t := v.(type) {
	case map[string]any:
		out := make(map[string]any, len(t))
		for k, val := range t {
			if _, secret := secretJSONKeys[strings.ToLower(k)]; secret {
				if s, ok := val.(string); ok && s != "" {
					out[k] = "[redacted]"
					continue
				}
			}
			out[k] = redactWalk(val)
		}
		return out
	case []any:
		for i := range t {
			t[i] = redactWalk(t[i])
		}
		return t
	default:
		return v
	}
}
