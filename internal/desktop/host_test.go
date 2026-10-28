package desktop

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestWriteHostManifest(t *testing.T) {
	path := filepath.Join(t.TempDir(), "host.json")
	if err := writeHostManifest(path, "/opt/site-wrapper/site-wrapper"); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var manifest hostManifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		t.Fatal(err)
	}
	if manifest.Name != NativeHostName || manifest.Path != "/opt/site-wrapper/site-wrapper" || len(manifest.AllowedExtensions) != 1 || manifest.AllowedExtensions[0] != ExtensionID {
		t.Fatalf("unexpected manifest: %+v", manifest)
	}
	if err := writeHostManifest(path, "relative/path"); err == nil {
		t.Fatal("relative executable accepted")
	}
}
