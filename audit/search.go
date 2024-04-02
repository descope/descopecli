package audit

import (
	"context"

	"github.com/descope/descopecli/shared"
	"github.com/descope/go-sdk/descope"
)

func Search(args []string) error {
	res, err := shared.Descope.Management.Audit().Search(context.Background(), &descope.AuditSearchOptions{Text: args[0]})
	if err != nil {
		return err
	}

	shared.ExitWithResults(res, "records", "Record", "Found", "record", "records")
	return nil
}
