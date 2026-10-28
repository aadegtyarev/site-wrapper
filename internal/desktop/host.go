package desktop

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

const NativeHostName = "io.github.aadegtyarev.sitewrapper"
const ExtensionID = "site-wrapper@aadegtyarev.github.io"

type hostManifest struct {
	Name              string   `json:"name"`
	Description       string   `json:"description"`
	Path              string   `json:"path"`
	Type              string   `json:"type"`
	AllowedExtensions []string `json:"allowed_extensions"`
}

func writeHostManifest(path, executable string) error {
	if !filepath.IsAbs(executable) {
		return errors.New("native host executable path must be absolute")
	}
	manifest := hostManifest{
		Name: NativeHostName, Description: "Site Wrapper settings helper", Path: executable,
		Type: "stdio", AllowedExtensions: []string{ExtensionID},
	}
	data, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}
	return os.WriteFile(path, append(data, '\n'), 0600)
}
