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
	credential, envVar := managementCredential()
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

func managementCredential() (credential string, envVar string) {
	if token := os.Getenv(EnvironmentVariableWorkloadIdentityToken); token != "" {
		return token, EnvironmentVariableWorkloadIdentityToken
	}
	return os.Getenv(descope.EnvironmentVariableManagementKey), descope.EnvironmentVariableManagementKey
}

// The management credential is either a static management key, which starts with "K", or a
// workload identity token, such as the OIDC token a CI job requests for itself,
// which is a compact JWT with three dot-separated parts. This is only a shape check to catch
// an obviously wrong environment variable early: the server decides which credential it got
// and whether it is trusted.
func isValidManagementCredential(key string) bool {
	return strings.HasPrefix(key, "K") || strings.Count(key, ".") == 2
}
