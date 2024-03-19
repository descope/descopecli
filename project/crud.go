package project

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/descope/descopecli/shared"
	"github.com/descope/go-sdk/descope"
)

func Clone(args []string) error {
	res, err := shared.Descope.Management.Project().Clone(context.Background(), args[1], descope.ProjectTag(Flags.Tag))
	if err == nil {
		b, _ := json.Marshal(res)
		fmt.Printf("* Cloned project: %s\n", string(b))
	}
	return err
}

func Delete(args []string) error {
	if !Flags.Force {
		fmt.Printf("Are you sure you want to delete project %s (this cannot be undone): [y/N] ", args[0])
		reader := bufio.NewReader(os.Stdin)
		response, _ := reader.ReadString('\n')
		response = strings.ToLower(strings.TrimSpace(response))
		if response != "y" && response != "yes" {
			return nil
		}
	}
	err := shared.Descope.Management.Project().Delete(context.Background())
	if err == nil {
		fmt.Printf("* Deleted project: %s\n", args[0])
	}
	return err
}
