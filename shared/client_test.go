package shared

import "testing"

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
