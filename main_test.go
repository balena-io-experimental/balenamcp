package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Example commands for each tool
var exampleCommands = map[string][]struct {
	name        string
	params      map[string]interface{}
	expectedCmd string
	expectedErr error
}{
	"device-list": {
		{
			name: "list all devices with JSON",
			params: map[string]interface{}{
				"json": true,
			},
			expectedCmd: "balena device list --json",
		},
		{
			name: "list devices in fleet",
			params: map[string]interface{}{
				"fleet": "my-fleet",
			},
			expectedCmd: "balena device list --fleet my-fleet",
		},
		{
			name: "list devices with help",
			params: map[string]interface{}{
				"help": true,
			},
			expectedCmd: "balena device list --help",
		},
	},
	"device-logs": {
		{
			name: "view device logs",
			params: map[string]interface{}{
				"device": "my-device",
			},
			expectedCmd: "balena device logs my-device",
		},
		{
			name: "view service logs",
			params: map[string]interface{}{
				"device":  "my-device",
				"service": "my-service",
			},
			expectedCmd: "balena device logs my-device --service my-service",
		},
		{
			name: "view system logs",
			params: map[string]interface{}{
				"device": "my-device",
				"system": true,
			},
			expectedCmd: "balena device logs my-device --system",
		},
		{
			name: "follow logs",
			params: map[string]interface{}{
				"device": "my-device",
				"tail":   true,
			},
			expectedCmd: "balena device logs my-device --tail",
		},
	},
	"fleet-list": {
		{
			name:        "list all fleets",
			params:      map[string]interface{}{},
			expectedCmd: "balena fleet list",
		},
		{
			name: "list fleets with JSON",
			params: map[string]interface{}{
				"json": true,
			},
			expectedCmd: "balena fleet list --json",
		},
		{
			name: "list fleets with help",
			params: map[string]interface{}{
				"help": true,
			},
			expectedCmd: "balena fleet list --help",
		},
	},
	"os-versions": {
		{
			name: "list OS versions",
			params: map[string]interface{}{
				"type": "raspberrypi4",
			},
			expectedCmd: "balena os versions raspberrypi4",
		},
		{
			name: "list ESR versions",
			params: map[string]interface{}{
				"type": "raspberrypi4",
				"esr":  true,
			},
			expectedCmd: "balena os versions raspberrypi4 --esr",
		},
		{
			name: "list draft versions",
			params: map[string]interface{}{
				"type":          "raspberrypi4",
				"include_draft": true,
			},
			expectedCmd: "balena os versions raspberrypi4 --include-draft",
		},
	},
	"release-list": {
		{
			name: "list all releases",
			params: map[string]interface{}{
				"fleet": "my-fleet",
			},
			expectedCmd: "balena release list my-fleet",
		},
		{
			name: "list releases with JSON",
			params: map[string]interface{}{
				"fleet": "my-fleet",
				"json":  true,
			},
			expectedCmd: "balena release list my-fleet --json",
		},
	},
	"release-info": {
		{
			name: "get release info",
			params: map[string]interface{}{
				"id": "123",
			},
			expectedCmd: "balena release 123",
		},
		{
			name: "get release info with JSON",
			params: map[string]interface{}{
				"id":   "123",
				"json": true,
			},
			expectedCmd: "balena release 123 --json",
		},
		{
			name: "get release composition",
			params: map[string]interface{}{
				"id":          "123",
				"composition": true,
			},
			expectedCmd: "balena release 123 --composition",
		},
	},
	"version": {
		{
			name:        "get version info",
			params:      map[string]interface{}{},
			expectedCmd: "balena version",
		},
		{
			name: "get version help",
			params: map[string]interface{}{
				"help": true,
			},
			expectedCmd: "balena version --help",
		},
	},
	"tag-list": {
		{
			name: "list device tags",
			params: map[string]interface{}{
				"resource": "device",
				"name":     "my-device",
			},
			expectedCmd: "balena tag list device my-device",
		},
		{
			name: "list fleet tags",
			params: map[string]interface{}{
				"resource": "fleet",
				"name":     "my-fleet",
			},
			expectedCmd: "balena tag list fleet my-fleet",
		},
		{
			name: "list release tags",
			params: map[string]interface{}{
				"resource": "release",
				"name":     "my-release",
			},
			expectedCmd: "balena tag list release my-release",
		},
	},
}

func runCommandTest(t *testing.T, tool string, testCase struct {
	name        string
	params      map[string]interface{}
	expectedCmd string
	expectedErr error
}) {
	// Start the MCP server as a subprocess in dry-run mode
	cmd := exec.Command("go", "run", "main.go", "-dry-run")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		t.Fatalf("Failed to get stdin pipe: %v", err)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatalf("Failed to get stdout pipe: %v", err)
	}
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		t.Fatalf("Failed to start MCP server: %v", err)
	}
	defer func() {
		stdin.Close()
		stdout.Close()
		cmd.Process.Kill()
		cmd.Wait()
	}()

	// Give the server a moment to start
	time.Sleep(200 * time.Millisecond)

	// Create tool request
	request := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "tools/call",
		"params": map[string]interface{}{
			"name":      tool,
			"arguments": testCase.params,
		},
	}

	// Marshal request to JSON
	requestJSON, err := json.Marshal(request)
	if err != nil {
		t.Fatalf("Failed to marshal request: %v", err)
	}

	// Send request
	if _, err := stdin.Write(append(requestJSON, '\n')); err != nil {
		t.Fatalf("Failed to write request: %v", err)
	}

	// Read response
	var response map[string]interface{}
	dec := json.NewDecoder(stdout)
	if err := dec.Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	// Check for error in response
	if errObj, ok := response["error"]; ok {
		if testCase.expectedErr != nil {
			assert.Error(t, fmt.Errorf("%v", errObj))
		} else {
			t.Fatalf("Unexpected error: %v", errObj)
		}
		return
	}

	// Get result from response
	result, ok := response["result"].(map[string]interface{})
	if !ok {
		t.Fatalf("Invalid response format: %v", response)
	}

	// Get content from result
	content, ok := result["content"].([]interface{})
	if !ok {
		t.Fatalf("Invalid result format: %v", result)
	}

	// Get text from content
	if len(content) == 0 {
		t.Fatalf("Empty content in result: %v", result)
	}

	contentItem, ok := content[0].(map[string]interface{})
	if !ok {
		t.Fatalf("Invalid content format: %v", content[0])
	}

	text, ok := contentItem["text"].(string)
	if !ok {
		t.Fatalf("Invalid text format: %v", contentItem)
	}

	// Verify command output
	assert.Contains(t, text, testCase.expectedCmd, "command should match expected")
}

// Main test functions
func TestDeviceCommands(t *testing.T) {
	cases := exampleCommands["device-list"]
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			runCommandTest(t, "device-list", tc)
		})
	}
}

func TestDeviceLogsCommands(t *testing.T) {
	cases := exampleCommands["device-logs"]
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			runCommandTest(t, "device-logs", tc)
		})
	}
}

func TestFleetCommands(t *testing.T) {
	cases := exampleCommands["fleet-list"]
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			runCommandTest(t, "fleet-list", tc)
		})
	}
}

func TestOSVersionsCommands(t *testing.T) {
	cases := exampleCommands["os-versions"]
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			runCommandTest(t, "os-versions", tc)
		})
	}
}

func TestReleaseCommands(t *testing.T) {
	cases := exampleCommands["release-list"]
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			runCommandTest(t, "release-list", tc)
		})
	}
}

func TestReleaseInfoCommands(t *testing.T) {
	cases := exampleCommands["release-info"]
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			runCommandTest(t, "release-info", tc)
		})
	}
}

func TestVersionCommands(t *testing.T) {
	cases := exampleCommands["version"]
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			runCommandTest(t, "version", tc)
		})
	}
}

func TestTagCommands(t *testing.T) {
	cases := exampleCommands["tag-list"]
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			runCommandTest(t, "tag-list", tc)
		})
	}
}
