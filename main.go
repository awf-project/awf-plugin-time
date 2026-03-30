package main

import (
	"context"
	"fmt"
	"time"

	"github.com/awf-project/cli/pkg/plugin/sdk"
)

// TimePlugin returns system date/time for workflow injection.
type TimePlugin struct {
	sdk.BasePlugin
}

func (p *TimePlugin) Operations() []string { return []string{"time"} }

func (p *TimePlugin) HandleOperation(_ context.Context, name string, inputs map[string]any) (*sdk.OperationResult, error) {
	if name != "time" {
		return nil, fmt.Errorf("unknown operation: %s", name)
	}

	format := sdk.GetStringDefault(inputs, "format", time.RFC3339)
	if format == "" {
		format = time.RFC3339
	}

	tzName := sdk.GetStringDefault(inputs, "tz", "Local")
	loc, err := time.LoadLocation(tzName)
	if err != nil {
		return sdk.NewErrorResultf("invalid timezone: %s", tzName), nil
	}

	now := time.Now().In(loc)

	return sdk.NewSuccessResult(now.Format(format), map[string]any{
		"timestamp": now.Format(time.RFC3339),
		"unix":      now.Unix(),
		"timezone":  loc.String(),
	}), nil
}

func main() {
	sdk.Serve(&TimePlugin{
		BasePlugin: sdk.BasePlugin{
			PluginName:    "awf-plugin-time",
			PluginVersion: "1.2.0",
		},
	})
}
