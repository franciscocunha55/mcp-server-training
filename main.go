package main

import (
	"fmt"
	//serverTest "github.com/franciscocunha55/mcp-server-training/serverTest"

	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func helloHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name, err := request.RequireString("name")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Hello, %s!", name)), nil
}



func main() {

	serverMcp := server.NewMCPServer(
		"Hello World",
		"1.0.0",
		server.WithToolCapabilities(false),
	)

	tool := mcp.NewTool("hello_world",
		mcp.WithDescription("Say hello to someone"),
		mcp.WithString("name",
			mcp.Required(),
			mcp.Description("Name of the person to greet"),
		),
	)

	serverMcp.AddTool(tool, helloHandler)

	// Start the stdio server
	if err := server.ServeStdio(serverMcp); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}

}
