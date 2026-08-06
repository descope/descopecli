package shared

import (
	"errors"
	"os"
	"strings"

	"github.com/descope/go-sdk/descope"
	"github.com/descope/go-sdk/descope/client"
	"github.com/spf13/cobra"
)

// The Descope SDK doesn't define this one, since it takes the credential as a config field rather
// than reading it from the environment itself.
const EnvironmentVariableWorkloadIdentityToken = "WORKLOAD_IDENTITY_TOKEN" // gitleaks:allow

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
	config := &client.Config{
		// optional as an environment variable in some commands
		ProjectID: os.Getenv(descope.EnvironmentVariableProjectID),
		// generate a management key in the Company section of the admin console: https://app.descope.com/settings/company
		ManagementKey: os.Getenv(descope.EnvironmentVariableManagementKey),
		// doesn't need to be specified in regular use
		DescopeBaseURL: os.Getenv(descope.EnvironmentVariableBaseURL),
	}

	// a workload identity token is a short lived alternative to a static management key, so a CI job
	// can authenticate with the OIDC token it requests for itself instead of a stored secret, it is
	// sent to the server the same way and takes precedence when both variables are set
	credentialEnvVar := descope.EnvironmentVariableManagementKey
	if token := os.Getenv(EnvironmentVariableWorkloadIdentityToken); token != "" {
		config.ManagementKey, credentialEnvVar = token, EnvironmentVariableWorkloadIdentityToken
	}

	if config.ManagementKey == "" {
		return nil, errors.New("the " + descope.EnvironmentVariableManagementKey + " or " + EnvironmentVariableWorkloadIdentityToken + " environment variable must be set")
	}
	if !isValidManagementCredential(config.ManagementKey) {
		return nil, errors.New("the " + credentialEnvVar + " environment variable must be a valid management key or workload identity token")
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

// The management credential is either a static management key, which starts with "K", or a
// workload identity token, such as the OIDC token a GitHub Actions job requests for itself,
// which is a compact JWS with three dot-separated parts. This is only a shape check to catch
// an obviously wrong environment variable early: the server decides which credential it got
// and whether it is trusted.
func isValidManagementCredential(key string) bool {
	return strings.HasPrefix(key, "K") || strings.Count(key, ".") == 2
}
