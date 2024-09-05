package project

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/descope/descopecli/shared"
	"github.com/descope/go-sdk/descope"
)

func Clone(args []string) error {
	env := descope.ProjectEnvironment(Flags.Environment)
	if env != descope.ProjectEnvironmentNone && env != descope.ProjectEnvironmentProduction {
		return errors.New(`the only valid value for the optional --environment flag is "production"`)
	}

	res, err := shared.Descope.Management.Project().Clone(context.Background(), args[1], env, Flags.Tags)
	if err != nil {
		return err
	}

	shared.ExitWithResult(res, "result", "Cloned project")
	return nil
}

func Delete(args []string) error {
	if !Flags.Force {
		if shared.Flags.Json {
			return errors.New("the --force flag is required when using --json to delete a project")
		}
		shared.PrintProgress(fmt.Sprintf("Are you sure you want to delete project %s (this cannot be undone): [y/N] ", args[0]))
		reader := bufio.NewReader(os.Stdin)
		response, _ := reader.ReadString('\n')
		response = strings.ToLower(strings.TrimSpace(response))
		if response != "y" && response != "yes" {
			return nil
		}
	}

	if err := shared.Descope.Management.Project().Delete(context.Background()); err != nil {
		return err
	}

	shared.PrintProgress(fmt.Sprintf("Deleted project %s", args[0]))
	return nil
}

func List(_ []string) error {
	res, err := shared.Descope.Management.Project().ListProjects(context.Background())
	if err != nil {
		return err
	}

	shared.ExitWithResults(res, "projects", "Project", "Loaded", "project", "projects")
	return nil
}
