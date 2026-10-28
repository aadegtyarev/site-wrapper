package settings

import (
	"strings"
	"testing"
)

func TestPasswordRoundTrip(t *testing.T) {
	hash, err := HashPassword("a-long-passphrase")
	if err != nil {
		t.Fatal(err)
	}
	if !VerifyPassword("a-long-passphrase", hash) {
		t.Fatal("correct password was rejected")
	}
	if VerifyPassword("wrong-password", hash) {
		t.Fatal("wrong password was accepted")
	}
	other, err := HashPassword("a-long-passphrase")
	if err != nil {
		t.Fatal(err)
	}
	if other == hash {
		t.Fatal("password hashes reused a salt")
	}
}

func TestPasswordRejectsShortAndMalformed(t *testing.T) {
	if _, err := HashPassword("too-short"); err == nil {
		t.Fatal("short password accepted")
	}
	for _, hash := range []string{"", "plain", "$argon2id$v=19$m=999999999,t=3,p=1$abc$abc", strings.Repeat("x", 1000)} {
		if VerifyPassword("a-long-passphrase", hash) {
			t.Fatalf("malformed verifier accepted: %q", hash)
		}
	}
}
