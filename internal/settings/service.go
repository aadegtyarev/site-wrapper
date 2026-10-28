package settings

import (
	"errors"
	"fmt"
	"os"
)

var ErrUnauthorized = errors.New("incorrect settings password")
var ErrPasswordRequired = errors.New("new settings password is required")

type Draft struct {
	URL         string   `json:"url"`
	ZoomPercent int      `json:"zoom_percent"`
	Mode        Mode     `json:"mode"`
	Language    Language `json:"language"`
}

func Apply(path, currentPassword, newPassword string, draft Draft) (Settings, error) {
	previous, err := Load(path)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return Settings{}, fmt.Errorf("load settings: %w", err)
	}
	firstSetup := errors.Is(err, os.ErrNotExist)
	if !firstSetup && !VerifyPassword(currentPassword, previous.PasswordHash) {
		return Settings{}, ErrUnauthorized
	}
	if firstSetup && newPassword == "" {
		return Settings{}, ErrPasswordRequired
	}
	next := Settings{
		Version:      Version,
		URL:          draft.URL,
		ZoomPercent:  draft.ZoomPercent,
		Mode:         draft.Mode,
		Language:     draft.Language,
		PasswordHash: previous.PasswordHash,
	}
	if newPassword != "" {
		next.PasswordHash, err = HashPassword(newPassword)
		if err != nil {
			return Settings{}, err
		}
	}
	if err := Save(path, next); err != nil {
		return Settings{}, fmt.Errorf("save settings: %w", err)
	}
	return next, nil
}
