package redact

import (
	"strings"
	"testing"
)

func TestNew_HasDefaultRules(t *testing.T) {
	r := New()
	if len(r.rules) == 0 {
		t.Fatal("expected default rules to be registered")
	}
}

func TestLine_BearerToken(t *testing.T) {
	r := New()
	input := `GET /api HTTP/1.1 Authorization: Bearer eyJhbGciOiJIUzI1NiJ9.payload.sig`
	out := r.Line(input)
	if strings.Contains(out, "eyJ") {
		t.Errorf("bearer token not redacted: %s", out)
	}
	if !strings.Contains(out, "[REDACTED]") {
		t.Errorf("expected [REDACTED] placeholder, got: %s", out)
	}
}

func TestLine_Password(t *testing.T) {
	r := New()
	cases := []string{
		`login password=supersecret more`,
		`login password: supersecret`,
	}
	for _, c := range cases {
		out := r.Line(c)
		if strings.Contains(out, "supersecret") {
			t.Errorf("password not redacted in %q, got: %s", c, out)
		}
	}
}

func TestLine_Email(t *testing.T) {
	r := New()
	input := `user alice@example.com logged in`
	out := r.Line(input)
	if strings.Contains(out, "alice@example.com") {
		t.Errorf("email not redacted: %s", out)
	}
	if !strings.Contains(out, "[EMAIL]") {
		t.Errorf("expected [EMAIL] placeholder, got: %s", out)
	}
}

func TestLine_NoSensitiveData(t *testing.T) {
	r := New()
	input := `2024-01-15T10:00:00Z INFO server started on :8080`
	out := r.Line(input)
	if out != input {
		t.Errorf("expected unchanged line, got: %s", out)
	}
}

func TestWithRule_CustomPattern(t *testing.T) {
	opt, err := WithRule(`\b\d{4}-\d{4}-\d{4}-\d{4}\b`, `[CARD]`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	r := New(opt)
	input := `charged card 4111-1111-1111-1111 ok`
	out := r.Line(input)
	if strings.Contains(out, "4111") {
		t.Errorf("card number not redacted: %s", out)
	}
	if !strings.Contains(out, "[CARD]") {
		t.Errorf("expected [CARD] placeholder, got: %s", out)
	}
}

func TestWithRule_InvalidPattern(t *testing.T) {
	_, err := WithRule(`[invalid`, `x`)
	if err == nil {
		t.Fatal("expected error for invalid regex")
	}
}

func TestAddRule_Runtime(t *testing.T) {
	r := New()
	if err := r.AddRule(`tok_[A-Za-z0-9]+`, `[TOKEN]`); err != nil {
		t.Fatalf("AddRule failed: %v", err)
	}
	out := r.Line(`api_key tok_abc123XYZ request`)
	if strings.Contains(out, "tok_abc123XYZ") {
		t.Errorf("custom token not redacted: %s", out)
	}
}

func TestAddRule_InvalidPattern(t *testing.T) {
	r := New()
	if err := r.AddRule(`(unclosed`, `x`); err == nil {
		t.Fatal("expected error for invalid regex")
	}
}
