package project

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/descope/descopecli/shared"
	"github.com/descope/go-sdk/descope"
)

func Clone(args []string) (err error) {
	res, err := shared.Descope.Management.Project().Clone(context.Background(), args[1], descope.ProjectTag(Flags.Tag))
	if err == nil {
		b, _ := json.Marshal(res)
		fmt.Printf("* Cloned project: %s\n", string(b))
	}
	return err
}
