package settings

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidateURL(t *testing.T) {
	for _, tc := range []struct {
		url string
		ok  bool
	}{
		{"https://example.org/path?q=1", true},
		{"http://localhost:3000/", true},
		{"file:///etc/passwd", false},
		{"javascript:alert(1)", false},
		{"https://user:pass@example.org", false},
		{"https://", false},
		{"https://example.org\n--kiosk", false},
		{"https://example.org:0/", false},
		{"https://example.org:65536/", false},
	} {
		err := ValidateURL(tc.url)
		if (err == nil) != tc.ok {
			t.Errorf("ValidateURL(%q) = %v, want ok=%v", tc.url, err, tc.ok)
		}
	}
}

func TestSaveLoadAndPermissions(t *testing.T) {
	path := filepath.Join(t.TempDir(), "settings.json")
	want := Settings{Version: Version, URL: "https://example.org", ZoomPercent: 125, Mode: Kiosk, Language: Kazakh, PasswordHash: "hash"}
	if err := Save(path, want); err != nil {
		t.Fatal(err)
	}
	got, err := Load(path)
	if err != nil || got != want {
		t.Fatalf("Load() = %+v, %v; want %+v", got, err, want)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0600 {
		t.Fatalf("mode = %o, want 600", info.Mode().Perm())
	}
}

func TestDecodeRejectsExtraData(t *testing.T) {
	data := []byte(`{"version":1,"url":"https://example.org","zoom_percent":100,"mode":"window","language":"en","password_hash":"hash","unexpected":true}`)
	if _, err := Decode(data); err == nil {
		t.Fatal("expected unknown field error")
	}
	data = []byte(`{"version":1,"url":"https://example.org","zoom_percent":100,"mode":"window","language":"en","password_hash":"hash"} {}`)
	if _, err := Decode(data); err == nil {
		t.Fatal("expected trailing data error")
	}
}

func TestEffectiveURL(t *testing.T) {
	got, err := EffectiveURL("https://saved.example", "https://override.example")
	if err != nil || got != "https://override.example" {
		t.Fatalf("got %q, %v", got, err)
	}
	got, err = EffectiveURL("https://saved.example", "")
	if err != nil || got != "https://saved.example" {
		t.Fatalf("got %q, %v", got, err)
	}
	if _, err := EffectiveURL("https://saved.example", "file:///etc/passwd"); err == nil {
		t.Fatal("expected invalid override error")
	}
}
