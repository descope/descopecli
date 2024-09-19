package user

import (
	"context"

	"github.com/descope/descopecli/shared"
)

func Activate(args []string) error {
	_, err := shared.Descope.Management.User().Activate(context.Background(), args[0])
	return err
}

func Deactivate(args []string) error {
	_, err := shared.Descope.Management.User().Deactivate(context.Background(), args[0])
	return err
}
