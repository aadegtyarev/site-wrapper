package bridge

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/aadegtyarev/firefox-site-wrapper/internal/settings"
)

const maxMessageBytes = 1024 * 1024

var ErrInvalidLaunchURL = errors.New("invalid launch URL")

type Request struct {
	Action      string         `json:"action"`
	Password    string         `json:"password,omitempty"`
	NewPassword string         `json:"new_password,omitempty"`
	Settings    settings.Draft `json:"settings,omitempty"`
}

type PublicSettings struct {
	URL         string            `json:"url"`
	ZoomPercent int               `json:"zoom_percent"`
	Mode        settings.Mode     `json:"mode"`
	Language    settings.Language `json:"language"`
}

type Response struct {
	OK         bool            `json:"ok"`
	ErrorCode  string          `json:"error_code,omitempty"`
	NeedsSetup bool            `json:"needs_setup,omitempty"`
	Settings   *PublicSettings `json:"settings,omitempty"`
	LaunchURL  string          `json:"launch_url,omitempty"`
	AppVersion string          `json:"app_version,omitempty"`
}

type Service struct {
	SettingsPath string
	LaunchPath   string
	AppVersion   string
}

func public(s settings.Settings) *PublicSettings {
	return &PublicSettings{URL: s.URL, ZoomPercent: s.ZoomPercent, Mode: s.Mode, Language: s.Language}
}

func (s Service) Handle(req Request) (Response, error) {
	switch req.Action {
	case "startup":
		launchURL, err := consumeLaunchURL(s.LaunchPath)
		if err != nil {
			if errors.Is(err, ErrInvalidLaunchURL) {
				return Response{ErrorCode: "invalid_launch_url"}, err
			}
			return Response{ErrorCode: "storage"}, err
		}
		current, err := settings.Load(s.SettingsPath)
		if errors.Is(err, os.ErrNotExist) {
			return Response{OK: true, NeedsSetup: true, LaunchURL: launchURL, AppVersion: s.AppVersion}, nil
		}
		if err != nil {
			return Response{ErrorCode: "storage"}, err
		}
		if launchURL == "" {
			launchURL = current.URL
		}
		return Response{OK: true, Settings: public(current), LaunchURL: launchURL, AppVersion: s.AppVersion}, nil
	case "unlock":
		current, err := settings.Load(s.SettingsPath)
		if err != nil {
			return Response{ErrorCode: "storage"}, err
		}
		if !settings.VerifyPassword(req.Password, current.PasswordHash) {
			return Response{ErrorCode: "incorrect_password"}, nil
		}
		return Response{OK: true, Settings: public(current), AppVersion: s.AppVersion}, nil
	case "save":
		current, err := settings.Apply(s.SettingsPath, req.Password, req.NewPassword, req.Settings)
		if err != nil {
			code := "storage"
			switch {
			case errors.Is(err, settings.ErrUnauthorized):
				code = "incorrect_password"
			case errors.Is(err, settings.ErrPasswordRequired):
				code = "password_required"
			default:
				var pathErr *os.PathError
				if !errors.As(err, &pathErr) {
					code = "invalid_settings"
				}
			}
			return Response{ErrorCode: code}, err
		}
		return Response{OK: true, Settings: public(current), AppVersion: s.AppVersion}, nil
	default:
		return Response{ErrorCode: "unknown_action"}, nil
	}
}

func ReadRequest(r io.Reader) (Request, error) {
	var length uint32
	if err := binary.Read(r, binary.LittleEndian, &length); err != nil {
		return Request{}, err
	}
	if length == 0 || length > maxMessageBytes {
		return Request{}, errors.New("native message has invalid length")
	}
	data := make([]byte, length)
	if _, err := io.ReadFull(r, data); err != nil {
		return Request{}, err
	}
	var req Request
	if err := json.Unmarshal(data, &req); err != nil {
		return Request{}, err
	}
	return req, nil
}

func WriteResponse(w io.Writer, res Response) error {
	data, err := json.Marshal(res)
	if err != nil {
		return err
	}
	if len(data) > maxMessageBytes {
		return errors.New("native response is too large")
	}
	if err := binary.Write(w, binary.LittleEndian, uint32(len(data))); err != nil {
		return err
	}
	_, err = w.Write(data)
	return err
}

func LaunchPath(settingsPath string) string {
	return filepath.Join(filepath.Dir(settingsPath), "launch-url")
}

func SetLaunchURL(path, url string) error {
	if err := settings.ValidateURL(url); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return fmt.Errorf("create launch directory: %w", err)
	}
	return os.WriteFile(path, []byte(url), 0600)
}

func consumeLaunchURL(path string) (string, error) {
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	if err := settings.ValidateURL(string(data)); err != nil {
		return "", fmt.Errorf("%w: %v", ErrInvalidLaunchURL, err)
	}
	if err := os.Remove(path); err != nil {
		return "", err
	}
	return string(data), nil
}
