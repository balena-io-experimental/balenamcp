# BalenaMCP

A simple MCP server that wraps the `balena` CLI to provide Balena capabilities
to supported clients.

## Prerequisites

- Go 1.21 or later
- `balena` CLI installed and available in PATH
- Claude Desktop installed

## Installation

```bash
# Clone the repository
git clone https://github.com/klutchell/balenamcp.git
cd balenamcp

# Install dependencies
go mod download

# Build the binary
go build -o bin/balenamcp
```

## Configuration

### Claude Desktop Setup

1. Open Claude Desktop
2. Go to the Claude menu and select "Settings..."
3. Click on "Developer" in the left-hand bar

   You may need to enable Developer Mode to see this option.

4. Click on "Edit Config"

   This will create or open a configuration file at:

   - macOS: `~/Library/Application Support/Claude/claude_desktop_config.json`
   - Windows: `%APPDATA%\Claude\claude_desktop_config.json`

5. Add the following configuration to the file:

   ```json
   {
     "mcpServers": {
       "balenamcp": {
         "command": "/full/path/to/your/balenamcp/bin/balenamcp",
         "args": []
       }
     }
   }
   ```

   Make sure to replace `/full/path/to/your/balenamcp` with the absolute path to
   your project directory.

6. Save the file and restart Claude Desktop

### Authentication

Before using the tools, make sure you're authenticated with Balena:

```bash
balena login
```

## Available MCP Commands

The following commands are available through BalenaMCP:

- `device-list`: List all devices
- `device-logs`: Show device logs (requires device UUID or name)
- `fleet-list`: List all fleets
- `os-versions`: Show available balenaOS versions for a device type (requires
  device type)
- `release-list`: List all releases of a fleet (requires fleet name)
- `release-info`: Get info for a release (requires release ID)
- `version`: Display version information for the balena CLI
- `tag-list`: List all tags for a fleet, device or release (requires resource
  type and name)

Each command can be accessed directly through Claude Desktop.

## Development

The server is built using the
[MCP Go library](https://github.com/mark3labs/mcp-go) and communicates with
Claude Desktop through standard input/output streams.

### Project Structure

- `main.go`: Main server implementation with tool definitions
- `server/setup.go`: Server setup and tool implementations
- `main_test.go`: Tests for the server implementation
- `bin/`: Built binaries (gitignored)

### Building

```bash
# Build for current platform
go build -o bin/balenamcp

# Build for specific platforms
GOOS=linux GOARCH=amd64 go build -o bin/balenamcp-linux-amd64
GOOS=windows GOARCH=amd64 go build -o bin/balenamcp-windows-amd64.exe
GOOS=darwin GOARCH=amd64 go build -o bin/balenamcp-darwin-amd64
GOOS=darwin GOARCH=arm64 go build -o bin/balenamcp-darwin-arm64
```

### Running Tests

```bash
go test ./...
```

## Troubleshooting

If you encounter issues with the MCP server:

1. Check the logs in:

   - macOS: `~/Library/Logs/Claude/mcp-server-balenamcp.log`
   - Windows: `%APPDATA%\Claude\logs\mcp-server-balenamcp.log`

2. Make sure the path to the balenamcp binary is correct and absolute
3. Ensure you have the necessary permissions to execute the binary
4. Verify that balena is installed and accessible in your PATH

## License

MIT
