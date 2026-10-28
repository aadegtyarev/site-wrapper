package desktop

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/aadegtyarev/site-wrapper/internal/settings"
)

func WriteProfilePreferences(profile string, language settings.Language) error {
	if language != settings.English && language != settings.German && language != settings.Kazakh && language != settings.Italian {
		return errors.New("unsupported language")
	}
	requested := string(language)
	if language != settings.English {
		requested += ",en"
	}
	data := []byte(fmt.Sprintf("user_pref(\"intl.accept_languages\", %q);\nuser_pref(\"browser.shell.checkDefaultBrowser\", false);\n", requested))
	if err := os.MkdirAll(profile, 0700); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(profile, "user.js"), data, 0600)
}
