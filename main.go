package main

import (
	"fmt"
	serverTest "github.com/franciscocunha55/mcp-serverTest-training/serverTest"

	"context"

	"log"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {

	fmt.Println("Hello World")

	mcpServer, err := serverTest.NewMCPSystemInfoServer("MCPServer-test", "1.0.0", 8080, []string{}, true, time.Now(), map[string]string{"configKey": "configValue"})
	if err != nil {
		fmt.Println("Error creating MCP serverTest:", err)
		return
	}
	fmt.Println(mcpServer.GetInfo())

	errChangeName := mcpServer.SetName("MCPServer-Updated")
	if errChangeName != nil {
		fmt.Println("Error:", errChangeName)
		return
	}
	fmt.Println(mcpServer.GetInfo())

	if status, startTime, err := mcpServer.GetStatus(); err != nil {
		fmt.Println("Error getting status:", err)
		return
	} else {
		fmt.Printf("Server Status: %s, Start Time: %s\n", status, startTime)
	}

	if config, err := mcpServer.GetConfig(); err != nil {
		fmt.Println("Error getting config:", err)
	} else {
		fmt.Println("Server Config:", config)
	}

	if err := mcpServer.Stop(); err != nil {
		fmt.Println("Error stopping serverTest:", err)
	} else {
		fmt.Println("Server stopped successfully")
	}
	if status, _, err := mcpServer.GetStatus(); err != nil {
		fmt.Println("Error getting status:", err)
		return
	} else {
		fmt.Printf("Server Status: %s\n", status)
	}

	serverMcp := server.NewMcpServer(
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
	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}

}

func helloHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name, err := request.RequireString("name")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Hello, %s!", name)), nil
}
