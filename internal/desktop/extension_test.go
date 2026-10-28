package desktop

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInstallExtensionCopiesPackageIntoProfile(t *testing.T) {
	dir := t.TempDir()
	executable := filepath.Join(dir, "sitewrap")
	source := filepath.Join(dir, ExtensionID+".xpi")
	if err := os.WriteFile(source, []byte("signed-xpi-placeholder"), 0600); err != nil {
		t.Fatal(err)
	}
	profile := filepath.Join(dir, "profile")
	if err := InstallExtension(executable, profile); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(filepath.Join(profile, "extensions", ExtensionID+".xpi"))
	if err != nil || string(data) != "signed-xpi-placeholder" {
		t.Fatalf("installed extension = %q, %v", data, err)
	}
}

func TestInstallExtensionReportsMissingPackage(t *testing.T) {
	if err := InstallExtension(filepath.Join(t.TempDir(), "sitewrap"), t.TempDir()); err == nil {
		t.Fatal("missing XPI was not reported")
	}
}
