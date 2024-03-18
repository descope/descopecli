package project

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/descope/descopecli/shared"
	"github.com/descope/go-sdk/descope"
)

func Import(args []string) error {
	return runImporter(args, false)
}

func Validate(args []string) error {
	return runImporter(args, true)
}

func runImporter(args []string, validate bool) error {
	root := Flags.Path
	if root == "" {
		root = "env-" + args[0]
	} else {
		root = filepath.Clean(root)
	}
	if info, err := os.Stat(root); os.IsNotExist(err) || !info.IsDir() {
		return errors.New("import path does not exist: " + root)
	}

	im := importer{root: root}
	return im.Run(validate)
}

type importer struct {
	root  string
	files map[string]any
}

func (im *importer) Run(validate bool) error {
	fmt.Println("* Reading files...")

	im.files = map[string]any{}
	if err := im.readFiles(im.root); err != nil {
		return err
	}

	if Flags.Debug {
		WriteDebugFile(im.root, "debug/import.log", im.files)
	}

	var secrets *descope.ImportProjectSecrets
	secrets, err := im.readSecrets(Flags.SecretsInput)
	if err != nil {
		return err
	}

	req := &descope.ImportProjectRequest{Files: im.files, InputSecrets: secrets}
	if validate {
		return im.Validate(req)
	}
	return im.Import(req)
}

func (im *importer) Import(req *descope.ImportProjectRequest) error {
	fmt.Println("* Importing project...")
	if err := shared.Descope.Management.Project().Import(context.Background(), req); err != nil {
		return fmt.Errorf("failed to import project: %w", err)
	}
	fmt.Println("* Done")
	return nil
}

func (im *importer) Validate(req *descope.ImportProjectRequest) error {
	fmt.Println("* Validating import data...")
	res, err := shared.Descope.Management.Project().ValidateImport(context.Background(), req)
	if err != nil {
		return fmt.Errorf("failed to validate project: %w", err)
	}
	if res.Ok {
		return nil
	}

	fmt.Println("* Validation failed:")
	for _, f := range res.Failures {
		fmt.Printf("  - %s\n", f)
	}
	if Flags.FailureOutput != "" {
		var b []byte
		for _, f := range res.Failures {
			b = append(b, f...)
			b = append(b, '\n')
		}
		if err := os.WriteFile(Flags.FailureOutput, b, 0644); err != nil {
			return fmt.Errorf("failed to write failure output file %s: %w", Flags.FailureOutput, err)
		}
	}

	if len(res.MissingSecrets.Connectors) > 0 || len(res.MissingSecrets.OAuthProviders) > 0 {
		if Flags.SecretsOutput == "" {
			return errors.New("validation failed with missing secrets but no secrets output path was specified")
		}
		im.writeSecrets(Flags.SecretsOutput, res.MissingSecrets)
	}

	// differentiate validation failures with status code 2, as opposed to 1 for all other errors
	os.Exit(2)
	return nil
}

func (im *importer) readFiles(path string) error {
	info, err := os.ReadDir(path)
	if err != nil {
		return fmt.Errorf("failed to read import files from path %s: %w", path, err)
	}

	for _, entry := range info {
		basename := entry.Name()
		if im.shouldIgnorePath(basename) {
			continue
		}
		fullpath := filepath.Join(path, basename)
		if entry.IsDir() {
			if err := im.readFiles(fullpath); err != nil {
				return err
			}
		} else {
			if err := im.readFile(fullpath); err != nil {
				return err
			}
		}
	}

	return nil
}

func (im *importer) shouldIgnorePath(basename string) bool {
	return strings.HasPrefix(basename, "__MACOSX") || strings.HasPrefix(basename, ".")
}

func (im *importer) readFile(fullpath string) error {
	relpath, err := filepath.Rel(im.root, fullpath)
	if err != nil {
		return fmt.Errorf("failed to parse import file path %s: %w", fullpath, err)
	}

	bytes, err := os.ReadFile(fullpath)
	if err != nil {
		return fmt.Errorf("failed to read import file %s: %w", relpath, err)
	}

	skipped := false

	switch filepath.Ext(fullpath) {
	case ".json":
		var m map[string]any
		if err = json.Unmarshal(bytes, &m); err != nil {
			return fmt.Errorf("failed to convert import json file %s: %w", relpath, err)
		}
		im.files[relpath] = m
	case ".txt", ".html":
		im.files[relpath] = string(bytes)
	case ".png", ".jpg", ".svg", ".webp":
		im.files[relpath] = base64.StdEncoding.EncodeToString(bytes)
	default:
		skipped = true
	}

	if !skipped {
		fmt.Printf("  - %s\n", relpath)
	}
	return nil
}

type secretEntries struct {
	ID      string            `json:"id"`
	Secrets map[string]string `json:"secrets"`
}

type secretsFile struct {
	Connectors     map[string]*secretEntries `json:"connectors,omitempty"`
	OAuthProviders map[string]*secretEntries `json:"oauthproviders,omitempty"`
}

func (im *importer) readSecrets(path string) (*descope.ImportProjectSecrets, error) {
	var file secretsFile

	fmt.Println("* Reading input secrets...")
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read secrets input file %s: %w", path, err)
	}
	if err = json.Unmarshal(bytes, &file); err != nil {
		return nil, fmt.Errorf("failed to convert secrets input json %s: %w", path, err)
	}

	secrets := &descope.ImportProjectSecrets{}
	if file.Connectors != nil {
		for k, v := range file.Connectors {
			for typ, value := range v.Secrets {
				secrets.Connectors = append(secrets.Connectors, &descope.ImportProjectSecret{ID: v.ID, Name: k, Type: typ, Value: value})
			}
		}
	}
	if file.OAuthProviders != nil {
		for k, v := range file.Connectors {
			for typ, value := range v.Secrets {
				secrets.Connectors = append(secrets.Connectors, &descope.ImportProjectSecret{ID: v.ID, Name: k, Type: typ, Value: value})
			}
		}
	}

	found := len(secrets.Connectors) + len(secrets.OAuthProviders)
	if found == 0 {
		fmt.Println("* No secrets found in input")
	} else if found == 1 {
		fmt.Println("* Found 1 secret in input")
	} else {
		fmt.Printf("* Found %d secrets in input\n", found)
	}

	return secrets, nil
}

func (im *importer) writeSecrets(path string, secrets *descope.ImportProjectSecrets) error {
	fmt.Println("* Writing missing secrets output...")
	file := secretsFile{}
	for _, v := range secrets.Connectors {
		if file.Connectors == nil {
			file.Connectors = map[string]*secretEntries{}
		}
		if _, ok := file.Connectors[v.Name]; !ok {
			file.Connectors[v.Name] = &secretEntries{ID: v.ID}
		}
		file.Connectors[v.Name].Secrets[v.Type] = ""
	}
	for _, v := range secrets.OAuthProviders {
		if file.OAuthProviders == nil {
			file.OAuthProviders = map[string]*secretEntries{}
		}
		if _, ok := file.OAuthProviders[v.Name]; !ok {
			file.OAuthProviders[v.Name] = &secretEntries{ID: v.ID}
		}
		file.OAuthProviders[v.Name].Secrets[v.Type] = ""
	}

	b, _ := json.MarshalIndent(file, "", "  ")
	b = append(b, '\n')
	if err := os.WriteFile(path, b, 0644); err != nil {
		return fmt.Errorf("failed to write secrets output file %s: %w", path, err)
	}

	return nil
}
