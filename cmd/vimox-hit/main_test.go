package main

import (
	"strings"
	"testing"
)

func TestFlagValue(t *testing.T) {
	if got := flagValue([]string{"--name", "campaign"}, "--name"); got != "campaign" {
		t.Fatalf("flagValue() = %q", got)
	}
}

func TestSlugRemovesPathSeparators(t *testing.T) {
	if got := slug("../campaign:one"); got != "campaign-one" {
		t.Fatalf("slug() = %q", got)
	}
}

func TestSanitizeRedactsCredentialsAndMedia(t *testing.T) {
	input := map[string]any{
		"authorization": "Bearer secret-token",
		"prompt":        "see https://cdn.example/a.png?sign=secret",
		"payload":       "data:image/png;base64,AAAA",
	}
	encoded, err := marshalDiagnostic(input)
	if err != nil {
		t.Fatal(err)
	}
	value := string(encoded)
	for _, secret := range []string{"secret-token", "sign=secret", "AAAA"} {
		if strings.Contains(value, secret) {
			t.Fatalf("diagnostic leaked %q: %s", secret, value)
		}
	}
	if !strings.Contains(value, "[REDACTED]") || !strings.Contains(value, "cdn.example/a.png") {
		t.Fatalf("expected sanitized diagnostic, got %s", value)
	}
}
