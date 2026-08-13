package shared

import (
	"strings"
	"testing"

	"github.com/descope/go-sdk/descope"
)

const (
	testManagementKey   = "K2abcdef1234567890"
	testManagementKeyID = "K2EHcqqqqqqqqqqqqqqqqqqqqqqq"
	testToken           = "eyJhbGciOiJSUzI1NiJ9.eyJpc3MiOiJodHRwczovL3Rva2VuLmFjdGlvbnMuZ2l0aHVidXNlcmNvbnRlbnQuY29tIn0.sig"
)

func TestIsValidManagementCredential(t *testing.T) {
	for _, tc := range []struct {
		name  string
		key   string
		valid bool
	}{
		{"management key", testManagementKey, true},
		// A workload identity token is only a credential once it names the key it acts as.
		{"oidc token paired with a key id", testToken + ":" + testManagementKeyID, true},
		{"bare oidc token", testToken, false},
		{"oidc token with an empty key id", testToken + ":", false},
		{"oidc token paired with a non-key", testToken + ":not-a-key", false},
		{"empty", "", false},
		{"garbage", "not-a-credential", false},
		{"partial jwt", "eyJhbGciOiJSUzI1NiJ9.eyJpc3MiOiJ4In0:" + testManagementKeyID, false},
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
// ID, so `<jwt>:<keyId>` is what produces the `Bearer <projectId>:<jwt>:<keyId>` the Management API expects.
func TestManagementCredentialPairsTokenWithKeyID(t *testing.T) {
	t.Setenv(descope.EnvironmentVariableManagementKey, "")
	t.Setenv(EnvironmentVariableWorkloadIdentityToken, testToken)
	t.Setenv(EnvironmentVariableManagementKeyID, testManagementKeyID)

	credential, envVar, err := managementCredential()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if want := testToken + ":" + testManagementKeyID; credential != want {
		t.Errorf("credential = %q, want %q", credential, want)
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
		keyID       string
		errContains string
	}{
		{"neither set", "", "", "", descope.EnvironmentVariableManagementKey + " or " + EnvironmentVariableWorkloadIdentityToken},
		{"management key only", testManagementKey, "", "", ""},
		{"token with its key id", "", testToken, testManagementKeyID, ""},
		{"token takes precedence over key", testManagementKey, testToken, testManagementKeyID, ""},
		// A token alone cannot be authorized: the API needs to know which key it is acting as.
		{"token without a key id", "", testToken, "", EnvironmentVariableManagementKeyID},
		// Pasting the whole key into the key-ID variable puts a secret somewhere unmasked.
		{"whole management key in the key id slot", "", testToken, strings.Repeat("K", 80), EnvironmentVariableManagementKeyID},
		// the error has to name the variable the bad value actually came from
		{"bad token with a valid key set", testManagementKey, "garbage", testManagementKeyID, EnvironmentVariableWorkloadIdentityToken},
		{"bad management key", "garbage", "", "", descope.EnvironmentVariableManagementKey},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv(descope.EnvironmentVariableManagementKey, tc.mgmtKey)
			t.Setenv(EnvironmentVariableWorkloadIdentityToken, tc.token)
			t.Setenv(EnvironmentVariableManagementKeyID, tc.keyID)
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
