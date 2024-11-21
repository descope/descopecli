package theme

import (
	"context"

	"github.com/descope/descopecli/shared"
)

func Export(args []string) error {
	theme, err := shared.Descope.Management.Flow().ExportTheme(context.Background())
	if err != nil {
		return err
	}

	if isEmpty := shared.OmitEmpty(theme["references"]); isEmpty {
		delete(theme, "references")
	}

	if len(args) != 0 {
		return shared.WriteJSONFile(args[0], theme)
	}

	if shared.Flags.Json {
		shared.ExitWithResult(theme, "theme", "Exported theme")
	}

	shared.PrintIndented(theme)
	return nil
}

func Import(args []string) error {
	theme, err := shared.ReadJSONFile[map[string]any](args[0])
	if err != nil {
		return err
	}
	return shared.Descope.Management.Flow().ImportTheme(context.Background(), theme)
}
