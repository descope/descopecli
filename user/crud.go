package user

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/descope/descopecli/shared"
	"github.com/descope/go-sdk/descope"
)

func Create(args []string) error {
	tenants := []*descope.AssociatedTenant{}
	for _, tenantID := range Flags.Tenants {
		tenants = append(tenants, &descope.AssociatedTenant{TenantID: tenantID})
	}

	user := &descope.UserRequest{}
	user.Email = Flags.Email
	user.Phone = Flags.Phone
	user.Name = Flags.Name
	user.Tenants = tenants

	res, err := shared.Descope.Management.User().Create(context.Background(), args[0], user)
	if err != nil {
		return err
	}

	b, _ := json.Marshal(res)
	fmt.Printf("* Created user: %s\n", string(b))

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

	b, _ := json.Marshal(user)
	fmt.Printf("* Loaded user: %s\n", string(b))

	return nil
}

func LoadAll(_ []string) error {
	res, err := shared.Descope.Management.User().SearchAll(context.Background(), &descope.UserSearchOptions{Limit: int32(Flags.Limit), Page: int32(Flags.Page)})
	if err != nil {
		return err
	}

	if len(res) == 1 {
		fmt.Printf("* Found 1 user\n")
	} else {
		fmt.Printf("* Found %d users\n", len(res))
	}
	for i, user := range res {
		b, _ := json.Marshal(user)
		fmt.Printf("  - User %d: %s\n", i+Flags.Page*Flags.Limit, string(b))
	}

	return nil
}
