package main

import (
	"fmt"
	"runtime"
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



func getSystemMemoryInfoHandler(ctx context.Context, request mcp.CallToolRequest)( *mcp.CallToolResult, error){
	// systemInfo, err := request.RequireString("system information")
	// if err != nil {
	// 	return mcp.NewToolResultError(err.Error()), nil
	// }
	
	var memory runtime.MemStats

	// Passes the pointer so the function can modify the original variable
	runtime.ReadMemStats(&memory)

	return mcp.NewToolResultText(fmt.Sprintf("currently allocated bytes: %d!, total system memory: %d, number of garbage collections %d, number of objects on heap: %d", memory.Alloc, memory.Sys / 1024, memory.NumGC, memory.HeapObjects)), nil
	
}

func getSystemCPUInfoHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error){
	cpuInformation := fmt.Sprintf("Number of logical CPUs: %d, Operating System: %s, Architecture: %s, Go Version: %s, Compiler: %s, , Max OS threads: %d", 
		runtime.NumCPU(), runtime.GOOS, runtime.GOARCH, runtime.Version(), runtime.Compiler, runtime.GOMAXPROCS(0))
	
	return mcp.NewToolResultText(cpuInformation), nil	
}

func main() {

	serverMcp := server.NewMCPServer(
		"MCP server Golang",
		"1.0.0",
		server.WithToolCapabilities(false),
	)

	toolHelloWorld := mcp.NewTool("hello_world",
		mcp.WithDescription("Say hello to someone"),
		mcp.WithString("name",
			mcp.Required(),
			mcp.Description("Name of the person to greet"),
		),
	)

	toolGetSystemMemoryInfo := mcp.NewTool("get_system_memory_info",
		mcp.WithDescription("Get detailed system memory usage including allocated bytes, heap objects, and garbage collection stats"),

	)	

	toolGetSystemCPUInfo := mcp.NewTool("get_system_cpu_info",
		mcp.WithDescription("Get detailed system CPU information including number of logical CPUs, OS, architecture, Go version and compiler"),
	)

	
	serverMcp.AddTool(toolHelloWorld, helloHandler)
	serverMcp.AddTool(toolGetSystemMemoryInfo, getSystemMemoryInfoHandler)
	serverMcp.AddTool(toolGetSystemCPUInfo, getSystemCPUInfoHandler)


	// Start the stdio server
	if err := server.ServeStdio(serverMcp); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}

}
