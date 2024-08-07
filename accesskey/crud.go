package accesskey

import (
	"context"
	"encoding/json"

	"github.com/descope/descopecli/shared"
	"github.com/descope/go-sdk/descope"
)

func Create(args []string) error {
	tenants := []*descope.AssociatedTenant{}
	for _, tenantID := range Flags.Tenants {
		tenants = append(tenants, &descope.AssociatedTenant{TenantID: tenantID})
	}

	cleartext, res, err := shared.Descope.Management.AccessKey().Create(context.Background(), args[0], Flags.Description, Flags.Expires, nil, tenants, Flags.UserId, nil, nil)
	if err != nil {
		return err
	}

	b, _ := json.Marshal(res)
	m := map[string]any{}
	json.Unmarshal(b, &m)
	m["cleartext"] = cleartext

	shared.ExitWithResult(m, "accessKey", "Created access key")
	return nil
}

func Delete(args []string) error {
	return shared.Descope.Management.AccessKey().Delete(context.Background(), args[0])
}

func Load(args []string) error {
	key, err := shared.Descope.Management.AccessKey().Load(context.Background(), args[0])
	if err != nil {
		return err
	}

	shared.ExitWithResult(key, "accessKey", "Loaded access key")
	return nil
}

func LoadAll(_ []string) error {
	res, err := shared.Descope.Management.AccessKey().SearchAll(context.Background(), nil)
	if err != nil {
		return err
	}

	shared.ExitWithResults(res, "accessKeys", "Access key", "Loaded", "access key", "access keys")
	return nil
}
