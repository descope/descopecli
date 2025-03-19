package flow

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/descope/descopecli/shared"
	"golang.org/x/exp/maps"
)

func Convert(args []string) error {
	sourcePath := filepath.Clean(args[0])
	sourceInfo, err := os.Stat(sourcePath)
	if os.IsNotExist(err) {
		return errors.New("source path does not exist: " + sourcePath)
	}

	if sourceInfo.IsDir() {
		if len(args) < 2 {
			return errors.New("target file path is required when converting from snapshot format")
		}

		targetPath := filepath.Clean(args[1])
		if targetInfo, err := os.Stat(targetPath); err == nil && targetInfo.IsDir() {
			targetPath = filepath.Join(targetPath, filepath.Base(sourcePath), ".json")
		}

		m, err := convertFlowSnapshotToExported(sourcePath)
		if err != nil {
			return err
		}

		return writeFlowFile(sourcePath, targetPath, m)
	}

	m, err := shared.ReadJSONFile[map[string]any](sourcePath)
	if err != nil {
		return abortUnsupported(sourcePath)
	}

	if m["flow"] != nil && m["screens"] != nil {
		targetPath := sourcePath
		if len(args) > 1 {
			targetPath = filepath.Clean(args[1])
			if targetInfo, err := os.Stat(targetPath); err == nil && targetInfo.IsDir() {
				targetPath = filepath.Join(targetPath, filepath.Base(sourcePath))
			}
		}

		m, err := convertFlowConsoleToExported(sourcePath)
		if err != nil {
			return err
		}

		return writeFlowFile(sourcePath, targetPath, m)
	}

	if m["contents"] != nil && m["metadata"] != nil {
		if len(args) < 2 {
			return errors.New("target directory path is required when converting to snapshot format")
		}

		targetPath := filepath.Clean(args[1])
		if _, err := os.Stat(targetPath); os.IsNotExist(err) {
			if err := os.Mkdir(targetPath, 0750); err != nil {
				return fmt.Errorf("failed to create target directory %s: %w", targetPath, err)
			}
		}

		m, err := convertFlowExportedToSnapshot(sourcePath)
		if err != nil {
			return err
		}

		shared.PrintProgress(fmt.Sprintf("Writing files: %s", targetPath))

		for filename, data := range m {
			path := filepath.Join(targetPath, filename)
			if err := shared.WriteJSONFile(path, data); err != nil {
				return err
			}
		}

		return nil
	}

	return abortUnsupported(sourcePath)
}

func abortUnsupported(path string) error {
	if Flags.Skip {
		shared.PrintProgress(fmt.Sprintf("Skipping unsupported input: %s", path))
		return nil
	}
	return fmt.Errorf("unsupported input format: %s", path)
}

func convertFlowSnapshotToExported(sourcePath string) (map[string]any, error) {
	flowID := filepath.Base(sourcePath)

	metadata, err := shared.ReadJSONFile[map[string]any](filepath.Join(sourcePath, "metadata.json"))
	if err != nil {
		return nil, abortUnsupported(sourcePath)
	}

	// extract the list of screen names (default to empty)
	screens, _ := metadata["screens"].([]any)
	delete(metadata, "screens")

	// extract the references
	references, _ := metadata["references"].(map[string]any)
	delete(metadata, "references")

	contents, err := shared.ReadJSONFile[map[string]any](filepath.Join(sourcePath, "contents.json"))
	if err != nil {
		return nil, abortUnsupported(sourcePath)
	}

	metadataStruct, err := shared.ReadJSONFile[exportedMetadata](filepath.Join(sourcePath, "metadata.json"))
	if err != nil {
		return nil, abortUnsupported(sourcePath)
	}

	flow := &exportedFlow{
		FlowId:     flowID,
		Metadata:   &metadataStruct,
		Contents:   contents,
		References: references,
	}

	for _, v := range screens {
		screenID, ok := v.(string)
		if !ok {
			return nil, fmt.Errorf("invalid value in list of screens in path: %s", sourcePath)
		}

		// look for the screen with the specified names from the list
		screen, err := shared.ReadJSONFile[exportedScreen](filepath.Join(sourcePath, "screen-"+screenID+".json"))
		if err != nil {
			return nil, fmt.Errorf("missing flow screen in path: %s", sourcePath)
		}
		flow.Screens = append(flow.Screens, &screen)
	}

	slices.SortFunc(flow.Screens, func(a, b *exportedScreen) int {
		return strings.Compare(a.ScreenId, b.ScreenId)
	})

	return convertJSONObject(flow), nil
}

func convertFlowExportedToSnapshot(sourcePath string) (map[string]map[string]any, error) {
	source, err := shared.ReadJSONFile[exportedFlow](sourcePath)
	if err != nil {
		return nil, err
	}

	metadata := convertJSONObject(source.Metadata)
	contents := source.Contents
	references := source.References

	// delete additional redundant fields we don't want in the exported data
	delete(metadata, "widget")
	if metadata["translation"] == nil {
		delete(metadata, "translation")
	}
	if len(source.Metadata.SharedInteractions) == 0 {
		delete(metadata, "sharedInteractions")
	}

	// collect the converted screen objects and convert them to maps
	screens := map[string]map[string]any{}
	for _, screen := range source.Screens {
		object := convertJSONObject(screen)
		if len(screen.Interactions) == 0 {
			delete(object, "interactions")
		}
		screens[screen.ScreenId] = object
	}

	// add references to metadata (if there's anything there)
	if len(source.References) != 0 {
		metadata["references"] = references
	}

	// we add a list of screens to the metadata so we know which files to look for during import
	screenIDs := maps.Keys(screens)
	slices.Sort(screenIDs)
	metadata["screens"] = screenIDs

	m := map[string]map[string]any{
		"metadata.json": metadata,
		"contents.json": contents,
	}
	for screenID, object := range screens {
		m["screen-"+screenID+".json"] = object
	}
	return m, nil
}

func convertFlowConsoleToExported(sourcePath string) (map[string]any, error) {
	source, err := shared.ReadJSONFile[consoleWrapper](sourcePath)
	if err != nil {
		return nil, err
	}

	componentsVersion := ""
	exportedScreens := []*exportedScreen{}
	for _, screen := range source.Screens {
		s := &exportedScreen{
			ScreenId:     screen.ID,
			Interactions: screen.Interactions,
			Contents:     screen.HTMLTemplate,
		}
		if screen.ComponentsVersion != "" {
			componentsVersion = screen.ComponentsVersion
		}
		exportedScreens = append(exportedScreens, s)
	}

	var translation *exportedTranslation
	if source.Flow.Translate || source.Flow.TranslateConnectorID != "" || source.Flow.TranslateSourceLang != "" || len(source.Flow.TranslateTargetLangs) > 0 {
		translation = &exportedTranslation{
			Enabled:         source.Flow.Translate,
			ConnectorId:     source.Flow.TranslateConnectorID,
			SourceLanguage:  source.Flow.TranslateSourceLang,
			TargetLanguages: source.Flow.TranslateTargetLangs,
		}
	}

	result := &exportedFlow{
		FlowId: source.Flow.ID,
		Metadata: &exportedMetadata{
			Name:               source.Flow.Name,
			Description:        source.Flow.Description,
			ComponentsVersion:  componentsVersion,
			Disabled:           source.Flow.Disabled,
			Fingerprint:        source.Flow.Fingerprint,
			Widget:             source.Flow.Widget,
			Translation:        translation,
			SharedInteractions: source.Flow.SharedInteractions,
		},
		Contents:   source.Flow.DSL,
		Screens:    exportedScreens,
		References: map[string]any{},
	}

	return convertJSONObject(result), nil
}

func writeFlowFile(sourcePath, targetPath string, object map[string]any) error {
	if object == nil {
		if Flags.Skip {
			shared.PrintProgress(fmt.Sprintf("Skipping unsupported file: %s", sourcePath))
			return nil
		}
		return fmt.Errorf("unsupported file format: %s", sourcePath)
	}

	if sourcePath == targetPath {
		shared.PrintProgress(fmt.Sprintf("Overwriting source file: %s", sourcePath))
	} else {
		shared.PrintProgress(fmt.Sprintf("Writing file: %s", targetPath))
	}

	return shared.WriteJSONFile(targetPath, object)
}

func convertJSONObject(v any) map[string]any {
	bytes, err := json.Marshal(v)
	if err != nil {
		panic("unexpected marshal failure")
	}

	var m map[string]any
	if err := json.Unmarshal(bytes, &m); err != nil {
		panic("unexpected unmarshal failure")
	}
	return m
}

type exportedFlow struct {
	FlowId     string            `json:"flowId,omitempty"`
	Metadata   *exportedMetadata `json:"metadata,omitempty"`
	Contents   map[string]any    `json:"contents,omitempty"`
	Screens    []*exportedScreen `json:"screens,omitempty"`
	References map[string]any    `json:"references,omitempty"`
}

type exportedMetadata struct {
	Name               string               `json:"name,omitempty"`
	Description        string               `json:"description,omitempty"`
	ComponentsVersion  string               `json:"componentsVersion,omitempty"`
	Disabled           bool                 `json:"disabled,omitempty"`
	Fingerprint        bool                 `json:"fingerprint,omitempty"`
	Widget             bool                 `json:"widget,omitempty"`
	Translation        *exportedTranslation `json:"translation,omitempty"`
	SharedInteractions []map[string]any     `json:"sharedInteractions,omitempty"`
}

type exportedTranslation struct {
	Enabled         bool     `json:"enabled,omitempty"`
	ConnectorId     string   `json:"connectorId,omitempty"`
	SourceLanguage  string   `json:"sourceLanguage,omitempty"`
	TargetLanguages []string `json:"targetLanguages,omitempty"`
}

type exportedScreen struct {
	ScreenId     string           `json:"screenId,omitempty"`
	Interactions []map[string]any `json:"interactions,omitempty"`
	Contents     map[string]any   `json:"contents,omitempty"`
}

type consoleWrapper struct {
	Flow    *consoleFlow     `json:"flow,omitempty"`
	Screens []*consoleScreen `json:"screens,omitempty"`
}

type consoleFlow struct {
	ID                   string           `json:"id,omitempty"`
	Version              int              `json:"version,omitempty"`
	Name                 string           `json:"name,omitempty"`
	Description          string           `json:"description,omitempty"`
	DSL                  map[string]any   `json:"dsl,omitempty"`
	ModifiedTime         int64            `json:"modifiedTime,omitempty"`
	ETag                 string           `json:"etag,omitempty"`
	Disabled             bool             `json:"disabled,omitempty"`
	Translate            bool             `json:"translate,omitempty"`
	TranslateConnectorID string           `json:"translateConnectorID,omitempty"`
	TranslateSourceLang  string           `json:"translateSourceLang,omitempty"`
	TranslateTargetLangs []string         `json:"translateTargetLangs,omitempty"`
	Fingerprint          bool             `json:"fingerprint,omitempty"`
	Widget               bool             `json:"widget,omitempty"`
	SharedInteractions   []map[string]any `json:"sharedInteractions,omitempty"`
}

type consoleScreen struct {
	ID                string           `json:"id,omitempty"`
	Version           int              `json:"version,omitempty"`
	FlowID            string           `json:"flowId,omitempty"`
	Inputs            []map[string]any `json:"inputs,omitempty"`
	Interactions      []map[string]any `json:"interactions,omitempty"`
	HTMLTemplate      map[string]any   `json:"htmlTemplate,omitempty"`
	ComponentsVersion string           `json:"componentsVersion,omitempty"`
}
