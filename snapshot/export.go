package snapshot

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/descope/descopecli/shared"
	"github.com/descope/go-sdk/descope"
	"golang.org/x/exp/maps"
)

func Export(args []string) error {
	mkdir := false

	root := Flags.Path
	if root == "" {
		root = "project-" + args[0]
		mkdir = true
	} else {
		root = filepath.Clean(root)
		if info, err := os.Stat(root); os.IsNotExist(err) {
			mkdir = true
		} else if err != nil || !info.IsDir() {
			return errors.New("invalid export path: " + root)
		}
	}

	if mkdir {
		if err := os.Mkdir(root, 0755); err != nil && !os.IsExist(err) {
			return errors.New("cannot create export path: " + root)
		}
	}

	ex := exporter{root: root}
	err := ex.Export()
	if mkdir && err != nil {
		os.Remove(root)
	}

	return err
}

type exporter struct {
	root string
}

func (ex *exporter) Export() error {
	shared.PrintProgress("Exporting snapshot...")

	req := &descope.ExportSnapshotRequest{}
	if Flags.NoAssets {
		req.Format = "plain"
	}

	res, err := shared.Descope.Management.Project().ExportSnapshot(context.Background(), req)
	if err != nil {
		return err
	}

	if Flags.Debug {
		WriteDebugFile(ex.root, "debug/export.log", res.Files)
	}

	shared.PrintProgress("Writing snapshot files:")

	paths := maps.Keys(res.Files)
	sort.Strings(paths)

	for _, path := range paths {
		shared.PrintItem(path)
		data := res.Files[path]
		if object, ok := data.(map[string]any); ok {
			if err := ex.writeObject(path, object); err != nil {
				return err
			}
		} else if asset, ok := data.(string); ok {
			if err := ex.writeAsset(path, asset); err != nil {
				return err
			}
		} else {
			return errors.New("unexpected exported file data: " + path)
		}
	}

	shared.PrintProgress("Done")
	return nil
}

func (ex *exporter) writeObject(path string, object map[string]any) error {
	bytes, err := json.MarshalIndent(object, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to format object file %s: %w", path, err)
	}
	bytes = append(bytes, '\n')
	return ex.writeBytes(path, bytes)
}

func (ex *exporter) writeAsset(path string, asset string) error {
	if filepath.Ext(path) == ".txt" || filepath.Ext(path) == ".html" {
		bytes := []byte(asset)
		if n := len(bytes); n != 0 && bytes[n-1] != '\n' {
			bytes = append(bytes, '\n')
		}
		return ex.writeBytes(path, bytes)
	}
	bytes, err := base64.StdEncoding.DecodeString(asset)
	if err != nil {
		return fmt.Errorf("failed to decode asset file %s: %w", path, err)
	}
	return ex.writeBytes(path, bytes)
}

func (ex *exporter) writeBytes(path string, bytes []byte) error {
	fullpath, err := ex.ensurePath(path)
	if err != nil {
		return err
	}
	if err = os.WriteFile(fullpath, bytes, 0644); err != nil {
		return fmt.Errorf("failed to write asset file %s: %w", path, err)
	}
	return nil
}

func (ex *exporter) ensurePath(path string) (string, error) {
	dir, file := filepath.Split(path)
	fullpath := ex.root
	if dir != "" {
		for _, d := range strings.Split(filepath.Clean(dir), string(filepath.Separator)) {
			fullpath = filepath.Join(fullpath, d)
			if err := os.Mkdir(fullpath, 0755); err != nil && !os.IsExist(err) {
				return "", fmt.Errorf("failed to create export subdirectory %s: %w", fullpath, err)
			}
		}
	}

	fullpath = filepath.Join(fullpath, file)
	return fullpath, nil
}

func WriteDebugFile(root, path string, object map[string]any) {
	ex := exporter{root: root}
	ex.writeObject(path, object)
}
