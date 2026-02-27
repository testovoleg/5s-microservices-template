package utils

import (
	"fmt"
	"math/rand/v2"

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
