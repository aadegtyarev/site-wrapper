//go:build windows

package desktop

import (
	"os"
	"path/filepath"

	"golang.org/x/sys/windows/registry"
)

func InstallHost(executable string) error {
	dir, err := os.UserConfigDir()
	if err != nil {
		return err
	}
	path := filepath.Join(dir, "firefox-site-wrapper", "native-host.json")
	if err := writeHostManifest(path, executable); err != nil {
		return err
	}
	key, _, err := registry.CreateKey(registry.CURRENT_USER, `Software\Mozilla\NativeMessagingHosts\`+NativeHostName, registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer key.Close()
	return key.SetStringValue("", path)
}
