package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/aadegtyarev/site-wrapper/internal/bridge"
	"github.com/aadegtyarev/site-wrapper/internal/desktop"
	"github.com/aadegtyarev/site-wrapper/internal/settings"
)

var version = "dev"

func main() {
	if err := execute(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func execute(args []string) error {
	if isNativeHostInvocation(args) {
		return serveNativeHost()
	}
	if len(args) > 0 {
		switch args[0] {
		case "host":
			if len(args) != 1 {
				return errors.New("host does not accept arguments")
			}
			return serveNativeHost()
		case "version", "--version":
			fmt.Println(version)
			return nil
		case "--help", "help":
			fmt.Println("Usage: site-wrapper [run] [HTTP(S) URL] | site-wrapper version")
			return nil
		case "run":
			args = args[1:]
		}
	}
	if len(args) > 1 {
		return errors.New("expected at most one URL")
	}
	override := ""
	if len(args) == 1 {
		override = args[0]
		if err := settings.ValidateURL(override); err != nil {
			return fmt.Errorf("URL: %w", err)
		}
	}
	return runFirefox(override)
}

func isNativeHostInvocation(args []string) bool {
	return len(args) >= 2 && args[1] == desktop.ExtensionID && strings.HasSuffix(strings.ToLower(args[0]), ".json")
}

func serveNativeHost() error {
	path, err := settings.Path()
	if err != nil {
		return err
	}
	req, err := bridge.ReadRequest(os.Stdin)
	if err != nil {
		return fmt.Errorf("read native request: %w", err)
	}
	service := bridge.Service{SettingsPath: path, LaunchPath: bridge.LaunchPath(path), AppVersion: version}
	response, handleErr := service.Handle(req)
	if err := bridge.WriteResponse(os.Stdout, response); err != nil {
		return fmt.Errorf("write native response: %w", err)
	}
	if handleErr != nil {
		fmt.Fprintf(os.Stderr, "site-wrapper native request failed: %s\n", response.ErrorCode)
	}
	return nil
}

func runFirefox(override string) error {
	path, err := settings.Path()
	if err != nil {
		return err
	}
	current, err := settings.Load(path)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("read settings: %w", err)
	}
	configured := err == nil
	if configured {
		override, err = settings.EffectiveURL(current.URL, override)
		if err != nil {
			return err
		}
	}
	firefox, err := desktop.FindFirefox()
	if errors.Is(err, desktop.ErrFirefoxMissing) {
		language := settings.English
		if configured {
			language = current.Language
		}
		message := desktop.MissingFirefoxMessage(language)
		if err == desktop.ErrFirefoxMissing {
			return errors.New(message)
		}
		return fmt.Errorf("%s (%v)", message, err)
	}
	if err != nil {
		return err
	}
	profile := filepath.Join(filepath.Dir(path), "profile")
	if err := os.MkdirAll(profile, 0700); err != nil {
		return fmt.Errorf("create Firefox profile: %w", err)
	}
	language := settings.English
	if configured {
		language = current.Language
	}
	if err := desktop.WriteProfilePreferences(profile, language); err != nil {
		return fmt.Errorf("set Firefox language: %w", err)
	}
	executable, err := os.Executable()
	if err != nil {
		return err
	}
	if err := desktop.InstallHost(executable); err != nil {
		return fmt.Errorf("install Firefox native host: %w", err)
	}
	if err := desktop.InstallExtension(executable, profile); err != nil {
		return fmt.Errorf("install Firefox extension: %w", err)
	}
	launchPath := bridge.LaunchPath(path)
	if override != "" {
		if err := bridge.SetLaunchURL(launchPath, override); err != nil {
			return err
		}
	} else if err := os.Remove(launchPath); err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	args := firefoxArgs(profile, configured && current.Mode == settings.Kiosk)
	cmd := exec.Command(firefox, args...)
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("Firefox exited: %w", err)
	}
	return nil
}

func firefoxArgs(profile string, kiosk bool) []string {
	args := []string{"--new-instance", "--profile", profile}
	if kiosk {
		args = append(args, "--kiosk")
	}
	return append(args, "about:blank")
}
