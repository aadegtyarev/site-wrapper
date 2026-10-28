package bridge

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/aadegtyarev/firefox-site-wrapper/internal/settings"
)

func TestNativeMessageRoundTrip(t *testing.T) {
	request := Request{Action: "unlock", Password: "test-password"}
	data, err := json.Marshal(request)
	if err != nil {
		t.Fatal(err)
	}
	var frame bytes.Buffer
	if err := binary.Write(&frame, binary.LittleEndian, uint32(len(data))); err != nil {
		t.Fatal(err)
	}
	frame.Write(data)
	got, err := ReadRequest(&frame)
	if err != nil || got.Action != request.Action || got.Password != request.Password {
		t.Fatalf("read = %+v, %v", got, err)
	}
	var output bytes.Buffer
	if err := WriteResponse(&output, Response{OK: true, AppVersion: "1.2.3"}); err != nil {
		t.Fatal(err)
	}
	var size uint32
	if err := binary.Read(&output, binary.LittleEndian, &size); err != nil {
		t.Fatal(err)
	}
	responseData, err := io.ReadAll(&output)
	if err != nil || int(size) != len(responseData) {
		t.Fatalf("response size = %d, bytes = %d, err = %v", size, len(responseData), err)
	}
}

func TestReadRequestRejectsOversizedFrame(t *testing.T) {
	var frame bytes.Buffer
	binary.Write(&frame, binary.LittleEndian, uint32(maxMessageBytes+1))
	if _, err := ReadRequest(&frame); err == nil {
		t.Fatal("oversized frame accepted")
	}
}

func TestServiceStartupAndAuthorization(t *testing.T) {
	dir := t.TempDir()
	service := Service{SettingsPath: filepath.Join(dir, "settings.json"), LaunchPath: filepath.Join(dir, "launch-url"), AppVersion: "1.2.3"}
	if err := SetLaunchURL(service.LaunchPath, "https://first.example"); err != nil {
		t.Fatal(err)
	}
	res, err := service.Handle(Request{Action: "startup"})
	if err != nil || !res.NeedsSetup || !res.OK || res.LaunchURL != "https://first.example" {
		t.Fatalf("first startup = %+v, %v", res, err)
	}
	draft := settings.Draft{URL: "https://saved.example", ZoomPercent: 125, Mode: settings.Kiosk, Language: settings.German}
	res, err = service.Handle(Request{Action: "save", NewPassword: "long-test-password", Settings: draft})
	if err != nil || !res.OK || res.Settings.URL != draft.URL {
		t.Fatalf("save = %+v, %v", res, err)
	}
	res, err = service.Handle(Request{Action: "unlock", Password: "wrong"})
	if err != nil || res.ErrorCode != "incorrect_password" || res.Settings != nil {
		t.Fatalf("wrong password = %+v, %v", res, err)
	}
	if err := SetLaunchURL(service.LaunchPath, "https://override.example"); err != nil {
		t.Fatal(err)
	}
	res, err = service.Handle(Request{Action: "startup"})
	if err != nil || res.LaunchURL != "https://override.example" || res.Settings.Language != settings.German {
		t.Fatalf("startup = %+v, %v", res, err)
	}
	if _, err := os.Stat(service.LaunchPath); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("launch override was not consumed: %v", err)
	}
}

func TestSaveRequestUsesWireFieldNames(t *testing.T) {
	var req Request
	data := []byte(`{"action":"save","new_password":"long-test-password","settings":{"url":"https://example.org","zoom_percent":150,"mode":"window","language":"it"}}`)
	if err := json.Unmarshal(data, &req); err != nil {
		t.Fatal(err)
	}
	if req.NewPassword != "long-test-password" || req.Settings.ZoomPercent != 150 || req.Settings.Language != settings.Italian {
		t.Fatalf("wire fields decoded incorrectly: %+v", req)
	}
}
