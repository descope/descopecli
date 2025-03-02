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

const errProjectNotFound = "E011008"

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
	// only works with company management key, ignore errors other than project not found
	projects, err := shared.Descope.Management.Project().ListProjects(context.Background())
	if err, ok := err.(*descope.Error); ok && err.Code == errProjectNotFound {
		return fmt.Errorf("cannot find project with id %s", args[0])
	}

	alias := args[0]
	name := ""
	for _, p := range projects {
		if p.ID == args[0] {
			alias = p.Name
			name = p.Name
		}
	}

	if !Flags.Force {
		if shared.Flags.Json {
			return errors.New("the --force flag is required when using --json to delete a project")
		}

		fmt.Printf("* Delete the project \"%s\"?\n", alias)
		fmt.Printf("* All project information will be deleted, including users, flows and analytics.\n")
		fmt.Printf("* This cannot be undone!\n")
		if name != "" {
			fmt.Printf("* Enter the project name to confirm: ")
		} else {
			fmt.Printf("* Enter the project id to confirm: ")
		}

		reader := bufio.NewReader(os.Stdin)
		response, _ := reader.ReadString('\n')
		if response = strings.TrimSpace(response); response != alias {
			fmt.Printf("* Cancelled\n")
			return nil
		}
	}

	if err := shared.Descope.Management.Project().Delete(context.Background()); err != nil {
		return err
	}

	shared.PrintProgress(fmt.Sprintf(`Deleted project "%s"`, alias))
	return nil
}

func List(_ []string) error {
	res, err := shared.Descope.Management.Project().ListProjects(context.Background())
	if err != nil {
		return err
	}

	shared.ExitWithSlice(res, "projects", "Project", "Loaded", "project", "projects")
	return nil
}
