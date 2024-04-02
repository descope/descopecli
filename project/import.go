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
		root = "project-" + args[0]
	} else {
		root = filepath.Clean(root)
	}
	if info, err := os.Stat(root); os.IsNotExist(err) || !info.IsDir() {
		return errors.New("snapshot path does not exist: " + root)
	}

	im := importer{root: root}
	return im.Run(validate)
}

type importer struct {
	root  string
	files map[string]any
}

func (im *importer) Run(validate bool) (err error) {
	shared.PrintProgress("Reading snapshot files:")

	im.files = map[string]any{}
	if err := im.readFiles(im.root); err != nil {
		return err
	}

	if Flags.Debug {
		WriteDebugFile(im.root, "debug/import.log", im.files)
	}

	var secrets *descope.SnapshotSecrets
	if Flags.SecretsInput != "" {
		secrets, err = im.readSecrets(Flags.SecretsInput)
		if err != nil {
			return err
		}
	}

	if validate {
		req := &descope.ValidateSnapshotRequest{Files: im.files, InputSecrets: secrets}
		return im.Validate(req)
	}

	req := &descope.ImportSnapshotRequest{Files: im.files, InputSecrets: secrets}
	return im.Import(req)
}

func (im *importer) Import(req *descope.ImportSnapshotRequest) error {
	shared.PrintProgress("Importing snapshot...")
	if err := shared.Descope.Management.Project().ImportSnapshot(context.Background(), req); err != nil {
		return err
	}
	shared.PrintProgress("Done")
	return nil
}

func (im *importer) Validate(req *descope.ValidateSnapshotRequest) error {
	shared.PrintProgress("Validating snapshot...")
	res, err := shared.Descope.Management.Project().ValidateSnapshot(context.Background(), req)
	if err != nil {
		return err
	}
	if res.Ok {
		shared.PrintProgress("Done")
		return nil
	}

	if shared.Flags.Json {
		shared.PrintIndented(res)
	} else {
		shared.PrintProgress("Validation failed:")
		for _, f := range res.Failures {
			shared.PrintItem(f)
		}
	}

	if Flags.FailuresOutput != "" {
		if err := im.writeFailures(Flags.FailuresOutput, res.Failures); err != nil {
			return err
		}
	}

	if Flags.SecretsOutput != "" && (len(res.MissingSecrets.Connectors) > 0 || len(res.MissingSecrets.OAuthProviders) > 0) {
		if err := im.writeSecrets(Flags.SecretsOutput, res.MissingSecrets); err != nil {
			return err
		}
	}

	shared.PrintProgress("Done")

	// differentiate between validation failures with status code 2, as opposed to 1 for all other errors
	os.Exit(2)
	return nil
}

func (im *importer) readFiles(path string) error {
	info, err := os.ReadDir(path)
	if err != nil {
		return fmt.Errorf("failed to read snapshot files from path %s: %w", path, err)
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
		return fmt.Errorf("failed to parse snapshot file path %s: %w", fullpath, err)
	}

	bytes, err := os.ReadFile(fullpath)
	if err != nil {
		return fmt.Errorf("failed to read snapshot file %s: %w", relpath, err)
	}

	skipped := false

	switch filepath.Ext(fullpath) {
	case ".json":
		var m map[string]any
		if err = json.Unmarshal(bytes, &m); err != nil {
			return fmt.Errorf("failed to convert snapshot json file %s: %w", relpath, err)
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
		shared.PrintItem(relpath)
	}
	return nil
}

const (
	connectorPrefix = "connector-"
	oauthPrefix     = "oauthprovider-"
)

type secretEntry struct {
	Name    string            `json:"name"`
	Secrets map[string]string `json:"secrets"`
}

func (im *importer) readSecrets(path string) (*descope.SnapshotSecrets, error) {
	var file map[string]*secretEntry

	shared.PrintProgress("Reading input secrets...")
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read secrets input file %s: %w", path, err)
	}
	if err = json.Unmarshal(bytes, &file); err != nil {
		return nil, fmt.Errorf("failed to convert secrets input json %s: %w", path, err)
	}

	secrets := &descope.SnapshotSecrets{}
	for k, entry := range file {
		if strings.HasPrefix(k, connectorPrefix) {
			id := strings.TrimPrefix(k, connectorPrefix)
			for typ, value := range entry.Secrets {
				secrets.Connectors = append(secrets.Connectors, &descope.SnapshotSecret{ID: id, Name: entry.Name, Type: typ, Value: value})
			}
		} else if strings.HasPrefix(k, oauthPrefix) {
			id := strings.TrimPrefix(k, oauthPrefix)
			for typ, value := range entry.Secrets {
				secrets.OAuthProviders = append(secrets.OAuthProviders, &descope.SnapshotSecret{ID: id, Name: entry.Name, Type: typ, Value: value})
			}
		} else {
			return nil, fmt.Errorf("unexpected secret type: %s", k)
		}
	}

	found := len(secrets.Connectors) + len(secrets.OAuthProviders)
	if found == 0 {
		shared.PrintProgress("No secrets found in input")
	} else if found == 1 {
		shared.PrintProgress("Found 1 secret in input")
	} else {
		shared.PrintProgress(fmt.Sprintf("Found %d secrets in input", found))
	}

	return secrets, nil
}

func (im *importer) writeFailures(path string, failures []string) error {
	shared.PrintProgress("Writing failures output...")
	var b []byte
	for _, f := range failures {
		b = append(b, f...)
		b = append(b, '\n')
	}
	if err := os.WriteFile(path, b, 0644); err != nil {
		return fmt.Errorf("failed to write failures output file %s: %w", path, err)
	}
	return nil
}

func (im *importer) writeSecrets(path string, secrets *descope.SnapshotSecrets) error {
	shared.PrintProgress("Writing secrets output...")
	file := map[string]*secretEntry{}
	for _, v := range secrets.Connectors {
		key := connectorPrefix + v.ID
		if _, ok := file[key]; !ok {
			file[key] = &secretEntry{Name: v.Name, Secrets: map[string]string{}}
		}
		file[key].Secrets[v.Type] = ""
	}
	for _, v := range secrets.OAuthProviders {
		key := oauthPrefix + v.ID
		if _, ok := file[key]; !ok {
			file[key] = &secretEntry{Name: v.Name, Secrets: map[string]string{}}
		}
		file[key].Secrets[v.Type] = ""
	}

	b, _ := json.MarshalIndent(file, "", "  ")
	b = append(b, '\n')
	if err := os.WriteFile(path, b, 0644); err != nil {
		return fmt.Errorf("failed to write secrets output file %s: %w", path, err)
	}

	return nil
}
