package main

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/awf-project/cli/pkg/plugin/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Compile-time interface assertions.
var (
	_ sdk.Plugin            = (*TimePlugin)(nil)
	_ sdk.OperationProvider = (*TimePlugin)(nil)
)

func newTimePlugin() *TimePlugin {
	return &TimePlugin{
		BasePlugin: sdk.BasePlugin{
			PluginName:    "awf-plugin-time",
			PluginVersion: "1.2.0",
		},
	}
}

func TestTimePlugin_Operations(t *testing.T) {
	p := newTimePlugin()

	ops := p.Operations()

	require.Len(t, ops, 1)
	assert.Equal(t, "time", ops[0])
}

func TestTimePlugin_Init(t *testing.T) {
	p := newTimePlugin()

	err := p.Init(context.Background(), map[string]any{})

	assert.NoError(t, err)
}

func TestTimePlugin_Shutdown(t *testing.T) {
	p := newTimePlugin()

	err := p.Shutdown(context.Background())

	assert.NoError(t, err)
}

func TestTimePlugin_HandleOperation(t *testing.T) {
	tests := []struct {
		name      string
		inputs    map[string]any
		wantErr   bool
		checkFunc func(t *testing.T, result *sdk.OperationResult)
	}{
		{
			name:   "default no inputs",
			inputs: map[string]any{},
			checkFunc: func(t *testing.T, result *sdk.OperationResult) {
				assert.True(t, result.Success)
				assert.NotEmpty(t, result.Output)
				_, err := time.Parse(time.RFC3339, result.Output)
				assert.NoError(t, err, "default output should be RFC3339")
				assertSuccessData(t, result)
				assert.Equal(t, time.Now().Location().String(), result.Data["timezone"])
			},
		},
		{
			name:   "custom format date only",
			inputs: map[string]any{"format": "2006-01-02"},
			checkFunc: func(t *testing.T, result *sdk.OperationResult) {
				assert.True(t, result.Success)
				assert.Len(t, result.Output, 10)
				_, err := time.Parse("2006-01-02", result.Output)
				assert.NoError(t, err, "output should match custom format")
				assertSuccessData(t, result)
			},
		},
		{
			name:   "UTC timezone",
			inputs: map[string]any{"tz": "UTC"},
			checkFunc: func(t *testing.T, result *sdk.OperationResult) {
				assert.True(t, result.Success)
				assert.Contains(t, result.Output, "Z", "UTC output should contain Z suffix")
				assertSuccessData(t, result)
				assert.Equal(t, "UTC", result.Data["timezone"])
			},
		},
		{
			name:   "specific timezone",
			inputs: map[string]any{"tz": "America/New_York"},
			checkFunc: func(t *testing.T, result *sdk.OperationResult) {
				assert.True(t, result.Success)
				assertSuccessData(t, result)
				assert.Equal(t, "America/New_York", result.Data["timezone"])
			},
		},
		{
			name:   "invalid timezone",
			inputs: map[string]any{"tz": "Invalid/Zone"},
			checkFunc: func(t *testing.T, result *sdk.OperationResult) {
				assert.False(t, result.Success)
				assert.Contains(t, result.Error, "invalid timezone")
				assert.Contains(t, result.Error, "Invalid/Zone")
			},
		},
		{
			name:   "empty format string falls back to RFC3339",
			inputs: map[string]any{"format": ""},
			checkFunc: func(t *testing.T, result *sdk.OperationResult) {
				assert.True(t, result.Success)
				_, err := time.Parse(time.RFC3339, result.Output)
				assert.NoError(t, err, "empty format should fall back to RFC3339")
				assertSuccessData(t, result)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := newTimePlugin()

			result, err := p.HandleOperation(context.Background(), "time", tt.inputs)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, result)
			tt.checkFunc(t, result)
		})
	}
}

func TestTimePlugin_HandleOperation_UnknownOperation(t *testing.T) {
	p := newTimePlugin()

	result, err := p.HandleOperation(context.Background(), "invalid", nil)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "unknown operation")
}

func TestTimePlugin_Context_Cancellation(t *testing.T) {
	p := newTimePlugin()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	result, err := p.HandleOperation(ctx, "time", map[string]any{})

	if err == nil {
		require.NotNil(t, result)
	}
}

func TestPluginYAMLManifest(t *testing.T) {
	content, err := os.ReadFile("plugin.yaml")
	require.NoError(t, err, "plugin.yaml must exist")

	manifest := string(content)
	assert.Contains(t, manifest, "name:")
	assert.Contains(t, manifest, "awf-plugin-time")
	assert.Contains(t, manifest, "version:")
	assert.Contains(t, manifest, "1.2.0")
	assert.Contains(t, manifest, "awf_version:")
	assert.Contains(t, manifest, "capabilities:")
	assert.Contains(t, manifest, "operations")
}

func assertSuccessData(t *testing.T, result *sdk.OperationResult) {
	t.Helper()
	require.NotNil(t, result.Data)

	ts, ok := result.Data["timestamp"].(string)
	require.True(t, ok, "timestamp should be a string")
	_, err := time.Parse(time.RFC3339, ts)
	assert.NoError(t, err, "timestamp should be valid RFC3339")

	// Note: int64 assertion works in unit tests (direct struct access).
	// Through gRPC serialization, numeric types may arrive as float64.
	unix, ok := result.Data["unix"]
	require.True(t, ok, "unix should be present")
	unixVal, ok := unix.(int64)
	require.True(t, ok, "unix should be int64")
	assert.Greater(t, unixVal, int64(0), "unix should be positive")

	tz, ok := result.Data["timezone"].(string)
	require.True(t, ok, "timezone should be a string")
	assert.NotEmpty(t, tz, "timezone should not be empty")
}
