package user

import (
	"context"

	"github.com/descope/descopecli/shared"
)

func SetActivePassword(args []string) error {
	loginID := args[0]
	password := args[1]
	return shared.Descope.Management.User().SetActivePassword(context.Background(), loginID, password)
}

func SetTemporaryPassword(args []string) error {
	loginID := args[0]
	password := args[1]
	return shared.Descope.Management.User().SetTemporaryPassword(context.Background(), loginID, password)
}

func ExpirePassword(args []string) error {
	loginID := args[0]
	return shared.Descope.Management.User().ExpirePassword(context.Background(), loginID)
}
