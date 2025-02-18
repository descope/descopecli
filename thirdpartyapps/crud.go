package thirdpartyapps

import (
	"context"
	"errors"
	"strings"

	"github.com/descope/descopecli/shared"
	"github.com/descope/go-sdk/descope"
)

func Create(args []string) error {
	req, err := createReq("", args[0], args[1])
	if err != nil {
		return err
	}

	appID, secret, err := shared.Descope.Management.ThirdPartyApplication().CreateApplication(context.Background(), req)
	if err != nil {
		return err
	}

	app, err := shared.Descope.Management.ThirdPartyApplication().LoadApplication(context.Background(), appID)
	if err != nil {
		return err
	}

	m := map[string]any{}
	m["app"] = app
	m["secret"] = secret

	shared.ExitWithMap(m, "Created new third party app")
	return nil
}

func Update(args []string) error {
	req, err := createReq(args[0], args[1], args[2])
	if err != nil {
		return err
	}

	if err := shared.Descope.Management.ThirdPartyApplication().UpdateApplication(context.Background(), req); err != nil {
		return err
	}

	app, err := shared.Descope.Management.ThirdPartyApplication().LoadApplication(context.Background(), args[0])
	if err != nil {
		return err
	}

	shared.ExitWithResult(app, "app", "Updated third party app")
	return nil
}

func Delete(args []string) error {
	return shared.Descope.Management.ThirdPartyApplication().DeleteApplication(context.Background(), args[0])
}

func Load(args []string) error {
	app, err := shared.Descope.Management.ThirdPartyApplication().LoadApplication(context.Background(), args[0])
	if err != nil {
		return err
	}

	shared.ExitWithResult(app, "app", "Loaded third party app")
	return nil
}

func LoadAll(_ []string) error {
	apps, err := shared.Descope.Management.ThirdPartyApplication().LoadAllApplications(context.Background())
	if err != nil {
		return err
	}

	shared.ExitWithSlice(apps, "apps", "Third Party App", "Loaded", "third party app", "third party apps")
	return nil
}

func LoadSecret(args []string) error {
	secret, err := shared.Descope.Management.ThirdPartyApplication().GetApplicationSecret(context.Background(), args[0])
	if err != nil {
		return err
	}

	shared.ExitWithResult(secret, "secret", "Loaded secret for third party app")
	return nil
}

func RotateSecret(args []string) error {
	secret, err := shared.Descope.Management.ThirdPartyApplication().RotateApplicationSecret(context.Background(), args[0])
	if err != nil {
		return err
	}

	shared.ExitWithResult(secret, "secret", "Rotated secret for third party app")
	return nil
}

func createReq(id, name, url string) (req *descope.ThirdPartyApplicationRequest, err error) {
	req = &descope.ThirdPartyApplicationRequest{
		ID:                   id,
		Name:                 name,
		Description:          Flags.Description,
		LoginPageURL:         url,
		ApprovedCallbackUrls: Flags.CallbackURLs,
	}

	if len(Flags.PermissionScopes) == 0 {
		return nil, errors.New("third party apps require at least one permission scope")
	}

	req.PermissionsScopes, err = parseScopes(Flags.PermissionScopes, "permission", "roles")
	if err != nil {
		return nil, err
	}

	req.AttributesScopes, err = parseScopes(Flags.AttributeScopes, "attribute", "user attributes")
	if err != nil {
		return nil, err
	}

	return req, nil
}

func parseScopes(strs []string, kind string, values string) ([]*descope.ThirdPartyApplicationScope, error) {
	if len(strs)%3 != 0 {
		return nil, errors.New(kind + " scopes must be a comma separated string with a name, a description, and a list of " + values + " separated by pipe characters (or a single hyphen character for no " + values + ")")
	}
	scopes := []*descope.ThirdPartyApplicationScope{}
	for i := 0; i < len(strs); i += 3 {
		scope := &descope.ThirdPartyApplicationScope{Name: strs[i], Description: strs[i+1]}
		if scope.Name == "" {
			return nil, errors.New("the " + kind + " scope name must not be empty")
		}
		if scope.Description == "" {
			return nil, errors.New("the " + kind + " scope description must not be empty")
		}
		if v := strs[i+2]; v != "-" && v != "" {
			scope.Values = strings.Split(v, "|")
			for _, v := range scope.Values {
				if v == "" {
					return nil, errors.New("the " + scope.Name + " " + kind + " scope must not have empty " + values)
				}
			}
		}
		scopes = append(scopes, scope)
	}
	return scopes, nil
}
