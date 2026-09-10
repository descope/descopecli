package shared

import (
	"errors"
	"os"
	"strings"

	"github.com/descope/go-sdk/descope"
	"github.com/descope/go-sdk/descope/client"
	"github.com/spf13/cobra"
)

// EnvironmentVariableWorkloadToken mirrors descope.EnvironmentVariableWorkloadToken from
// go-sdk#849, name and value. Replace this with the SDK's own constant once that is released.
const EnvironmentVariableWorkloadToken = "DESCOPE_WORKLOAD_TOKEN" // gitleaks:allow

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
		// a workload token goes in this same field: the SDK sends either one as the credential in the
		// authorization header, so from its point of view they are interchangeable and mutually exclusive
		ManagementKey: credential,
		// doesn't need to be specified in regular use
		DescopeBaseURL: os.Getenv(descope.EnvironmentVariableBaseURL),
	}

	if config.ManagementKey == "" {
		return nil, errors.New("the " + descope.EnvironmentVariableManagementKey + " or " + EnvironmentVariableWorkloadToken + " environment variable must be set")
	}
	if !isValidManagementCredential(config.ManagementKey) {
		return nil, errors.New("the " + envVar + " environment variable must be a valid management key or workload token")
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

func managementCredential() (credential string, envVar string, err error) {
	token := os.Getenv(EnvironmentVariableWorkloadToken)
	managementKey := os.Getenv(descope.EnvironmentVariableManagementKey)

	if token != "" && managementKey != "" {
		return "", EnvironmentVariableWorkloadToken, errors.New("the " + descope.EnvironmentVariableManagementKey +
			" and " + EnvironmentVariableWorkloadToken + " environment variables are exclusive, set only one:" +
			" a management key that a workload token acts as does not accept its own secret")
	}
	if token != "" {
		return token, EnvironmentVariableWorkloadToken, nil
	}
	return managementKey, descope.EnvironmentVariableManagementKey, nil
}

func isValidManagementCredential(key string) bool {
	if strings.HasPrefix(key, "K") {
		// a static management key
		return true
	}
	// a workload token from a trusted issuer (JWT)
	return strings.Count(key, ".") == 2
}
