package tenant

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/descope/descopecli/shared"
	"github.com/descope/go-sdk/descope"
)

func Create(args []string) (err error) {
	tr := &descope.TenantRequest{Name: args[0], SelfProvisioningDomains: Flags.Domains}

	var tenantID string
	if Flags.TenantID != "" {
		tenantID, err = Flags.TenantID, shared.Descope.Management.Tenant().CreateWithID(context.Background(), Flags.TenantID, tr)
	} else {
		tenantID, err = shared.Descope.Management.Tenant().Create(context.Background(), tr)
	}

	if err == nil {
		fmt.Printf("* Created new tenant with id: %s\n", tenantID)
	}
	return err
}

func Delete(args []string) error {
	return shared.Descope.Management.Tenant().Delete(context.Background(), args[0])
}

func Load(args []string) error {
	tenant, err := shared.Descope.Management.Tenant().Load(context.Background(), args[0])
	if err == nil {
		b, _ := json.Marshal(tenant)
		fmt.Printf("* Loaded tenant: %s\n", string(b))
	}
	return err
}

func LoadAll(args []string) error {
	res, err := shared.Descope.Management.Tenant().LoadAll(context.Background())
	if err == nil {
		if len(res) == 1 {
			fmt.Printf("* Found 1 tenant\n")
		} else {
			fmt.Printf("* Found %d tenants\n", len(res))
		}
		for i, tenant := range res {
			b, _ := json.Marshal(tenant)
			fmt.Printf("  - Tenant %d: %s\n", i, string(b))
		}
	}
	return err
}
