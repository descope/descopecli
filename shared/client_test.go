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
		{"oidc token", testToken, true},
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

// The credential handed to the SDK matters as much as whether it was accepted: the SDK joins it onto the project
// ID, so the token has to go out alone to produce the `Bearer <projectId>:<jwt>` the Management API expects. The
// key being assumed travels in the token's `aud`, not in the credential.
func TestManagementCredentialSendsTokenAlone(t *testing.T) {
	t.Setenv(descope.EnvironmentVariableManagementKey, "")
	t.Setenv(EnvironmentVariableWorkloadIdentityToken, testToken)

	credential, envVar, err := managementCredential()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if credential != testToken {
		t.Errorf("credential = %q, want %q", credential, testToken)
	}
	if envVar != EnvironmentVariableWorkloadIdentityToken {
		t.Errorf("envVar = %q, want %q", envVar, EnvironmentVariableWorkloadIdentityToken)
	}
}

func TestCredentialEnvironmentVariables(t *testing.T) {
	for _, tc := range []struct {
		name        string
		mgmtKey     string
		token       string
		errContains string
	}{
		{"neither set", "", "", descope.EnvironmentVariableManagementKey + " or " + EnvironmentVariableWorkloadIdentityToken},
		{"management key only", testManagementKey, "", ""},
		{"token only", "", testToken, ""},
		// The credentials are exclusive on the server, so setting both is rejected here rather than silently
		// resolved in favour of one of them.
		{"management key and token together", testManagementKey, testToken,
			descope.EnvironmentVariableManagementKey + " and " + EnvironmentVariableWorkloadIdentityToken},
		// the error has to name the variable the bad value actually came from
		{"bad token", "", "garbage", EnvironmentVariableWorkloadIdentityToken},
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
