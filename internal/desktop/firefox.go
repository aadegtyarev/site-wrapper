package desktop

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/aadegtyarev/site-wrapper/internal/settings"
)

const DownloadURL = "https://www.firefox.com/"

var ErrFirefoxMissing = errors.New("Firefox is not installed")

func MissingFirefoxMessage(language settings.Language) string {
	messages := map[settings.Language]string{
		settings.English: "Firefox is required. Download it from",
		settings.German:  "Firefox wird benötigt. Herunterladen von",
		settings.Kazakh:  "Firefox қажет. Мына жерден жүктеп алыңыз:",
		settings.Italian: "Firefox è necessario. Scaricalo da",
	}
	message := messages[language]
	if message == "" {
		message = messages[settings.English]
	}
	return message + " " + DownloadURL
}

func FindFirefox() (string, error) {
	return findFirefox(runtime.GOOS, os.Getenv, exec.LookPath, regularFile)
}

func regularFile(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.Mode().IsRegular()
}

func findFirefox(goos string, getenv func(string) string, lookPath func(string) (string, error), exists func(string) bool) (string, error) {
	if override := getenv("FIREFOX_BINARY"); override != "" {
		if exists(override) {
			return override, nil
		}
		return "", fmt.Errorf("%w: FIREFOX_BINARY points to a missing file: %s", ErrFirefoxMissing, override)
	}
	for _, name := range []string{"firefox", "firefox-esr", "firefox.exe"} {
		if path, err := lookPath(name); err == nil && exists(path) {
			return path, nil
		}
	}
	if goos == "windows" {
		for _, root := range []string{getenv("ProgramFiles"), getenv("ProgramFiles(x86)"), getenv("LOCALAPPDATA")} {
			if root == "" {
				continue
			}
			for _, path := range []string{
				filepath.Join(root, "Mozilla Firefox", "firefox.exe"),
				filepath.Join(root, "Programs", "Mozilla Firefox", "firefox.exe"),
			} {
				if exists(path) {
					return path, nil
				}
			}
		}
	}
	return "", ErrFirefoxMissing
}
