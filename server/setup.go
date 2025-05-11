package server

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// ServerConfig holds the server configuration
type ServerConfig struct {
	DryRun bool
}

// Global server configuration
var Config = ServerConfig{
	DryRun: false,
}

// executeCommand executes a balena command, either for real or in dry-run mode
func executeCommand(args []string) (string, error) {
	if Config.DryRun {
		cmdStr := "balena " + strings.Join(args, " ")
		fmt.Fprintf(os.Stderr, "[DRY RUN] Would execute: %s\n", cmdStr)
		return fmt.Sprintf("[DRY RUN] %s", cmdStr), nil
	}

	cmd := exec.Command("balena", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to execute command: %v\nOutput: %s", err, string(output))
	}
	return string(output), nil
}

// SetupServer creates and configures an MCP server with all tools
func SetupServer() *server.MCPServer {
	// Create MCP server
	srv := server.NewMCPServer(
		"BalenaMCP",
		"1.0.0",
		server.WithLogging(),
		server.WithToolCapabilities(true),
	)

	// Register tools
	srv.AddTool(mcp.NewTool("device-list",
		mcp.WithDescription("List all devices"),
		mcp.WithBoolean("help",
			mcp.Description("Show help for this command"),
		),
		mcp.WithString("fleet",
			mcp.Description("fleet name or slug (preferred)"),
		),
		mcp.WithBoolean("json",
			mcp.Description("produce JSON output instead of tabular output"),
		),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {

		// Extract arguments
		var args []string
		args = append(args, "device", "list")

		// Extract flags
		var flags []string
		if json, ok := request.Params.Arguments["json"].(bool); ok && json {
			flags = append(flags, "--json")
		}
		if fleet, ok := request.Params.Arguments["fleet"].(string); ok {
			flags = append(flags, "--fleet", fleet)
		}
		if help, ok := request.Params.Arguments["help"].(bool); ok && help {
			flags = append(flags, "--help")
		}

		output, err := executeCommand(append(args, flags...))
		if err != nil {
			return nil, err
		}

		return mcp.NewToolResultText(output), nil
	})

	// Register device-logs capability
	srv.AddTool(mcp.NewTool("device-logs",
		mcp.WithDescription("Show device logs"),
		mcp.WithBoolean("help",
			mcp.Description("Show help for this command"),
		),
		mcp.WithBoolean("system",
			mcp.Description("Display system logs only"),
		),
		mcp.WithBoolean("tail",
			mcp.Description("continuously stream output"),
		),
		mcp.WithString("service",
			mcp.Description("Service name to filter logs"),
		),
		mcp.WithString("device",
			mcp.Description("UUID or name of the device"),
		),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Extract arguments
		var args []string
		args = append(args, "device", "logs")
		if device, ok := request.Params.Arguments["device"].(string); ok {
			args = append(args, device)
		}

		// Extract flags
		var flags []string
		if help, ok := request.Params.Arguments["help"].(bool); ok && help {
			flags = append(flags, "--help")
		}
		if system, ok := request.Params.Arguments["system"].(bool); ok && system {
			flags = append(flags, "--system")
		}
		if tail, ok := request.Params.Arguments["tail"].(bool); ok && tail {
			flags = append(flags, "--tail")
		}
		if service, ok := request.Params.Arguments["service"].(string); ok {
			flags = append(flags, "--service", service)
		}

		output, err := executeCommand(append(args, flags...))
		if err != nil {
			return nil, err
		}

		return mcp.NewToolResultText(output), nil
	})

	// Register fleet-list capability
	srv.AddTool(mcp.NewTool("fleet-list",
		mcp.WithDescription("List all fleets"),
		mcp.WithBoolean("help",
			mcp.Description("Show help for this command"),
		),
		mcp.WithBoolean("json",
			mcp.Description("produce JSON output instead of tabular output"),
		),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Extract arguments
		var args []string
		args = append(args, "fleet", "list")

		// Extract flags
		var flags []string
		if help, ok := request.Params.Arguments["help"].(bool); ok && help {
			flags = append(flags, "--help")
		}
		if json, ok := request.Params.Arguments["json"].(bool); ok && json {
			flags = append(flags, "--json")
		}

		output, err := executeCommand(append(args, flags...))
		if err != nil {
			return nil, err
		}

		return mcp.NewToolResultText(output), nil
	})

	// Register os-versions capability
	srv.AddTool(mcp.NewTool("os-versions",
		mcp.WithDescription("Show available balenaOS versions for the given device type"),
		mcp.WithBoolean("help",
			mcp.Description("Show help for this command"),
		),
		mcp.WithBoolean("esr",
			mcp.Description("select balenaOS ESR versions"),
		),
		mcp.WithBoolean("include_draft",
			mcp.Description("include pre-release balenaOS versions"),
		),
		mcp.WithString("type",
			mcp.Description("device type"),
		),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Extract arguments
		var args []string
		args = append(args, "os", "versions")
		if deviceType, ok := request.Params.Arguments["type"].(string); ok {
			args = append(args, deviceType)
		}

		// Extract flags
		var flags []string
		if help, ok := request.Params.Arguments["help"].(bool); ok && help {
			flags = append(flags, "--help")
		}
		if esr, ok := request.Params.Arguments["esr"].(bool); ok && esr {
			flags = append(flags, "--esr")
		}
		if includeDraft, ok := request.Params.Arguments["include_draft"].(bool); ok && includeDraft {
			flags = append(flags, "--include-draft")
		}

		output, err := executeCommand(append(args, flags...))
		if err != nil {
			return nil, err
		}

		return mcp.NewToolResultText(output), nil
	})

	// Register release-list capability
	srv.AddTool(mcp.NewTool("release-list",
		mcp.WithDescription("List all releases of a fleet"),
		mcp.WithBoolean("help",
			mcp.Description("Show help for this command"),
		),
		mcp.WithBoolean("json",
			mcp.Description("produce JSON output instead of tabular output"),
		),
		mcp.WithString("fleet",
			mcp.Description("Name of the fleet"),
		),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Extract arguments
		var args []string
		args = append(args, "release", "list")
		if fleet, ok := request.Params.Arguments["fleet"].(string); ok {
			args = append(args, fleet)
		}

		// Extract flags
		var flags []string
		if help, ok := request.Params.Arguments["help"].(bool); ok && help {
			flags = append(flags, "--help")
		}
		if json, ok := request.Params.Arguments["json"].(bool); ok && json {
			flags = append(flags, "--json")
		}

		output, err := executeCommand(append(args, flags...))
		if err != nil {
			return nil, err
		}

		return mcp.NewToolResultText(output), nil
	})

	// Register release-info capability
	srv.AddTool(mcp.NewTool("release-info",
		mcp.WithDescription("Get info for a release"),
		mcp.WithBoolean("help",
			mcp.Description("Show help for this command"),
		),
		mcp.WithBoolean("json",
			mcp.Description("produce JSON output instead of tabular output"),
		),
		mcp.WithBoolean("composition",
			mcp.Description("Return the release composition"),
		),
		mcp.WithString("id",
			mcp.Description("ID of the release"),
		),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Extract arguments
		var args []string
		args = append(args, "release")
		if id, ok := request.Params.Arguments["id"].(string); ok {
			args = append(args, id)
		}

		// Extract flags
		var flags []string
		if help, ok := request.Params.Arguments["help"].(bool); ok && help {
			flags = append(flags, "--help")
		}
		if json, ok := request.Params.Arguments["json"].(bool); ok && json {
			flags = append(flags, "--json")
		}
		if composition, ok := request.Params.Arguments["composition"].(bool); ok && composition {
			flags = append(flags, "--composition")
		}

		output, err := executeCommand(append(args, flags...))
		if err != nil {
			return nil, err
		}

		return mcp.NewToolResultText(output), nil
	})

	// Register version capability
	srv.AddTool(mcp.NewTool("version",
		mcp.WithDescription("Display version information for the balena CLI"),
		mcp.WithBoolean("help",
			mcp.Description("Show help for this command"),
		),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Extract arguments
		var args []string
		args = append(args, "version")

		// Extract flags
		var flags []string
		if help, ok := request.Params.Arguments["help"].(bool); ok && help {
			flags = append(flags, "--help")
		}

		output, err := executeCommand(append(args, flags...))
		if err != nil {
			return nil, err
		}

		return mcp.NewToolResultText(output), nil
	})

	// Register tag-list capability
	srv.AddTool(mcp.NewTool("tag-list",
		mcp.WithDescription("List all tags for a fleet, device or release"),
		mcp.WithBoolean("help",
			mcp.Description("Show help for this command"),
		),
		mcp.WithString("resource",
			mcp.Description("Resource type (fleet, device, or release)"),
		),
		mcp.WithString("name",
			mcp.Description("Name of the resource (fleet name, device UUID, release ID)"),
		),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Extract arguments
		var args []string
		args = append(args, "tag", "list")
		if resource, ok := request.Params.Arguments["resource"].(string); ok {
			args = append(args, resource)
		}
		if name, ok := request.Params.Arguments["name"].(string); ok {
			args = append(args, name)
		}

		// Extract flags
		var flags []string
		if help, ok := request.Params.Arguments["help"].(bool); ok && help {
			flags = append(flags, "--help")
		}

		output, err := executeCommand(append(args, flags...))
		if err != nil {
			return nil, err
		}

		return mcp.NewToolResultText(output), nil
	})

	return srv
}
