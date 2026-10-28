package settings

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestApplyFirstSetupAndUpdate(t *testing.T) {
	path := filepath.Join(t.TempDir(), "settings.json")
	draft := Draft{URL: "https://example.org", ZoomPercent: 100, Mode: Window, Language: English}
	if _, err := Apply(path, "", "", draft); !errors.Is(err, ErrPasswordRequired) {
		t.Fatalf("missing initial password: %v", err)
	}
	first, err := Apply(path, "", "my-long-password", draft)
	if err != nil {
		t.Fatal(err)
	}
	if !VerifyPassword("my-long-password", first.PasswordHash) {
		t.Fatal("initial password not saved")
	}
	draft.Mode = Kiosk
	if _, err := Apply(path, "wrong-password", "", draft); !errors.Is(err, ErrUnauthorized) {
		t.Fatalf("wrong password: %v", err)
	}
	unchanged, err := Load(path)
	if err != nil || unchanged.Mode != Window {
		t.Fatalf("wrong password changed settings: %+v, %v", unchanged, err)
	}
	updated, err := Apply(path, "my-long-password", "another-long-password", draft)
	if err != nil || updated.Mode != Kiosk || !VerifyPassword("another-long-password", updated.PasswordHash) {
		t.Fatalf("update failed: %+v, %v", updated, err)
	}
}

func TestApplyInvalidInputDoesNotWrite(t *testing.T) {
	path := filepath.Join(t.TempDir(), "settings.json")
	draft := Draft{URL: "https://example.org", ZoomPercent: 100, Mode: Window, Language: English}
	_, err := Apply(path, "", "my-long-password", draft)
	if err != nil {
		t.Fatal(err)
	}
	before, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	draft.URL = "file:///etc/passwd"
	if _, err := Apply(path, "my-long-password", "", draft); err == nil {
		t.Fatal("invalid URL accepted")
	}
	after, err := os.ReadFile(path)
	if err != nil || string(after) != string(before) {
		t.Fatalf("invalid input changed stored settings: %v", err)
	}
}
