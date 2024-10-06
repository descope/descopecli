package user

import (
	"context"

	"github.com/descope/descopecli/shared"
)

type enchantedLink struct {
	Link       string `json:"link,omitempty"`
	PendingRef string `json:"pendingRef,omitempty"`
}

func CreateTestUser(args []string) error {
	return createUser(args, true)
}

func DeleteAllTestUsers(_ []string) error {
	return shared.Descope.Management.User().DeleteAllTestUsers(context.Background())
}
