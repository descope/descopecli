package shared

import (
	"strings"
	"testing"

	"github.com/descope/go-sdk/descope"
)

const (
	testManagementKey = "K2abcdef1234567890"
	testToken         = "eyJhbGciOiJSUzI1NiJ9.eyJpc3MiOiJodHRwczovL3Rva2VuLmFjdGlvbnMuZ2l0aHVidXNlcmNvbnRlbnQuY29tIn0.sig"
)

func TestIsValidManagementCredential(t *testing.T) {
	for _, tc := range []struct {
		name  string
		key   string
		valid bool
	}{
		{"management key", testManagementKey, true},
		{"workload token", testToken, true},
		{"empty", "", false},
		{"garbage", "not-a-credential", false},
		{"partial jwt", "eyJhbGciOiJSUzI1NiJ9.eyJpc3MiOiJ4In0", false},
		{"project id", "P2abcdef1234567890", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := isValidManagementCredential(tc.key); got != tc.valid {
				t.Errorf("isValidManagementCredential(%q) = %v, want %v", tc.key, got, tc.valid)
			}
		})
	}
}

// Both credentials arrive in the same variable, so the client has to accept either one there.
func TestManagementKeyEnvironmentVariable(t *testing.T) {
	for _, tc := range []struct {
		name        string
		credential  string
		errContains string
	}{
		{"not set", "", "must be set"},
		{"management key", testManagementKey, ""},
		{"workload token", testToken, ""},
		{"garbage", "garbage", "must be a valid management key or workload token"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv(descope.EnvironmentVariableManagementKey, tc.credential)
			t.Setenv(descope.EnvironmentVariableProjectID, "P2abcdef1234567890")
			// the credential checks all run before the client is built, so keep the client itself
			// from reaching the network on the cases that get that far
			t.Setenv(descope.EnvironmentVariableBaseURL, "http://127.0.0.1:1")

			_, err := createDescopeClient(nil, false, false)
			switch {
			case tc.errContains == "" && err != nil:
				t.Errorf("unexpected error: %v", err)
			case tc.errContains != "" && err == nil:
				t.Errorf("expected an error mentioning %q, got none", tc.errContains)
			case tc.errContains != "" && !strings.Contains(err.Error(), tc.errContains):
				t.Errorf("error %q does not mention %q", err, tc.errContains)
			}
		})
	}
}
