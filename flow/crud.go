package flow

import (
	"context"

	"github.com/descope/descopecli/shared"
)

func List(_ []string) error {
	res, err := shared.Descope.Management.Flow().ListFlows(context.Background())
	if err != nil {
		return err
	}

	shared.ExitWithResults(res.Flows, "flows", "Flow", "Loaded", "flow", "flows")
	return nil
}

func Export(args []string) error {
	flow, err := shared.Descope.Management.Flow().ExportFlow(context.Background(), args[0])
	if err != nil {
		return err
	}

	shared.OmitEmpty(flow["metadata"])
	if isEmpty := shared.OmitEmpty(flow["references"]); isEmpty {
		delete(flow, "references")
	}
	if screens, ok := flow["screens"].([]any); ok {
		for i := range screens {
			if screen, ok := screens[i].(map[string]any); ok {
				shared.OmitEmpty(screen["interactions"])
			}
		}
	}

	if len(args) == 2 {
		return shared.WriteJSONFile(args[1], flow)
	}

	if shared.Flags.Json {
		shared.ExitWithResult(flow, "flow", "Exported flow")
	}

	shared.PrintIndented(flow)
	return nil
}

func Import(args []string) error {
	flow, err := shared.ReadJSONFile[map[string]any](args[1])
	if err != nil {
		return err
	}
	return shared.Descope.Management.Flow().ImportFlow(context.Background(), args[0], flow)
}
