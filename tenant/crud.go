package tenant

import (
	"context"

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
	if err != nil {
		return err
	}

	shared.ExitWithResult(tenantID, "tenantId", "Created new tenant with id")
	return nil
}

func Delete(args []string) error {
	return shared.Descope.Management.Tenant().Delete(context.Background(), args[0], false)
}

func Load(args []string) error {
	tenant, err := shared.Descope.Management.Tenant().Load(context.Background(), args[0])
	if err != nil {
		return err
	}

	shared.ExitWithResult(tenant, "tenant", "Loaded tenant")
	return nil
}

func LoadAll(_ []string) error {
	res, err := shared.Descope.Management.Tenant().LoadAll(context.Background())
	if err != nil {
		return err
	}

	shared.ExitWithResults(res, "tenants", "Tenant", "Loaded", "tenant", "tenants")
	return nil
}
