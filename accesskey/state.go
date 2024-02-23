package accesskey

import (
	"context"

	"github.com/descope/descopecli/shared"
)

func Deactivate(args []string) error {
	return shared.Descope.Management.AccessKey().Deactivate(context.Background(), args[0])
}

func Activate(args []string) error {
	return shared.Descope.Management.AccessKey().Activate(context.Background(), args[0])
}
