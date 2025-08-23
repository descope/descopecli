package flow

import (
	"context"
	"encoding/json"

	"github.com/descope/descopecli/shared"
	"github.com/descope/go-sdk/descope"
)

func RunManagementFlow(args []string) error {
	options := &descope.MgmtFlowOptions{}
	if len(args) > 1 && args[1] != "" {
		if err := json.Unmarshal([]byte(args[1]), &options); err != nil {
			return err
		}
	}
	output, err := shared.Descope.Management.Flow().RunManagementFlow(context.Background(), args[0], options)
	if err != nil {
		return err
	}

	shared.ExitWithMap(output, "Management flow execution result")
	return nil
}

func RunManagementFlowAsync(args []string) error {
	options := &descope.MgmtFlowOptions{}
	if len(args) > 1 && args[1] != "" {
		if err := json.Unmarshal([]byte(args[1]), &options); err != nil {
			return err
		}
	}
	executionRef, err := shared.Descope.Management.Flow().RunManagementFlowAsync(context.Background(), args[0], options)
	if err != nil {
		return err
	}

	shared.ExitWithResult(executionRef, "executionRef", "Management flow async execution ref")
	return nil
}

func GetManagementFlowAsyncResult(args []string) error {
	output, err := shared.Descope.Management.Flow().GetManagementFlowAsyncResult(context.Background(), args[0])
	if err != nil {
		return err
	}

	shared.ExitWithMap(output, "Management flow async execution result")
	return nil
}

func List(_ []string) error {
	res, err := shared.Descope.Management.Flow().ListFlows(context.Background())
	if err != nil {
		return err
	}

	shared.ExitWithSlice(res.Flows, "flows", "Flow", "Loaded", "flow", "flows")
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
