package accesskey

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/descope/descopecli/shared"
	"github.com/descope/go-sdk/descope"
)

func Create(args []string) error {
	tenants := []*descope.AssociatedTenant{}
	for _, tenantID := range Flags.Tenants {
		tenants = append(tenants, &descope.AssociatedTenant{TenantID: tenantID})
	}

	cleartext, res, err := shared.Descope.Management.AccessKey().Create(context.Background(), args[0], Flags.Expires, nil, tenants, Flags.UserId)
	if err != nil {
		return err
	}

	b, _ := json.Marshal(res)
	m := map[string]any{}
	json.Unmarshal(b, &m)
	m["cleartext"] = cleartext
	b, _ = json.Marshal(m)
	fmt.Printf("* Created access key: %s\n", string(b))

	return nil
}

func Delete(args []string) error {
	return shared.Descope.Management.AccessKey().Delete(context.Background(), args[0])
}

func Load(args []string) error {
	res, err := shared.Descope.Management.AccessKey().Load(context.Background(), args[0])
	if err != nil {
		return err
	}

	b, _ := json.Marshal(res)
	fmt.Printf("* Loaded access key: %s\n", string(b))

	return nil
}

func LoadAll(_ []string) error {
	res, err := shared.Descope.Management.AccessKey().SearchAll(context.Background(), nil)
	if err != nil {
		return err
	}

	if len(res) == 1 {
		fmt.Printf("* Found 1 access key\n")
	} else {
		fmt.Printf("* Found %d access keys\n", len(res))
	}
	for i, key := range res {
		b, _ := json.Marshal(key)
		fmt.Printf("  - Access key %d: %s\n", int64(i), string(b))
	}

	return nil
}
