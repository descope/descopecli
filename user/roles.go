package user

import (
	"context"
	"errors"

	"github.com/descope/descopecli/shared"
)

func SetRoles(args []string) error {
	loginID := args[0]
	var err error
	if Flags.TenantID != "" {
		_, err = shared.Descope.Management.User().SetTenantRoles(context.Background(), loginID, Flags.TenantID, Flags.Roles)
	} else {
		_, err = shared.Descope.Management.User().SetRoles(context.Background(), loginID, Flags.Roles)
	}
	return err
}

func AddRoles(args []string) error {
	loginID := args[0]
	if len(Flags.Roles) == 0 {
		return errors.New("this command requires at least one role to be specified")
	}
	var err error
	if Flags.TenantID != "" {
		_, err = shared.Descope.Management.User().AddTenantRoles(context.Background(), loginID, Flags.TenantID, Flags.Roles)
	} else {
		_, err = shared.Descope.Management.User().AddRoles(context.Background(), loginID, Flags.Roles)
	}
	return err
}

func RemoveRoles(args []string) error {
	loginID := args[0]
	if len(Flags.Roles) == 0 {
		return errors.New("this command requires at least one role to be specified")
	}
	var err error
	if Flags.TenantID != "" {
		_, err = shared.Descope.Management.User().RemoveTenantRoles(context.Background(), loginID, Flags.TenantID, Flags.Roles)
	} else {
		_, err = shared.Descope.Management.User().RemoveRoles(context.Background(), loginID, Flags.Roles)
	}
	return err
}
