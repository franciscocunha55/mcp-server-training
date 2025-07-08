package main

import (
"encoding/base64"
"fmt"
"io"
"os"
"runtime"
"strings"

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

func getSystemMemoryInfoHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var memory runtime.MemStats
	runtime.ReadMemStats(&memory)
	return mcp.NewToolResultText(fmt.Sprintf("currently allocated bytes: %d!, total system memory: %d, number of garbage collections %d, number of objects on heap: %d", memory.Alloc, memory.Sys/1024, memory.NumGC, memory.HeapObjects)), nil
}

func getSystemCPUInfoHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	cpuInformation := fmt.Sprintf("Number of logical CPUs: %d, Operating System: %s, Architecture: %s, Go Version: %s, Compiler: %s, Max OS threads: %d",
runtime.NumCPU(), runtime.GOOS, runtime.GOARCH, runtime.Version(), runtime.Compiler, runtime.GOMAXPROCS(0))
	return mcp.NewToolResultText(cpuInformation), nil
}

func readFileHandler(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	filePath := strings.Replace(request.Params.URI, "file:///", "", 1)
	openPDF, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer openPDF.Close()
	binaryData, err := io.ReadAll(openPDF)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	base64EncodedData := base64.StdEncoding.EncodeToString(binaryData)

	return []mcp.ResourceContents{
		&mcp.BlobResourceContents{
			URI:      request.Params.URI,
			MIMEType: "application/pdf",
			Blob:     base64EncodedData,
		},
	}, nil
}

func readTextFileHandler(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	filePath := strings.Replace(request.Params.URI, "file:///", "", 1)
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return []mcp.ResourceContents{
		&mcp.TextResourceContents{
			URI:      request.Params.URI,
			MIMEType: "text/plain",
			Text:     string(data),
		},
	}, nil
}

func main() {
	serverMcp := server.NewMCPServer(
"MCP server Golang",
"1.0.0",
server.WithToolCapabilities(true),
server.WithResourceCapabilities(true, true),
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

	resourcePDFOllama := mcp.NewResource(
"file:///Users/francisco.a.cunha/Downloads/Ollama_fundamentals.pdf",
"Ollama Fundamentals PDF",
mcp.WithResourceDescription("PDF file containing the fundamentals of Ollama."),
mcp.WithMIMEType("application/pdf"),
)

	resourceTextTest := mcp.NewResource(
"file:///Users/francisco.a.cunha/Downloads/ollama_test.txt",
"Ollama Test Text",
mcp.WithResourceDescription("Test text file with Ollama information for MCP testing."),
mcp.WithMIMEType("text/plain"),
)

	serverMcp.AddResource(resourcePDFOllama, readFileHandler)
	serverMcp.AddResource(resourceTextTest, readTextFileHandler)

	if err := server.ServeStdio(serverMcp); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
