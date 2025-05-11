package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/klutchell/balenamcp/server"
	mcpserver "github.com/mark3labs/mcp-go/server"
)

func main() {
	// Parse command line flags
	dryRun := flag.Bool("dry-run", false, "Enable dry-run mode")
	flag.Parse()

	// Set dry-run mode
	server.Config.DryRun = *dryRun

	// Create and setup MCP server
	srv := server.SetupServer()

	fmt.Fprintf(os.Stderr, "Starting BalenaMCP server...\n")
	if err := mcpserver.ServeStdio(srv); err != nil {
		fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
	}
}
