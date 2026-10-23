package v1

import (
	"html/template"
	"strings"
	"testing"
)

// The Swagger page is a template in a separate file, so a typo in it would otherwise
// only surface as a panic in swaggerHandler at runtime. Render it here instead.
func TestSwaggerUITemplateRenders(t *testing.T) {
	tpl, err := template.New("swagger_ui").Parse(swaggerUIHTML)
	if err != nil {
		t.Fatalf("parse swagger_ui.html: %v", err)
	}

	var buf strings.Builder
	if err := tpl.Execute(&buf, swaggerUIData{AppName: "App"}); err != nil {
		t.Fatalf("execute swagger_ui.html: %v", err)
	}

	out := buf.String()
	if strings.Contains(out, "{{") {
		t.Error("rendered page still contains template actions")
	}
	for _, want := range []string{"<title>App — Swagger UI</title>", "SwaggerUIBundle("} {
		if !strings.Contains(out, want) {
			t.Errorf("rendered page misses %q", want)
		}
	}
}
