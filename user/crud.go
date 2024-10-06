package user

import (
	"context"
	"errors"

	"github.com/descope/descopecli/shared"
	"github.com/descope/go-sdk/descope"
)

func Create(args []string) error {
	return createUser(args, false)
}

func createUser(args []string, test bool) error {
	tenants := []*descope.AssociatedTenant{}
	for _, tenantID := range Flags.Tenants {
		tenants = append(tenants, &descope.AssociatedTenant{TenantID: tenantID})
	}

	u := &descope.UserRequest{}
	u.Email = Flags.Email
	u.Phone = Flags.Phone
	u.Name = Flags.Name
	u.Tenants = tenants

	var err error
	var user *descope.UserResponse
	if test {
		user, err = shared.Descope.Management.User().CreateTestUser(context.Background(), args[0], u)
	} else {
		user, err = shared.Descope.Management.User().Create(context.Background(), args[0], u)
	}
	if err != nil {
		return err
	}

	shared.ExitWithResult(user, "user", "Created user")
	return nil
}

func Delete(_ []string) error {
	if (Flags.LoginID == "" && Flags.UserID == "") || (Flags.LoginID != "" && Flags.UserID != "") {
		return errors.New("this command requires 1 flag to identify the user")
	}
	if Flags.LoginID != "" {
		return shared.Descope.Management.User().Delete(context.Background(), Flags.LoginID)
	}
	return shared.Descope.Management.User().DeleteByUserID(context.Background(), Flags.UserID)
}

func Load(_ []string) error {
	if (Flags.LoginID == "" && Flags.UserID == "") || (Flags.LoginID != "" && Flags.UserID != "") {
		return errors.New("this command requires 1 flag to identify the user")
	}

	var user *descope.UserResponse
	var err error
	if Flags.LoginID != "" {
		user, err = shared.Descope.Management.User().Load(context.Background(), Flags.LoginID)
	} else {
		user, err = shared.Descope.Management.User().LoadByUserID(context.Background(), Flags.UserID)
	}
	if err != nil {
		return err
	}

	shared.ExitWithResult(user, "user", "Loaded user")
	return nil
}

func LoadAll(_ []string) error {
	res, _, err := shared.Descope.Management.User().SearchAll(context.Background(), &descope.UserSearchOptions{Limit: int32(Flags.Limit), Page: int32(Flags.Page)})
	if err != nil {
		return err
	}

	shared.ExitWithResults(res, "users", "User", "Loaded", "user", "users")
	return nil
}
