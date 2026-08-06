package shared

import (
	"strings"
	"testing"

	"github.com/descope/go-sdk/descope"
)

func TestIsValidManagementCredential(t *testing.T) {
	for _, tc := range []struct {
		name  string
		key   string
		valid bool
	}{
		{"management key", "K2abcdef1234567890", true},
		{"github actions oidc token", "eyJhbGciOiJSUzI1NiJ9.eyJpc3MiOiJodHRwczovL3Rva2VuLmFjdGlvbnMuZ2l0aHVidXNlcmNvbnRlbnQuY29tIn0.sig", true},
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

func TestCredentialEnvironmentVariables(t *testing.T) {
	const managementKey = "K2abcdef1234567890"
	const token = "eyJhbGciOiJSUzI1NiJ9.eyJpc3MiOiJodHRwczovL3Rva2VuLmFjdGlvbnMuZ2l0aHVidXNlcmNvbnRlbnQuY29tIn0.sig"

	for _, tc := range []struct {
		name        string
		mgmtKey     string
		token       string
		errContains string
	}{
		{"neither set", "", "", descope.EnvironmentVariableManagementKey + " or " + EnvironmentVariableWorkloadIdentityToken},
		{"management key only", managementKey, "", ""},
		{"workload identity token only", "", token, ""},
		{"token takes precedence over key", managementKey, token, ""},
		// the error has to name the variable the bad value actually came from
		{"bad token with a valid key set", managementKey, "garbage", EnvironmentVariableWorkloadIdentityToken},
		{"bad management key", "garbage", "", descope.EnvironmentVariableManagementKey},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv(descope.EnvironmentVariableManagementKey, tc.mgmtKey)
			t.Setenv(EnvironmentVariableWorkloadIdentityToken, tc.token)
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
