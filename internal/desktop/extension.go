package desktop

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func InstallExtension(executable, profile string) error {
	source := filepath.Join(filepath.Dir(executable), ExtensionID+".xpi")
	in, err := os.Open(source)
	if err != nil {
		return fmt.Errorf("open Firefox extension %s: %w", source, err)
	}
	defer in.Close()
	dir := filepath.Join(profile, "extensions")
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}
	out, err := os.CreateTemp(dir, ".sitewrap-*.xpi")
	if err != nil {
		return err
	}
	tmp := out.Name()
	defer os.Remove(tmp)
	if _, err := io.Copy(out, in); err != nil {
		out.Close()
		return err
	}
	if err := out.Sync(); err != nil {
		out.Close()
		return err
	}
	if err := out.Close(); err != nil {
		return err
	}
	return os.Rename(tmp, filepath.Join(dir, ExtensionID+".xpi"))
}
