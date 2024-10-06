package user

import (
	"context"
	"errors"
	"net/url"

	"github.com/descope/descopecli/shared"
	"github.com/descope/go-sdk/descope"
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

func GenerarteTestUserOTP(args []string) error {
	method := descope.DeliveryMethod(args[0])
	if method != "email" && method != "sms" && method != "voice" {
		return errors.New("method must be either email, sms, or voice")
	}
	code, err := shared.Descope.Management.User().GenerateOTPForTestUser(context.Background(), method, args[1], nil)
	if err != nil {
		return err
	}
	shared.ExitWithResult(code, "code", "Generated OTP for test user")
	return nil
}

func GenerarteTestUserMagicLink(args []string) error {
	method := descope.DeliveryMethod(args[0])
	if method != "email" && method != "sms" {
		return errors.New("method must be either email or sms")
	}
	link, err := shared.Descope.Management.User().GenerateMagicLinkForTestUser(context.Background(), method, args[1], Flags.RedirectURL, nil)
	if err != nil {
		return err
	}
	if _, err := url.ParseRequestURI(link); err != nil {
		return errors.New("ensure a redirect URL is configured or specify one with the -u flag")
	}
	shared.ExitWithResult(link, "link", "Generated magic link for test user")
	return nil
}

func GenerarteTestUserEnchantedLink(args []string) error {
	link, pendingRef, err := shared.Descope.Management.User().GenerateEnchantedLinkForTestUser(context.Background(), args[0], Flags.RedirectURL, nil)
	if err != nil {
		return err
	}
	if _, err := url.ParseRequestURI(link); err != nil {
		return errors.New("ensure a redirect URL is configured or specify one with the -u flag")
	}
	result := enchantedLink{Link: link, PendingRef: pendingRef}
	shared.ExitWithResult(result, "result", "Generated enchanted link for test user")
	return nil
}
