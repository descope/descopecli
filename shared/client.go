package shared

import (
	"errors"
	"os"
	"strings"

	"github.com/descope/go-sdk/descope"
	"github.com/descope/go-sdk/descope/client"
	"github.com/spf13/cobra"
)

const EnvironmentVariableWorkloadIdentityToken = "DESCOPE_WORKLOAD_IDENTITY_TOKEN" // gitleaks:allow

// EnvironmentVariableManagementKeyID names the management key a workload identity token acts as. It is an
// identifier, not a secret, so it belongs in a CI variable rather than a CI secret.
const EnvironmentVariableManagementKeyID = "DESCOPE_MANAGEMENT_KEY_ID"

// A management key ID is far shorter than a whole management key. Used only to catch the easy mistake of pasting
// the key itself into the key-ID variable, which would put a secret somewhere unmasked.
const maxManagementKeyIDLength = 70

var Descope *client.DescopeClient

func DefaultPreRun(cmd *cobra.Command, args []string) (err error) {
	cmd.SilenceUsage = true
	Descope, err = createDescopeClient(args, false, false)
	return err
}

func ProjectPreRun(cmd *cobra.Command, args []string) (err error) {
	cmd.SilenceUsage = true
	Descope, err = createDescopeClient(args, false, true)
	return err
}

func CompanyPreRun(cmd *cobra.Command, args []string) (err error) {
	cmd.SilenceUsage = true
	Descope, err = createDescopeClient(args, true, false)
	return err
}

func StandalonePreRun(cmd *cobra.Command, _ []string) error {
	cmd.SilenceUsage = true
	return nil
}

func createDescopeClient(args []string, company bool, project bool) (*client.DescopeClient, error) {
	credential, envVar, err := managementCredential()
	if err != nil {
		return nil, err
	}
	config := &client.Config{
		// optional as an environment variable in some commands
		ProjectID: os.Getenv(descope.EnvironmentVariableProjectID),
		// generate a management key in the Company section of the admin console: https://app.descope.com/settings/company
		ManagementKey: credential,
		// doesn't need to be specified in regular use
		DescopeBaseURL: os.Getenv(descope.EnvironmentVariableBaseURL),
	}

	if config.ManagementKey == "" {
		return nil, errors.New("the " + descope.EnvironmentVariableManagementKey + " or " + EnvironmentVariableWorkloadIdentityToken + " environment variable must be set")
	}
	if !isValidManagementCredential(config.ManagementKey) {
		return nil, errors.New("the " + envVar + " environment variable must be a valid management key or workload identity token")
	}

	if company {
		config.ProjectID = ""
		config.AllowEmptyProjectID = true
	} else if project {
		config.ProjectID = args[0]
		if !strings.HasPrefix(config.ProjectID, "P") {
			return nil, errors.New("the command argument must be a valid projectId")
		}
	} else {
		if config.ProjectID == "" {
			return nil, errors.New("the " + descope.EnvironmentVariableProjectID + " environment variable must be set")
		}
		if !strings.HasPrefix(config.ProjectID, "P") {
			return nil, errors.New("the " + descope.EnvironmentVariableProjectID + " environment variable must be a valid projectId")
		}
	}

	return client.NewWithConfig(config)
}

// managementCredential resolves the credential the Management API is called with.
//
// A static management key is used as-is. A workload identity token is not a credential on its own: the API
// authorizes it as one specific management key, the way sts:AssumeRoleWithWebIdentity requires a RoleArn. So the
// token is paired with that key's ID, and the SDK joins the pair onto the project ID to produce
// `Bearer <projectId>:<jwt>:<keyId>`.
func managementCredential() (credential string, envVar string, err error) {
	if token := os.Getenv(EnvironmentVariableWorkloadIdentityToken); token != "" {
		keyID := os.Getenv(EnvironmentVariableManagementKeyID)
		if keyID == "" {
			return "", EnvironmentVariableWorkloadIdentityToken, errors.New("the " + EnvironmentVariableManagementKeyID +
				" environment variable must be set alongside " + EnvironmentVariableWorkloadIdentityToken +
				": a workload identity token acts as a specific management key, so it has to name one")
		}
		if len(keyID) >= maxManagementKeyIDLength {
			return "", EnvironmentVariableManagementKeyID, errors.New("the " + EnvironmentVariableManagementKeyID +
				" environment variable looks like a whole management key rather than a key ID")
		}
		return token + ":" + keyID, EnvironmentVariableWorkloadIdentityToken, nil
	}
	return os.Getenv(descope.EnvironmentVariableManagementKey), descope.EnvironmentVariableManagementKey, nil
}

// The management credential is either a static management key, which starts with "K", or a
// workload identity token, such as the OIDC token a CI job requests for itself,
// which is a compact JWT with three dot-separated parts. This is only a shape check to catch
// an obviously wrong environment variable early: the server decides which credential it got
// and whether it is trusted.
func isValidManagementCredential(key string) bool {
	if strings.HasPrefix(key, "K") {
		return true
	}
	// A workload identity credential is the token and the key it acts as, joined.
	token, keyID, paired := strings.Cut(key, ":")
	return paired && strings.Count(token, ".") == 2 && strings.HasPrefix(keyID, "K")
}
