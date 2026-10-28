package desktop

import (
	"errors"
	"path/filepath"
	"testing"
)

func TestFindFirefoxUsesExplicitBinary(t *testing.T) {
	want := filepath.Join(t.TempDir(), "firefox")
	getenv := func(key string) string {
		if key == "FIREFOX_BINARY" {
			return want
		}
		return ""
	}
	got, err := findFirefox("linux", getenv, func(string) (string, error) { t.Fatal("PATH should not be consulted"); return "", nil }, func(path string) bool { return path == want })
	if err != nil || got != want {
		t.Fatalf("got %q, %v", got, err)
	}
}

func TestFindFirefoxWindowsInstall(t *testing.T) {
	root := filepath.Join(t.TempDir(), "Program Files")
	want := filepath.Join(root, "Mozilla Firefox", "firefox.exe")
	getenv := func(key string) string {
		if key == "ProgramFiles" {
			return root
		}
		return ""
	}
	got, err := findFirefox("windows", getenv, func(string) (string, error) { return "", errors.New("missing") }, func(path string) bool { return path == want })
	if err != nil || got != want {
		t.Fatalf("got %q, %v", got, err)
	}
}

func TestFindFirefoxMissing(t *testing.T) {
	_, err := findFirefox("linux", func(string) string { return "" }, func(string) (string, error) { return "", errors.New("missing") }, func(string) bool { return false })
	if !errors.Is(err, ErrFirefoxMissing) {
		t.Fatalf("got %v, want ErrFirefoxMissing", err)
	}
}
