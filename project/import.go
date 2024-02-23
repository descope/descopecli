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
)

func Import(args []string) error {
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
	return im.Import()
}

type importer struct {
	root  string
	files map[string]any
}

func (im *importer) Import() error {
	fmt.Println("* Reading files...")

	im.files = map[string]any{}
	if err := im.readFiles(im.root); err != nil {
		return err
	}

	if Flags.Debug {
		WriteDebugFile(im.root, "debug/import.log", im.files)
	}

	fmt.Println("* Importing project...")
	if err := shared.Descope.Management.Project().Import(context.Background(), im.files); err != nil {
		return fmt.Errorf("failed to import project: %w", err)
	}

	fmt.Println("* Done")
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
