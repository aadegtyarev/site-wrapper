//go:build linux

package desktop

import (
	"os"
	"path/filepath"
)

func InstallHost(executable string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	return writeHostManifest(filepath.Join(home, ".mozilla", "native-messaging-hosts", NativeHostName+".json"), executable)
}
