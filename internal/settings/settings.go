package settings

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const Version = 1

type Mode string
type Language string

const (
	Window  Mode     = "window"
	Kiosk   Mode     = "kiosk"
	English Language = "en"
	German  Language = "de"
	Kazakh  Language = "kk"
	Italian Language = "it"
)

type Settings struct {
	Version      int      `json:"version"`
	URL          string   `json:"url"`
	ZoomPercent  int      `json:"zoom_percent"`
	Mode         Mode     `json:"mode"`
	Language     Language `json:"language"`
	PasswordHash string   `json:"password_hash"`
}

func Defaults() Settings {
	return Settings{Version: Version, ZoomPercent: 100, Mode: Window, Language: English}
}

func ValidateURL(raw string) error {
	if raw == "" || strings.ContainsAny(raw, "\x00\r\n\t") {
		return errors.New("invalid URL")
	}
	u, err := url.Parse(raw)
	if err != nil || u == nil || (u.Scheme != "http" && u.Scheme != "https") || u.Hostname() == "" || u.User != nil {
		return errors.New("expected an absolute HTTP(S) URL without credentials")
	}
	if u.Opaque != "" || strings.HasSuffix(u.Hostname(), ".") || strings.Contains(u.Host, " ") {
		return errors.New("invalid URL host")
	}
	if u.Port() != "" {
		// url.Parse validates port syntax but not range.
		port, err := strconv.Atoi(u.Port())
		if err != nil || port < 1 || port > 65535 {
			return errors.New("invalid URL port")
		}
	}
	return nil
}

func (s Settings) Validate() error {
	if s.Version != Version {
		return fmt.Errorf("unsupported settings version %d", s.Version)
	}
	if err := ValidateURL(s.URL); err != nil {
		return fmt.Errorf("url: %w", err)
	}
	if s.ZoomPercent < 30 || s.ZoomPercent > 500 {
		return errors.New("zoom_percent must be between 30 and 500")
	}
	if s.Mode != Window && s.Mode != Kiosk {
		return errors.New("mode must be window or kiosk")
	}
	if s.Language != English && s.Language != German && s.Language != Kazakh && s.Language != Italian {
		return errors.New("language must be en, de, kk or it")
	}
	if s.PasswordHash == "" {
		return errors.New("password_hash is required")
	}
	return nil
}

func Decode(data []byte) (Settings, error) {
	var s Settings
	d := json.NewDecoder(bytes.NewReader(data))
	d.DisallowUnknownFields()
	if err := d.Decode(&s); err != nil {
		return Settings{}, err
	}
	if err := d.Decode(new(any)); err != io.EOF {
		return Settings{}, errors.New("unexpected data after settings")
	}
	if err := s.Validate(); err != nil {
		return Settings{}, err
	}
	return s, nil
}

func Load(path string) (Settings, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Settings{}, err
	}
	return Decode(data)
}

func Save(path string, s Settings) error {
	if err := s.Validate(); err != nil {
		return err
	}
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}
	f, err := os.CreateTemp(filepath.Dir(path), ".settings-*")
	if err != nil {
		return err
	}
	tmp := f.Name()
	defer os.Remove(tmp)
	if err := f.Chmod(0600); err != nil {
		f.Close()
		return err
	}
	if _, err := f.Write(data); err != nil {
		f.Close()
		return err
	}
	if err := f.Sync(); err != nil {
		f.Close()
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

func Path() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "site-wrapper", "settings.json"), nil
}

func EffectiveURL(saved, override string) (string, error) {
	if override == "" {
		override = saved
	}
	if err := ValidateURL(override); err != nil {
		return "", err
	}
	return override, nil
}
