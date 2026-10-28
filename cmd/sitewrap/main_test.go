package main

import (
	"reflect"
	"testing"
)

func TestFirefoxArgs(t *testing.T) {
	for _, tc := range []struct {
		kiosk bool
		want  []string
	}{
		{false, []string{"--new-instance", "--profile", "/profile", "about:blank"}},
		{true, []string{"--new-instance", "--profile", "/profile", "--kiosk", "about:blank"}},
	} {
		got := firefoxArgs("/profile", tc.kiosk)
		if !reflect.DeepEqual(got, tc.want) {
			t.Fatalf("firefoxArgs = %v, want %v", got, tc.want)
		}
	}
}

func TestNativeHostInvocation(t *testing.T) {
	if !isNativeHostInvocation([]string{"/home/user/.mozilla/native-messaging-hosts/io.github.aadegtyarev.sitewrap.json", "sitewrap@aadegtyarev.github.io"}) {
		t.Fatal("Firefox native host arguments not recognized")
	}
	if isNativeHostInvocation([]string{"https://example.org", "sitewrap@aadegtyarev.github.io"}) {
		t.Fatal("URL treated as native host invocation")
	}
}
