package audit

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/descope/descopecli/shared"
	"github.com/descope/go-sdk/descope"
)

func Search(args []string) error {
	res, err := shared.Descope.Management.Audit().Search(context.Background(), &descope.AuditSearchOptions{Text: args[0]})
	if err != nil {
		return err
	}

	if len(res) == 1 {
		fmt.Printf("* Found 1 record\n")
	} else {
		fmt.Printf("* Found %d records\n", len(res))
	}
	for i, record := range res {
		b, _ := json.MarshalIndent(record, "    ", "  ")
		fmt.Printf("  - Record %d:\n    %s\n", i, string(b))
	}
	return nil
}
