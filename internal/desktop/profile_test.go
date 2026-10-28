package desktop

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/aadegtyarev/firefox-site-wrapper/internal/settings"
)

func TestWriteProfilePreferences(t *testing.T) {
	profile := filepath.Join(t.TempDir(), "profile")
	if err := WriteProfilePreferences(profile, settings.Kazakh); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(filepath.Join(profile, "user.js"))
	if err != nil || !strings.Contains(string(data), `"kk,en"`) {
		t.Fatalf("language preference = %q, %v", data, err)
	}
	if err := WriteProfilePreferences(profile, "unknown"); err == nil {
		t.Fatal("unknown language accepted")
	}
}
