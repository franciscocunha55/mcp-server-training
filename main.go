package main

import (
	"fmt"
	"github.com/franciscocunha55/mcp-server-training/server"
	"time"
)

func main() {

	fmt.Println("Hello World")

	mcpServer, err := server.NewMCPSystemInfoServer("MCPServer-test", "1.0.0", 8080, []string{}, true, time.Now(), map[string]string{"configKey": "configValue"})
	if err != nil {
		fmt.Println("Error creating MCP server:", err)
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
		fmt.Println("Error stopping server:", err)
	} else {
		fmt.Println("Server stopped successfully")
	}
	if status, _, err := mcpServer.GetStatus(); err != nil {
		fmt.Println("Error getting status:", err)
		return
	} else {
		fmt.Printf("Server Status: %s\n", status)
	}
}
