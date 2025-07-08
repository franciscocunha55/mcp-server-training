package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"

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

	// Passes the pointer so the function can modify the original variable
	runtime.ReadMemStats(&memory)

	return mcp.NewToolResultText(fmt.Sprintf("currently allocated bytes: %d!, total system memory: %d, number of garbage collections %d, number of objects on heap: %d", memory.Alloc, memory.Sys/1024, memory.NumGC, memory.HeapObjects)), nil
}

func getSystemCPUInfoHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	cpuInformation := fmt.Sprintf("Number of logical CPUs: %d, Operating System: %s, Architecture: %s, Go Version: %s, Compiler: %s, Max OS threads: %d",
		runtime.NumCPU(), runtime.GOOS, runtime.GOARCH, runtime.Version(), runtime.Compiler, runtime.GOMAXPROCS(0))

	return mcp.NewToolResultText(cpuInformation), nil
}

func readFileHandler(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {

	filePath := "/Users/francisco.a.cunha/Downloads/Ollama_fundamentals.pdf"

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

func main() {

	mode := "stdio"
	if len(os.Args) > 1 && os.Args[1] == "http" {
		mode = "http"
	}

	serverMcp := server.NewMCPServer(
		"MCP server Golang",
		"1.0.0",
		server.WithToolCapabilities(true),
		server.WithResourceCapabilities(true, true),
	)

	// Tools
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

	// Resources
	resourcePDFOllama := mcp.NewResource(
		"file:///Users/francisco.a.cunha/Downloads/Ollama_fundamentals.pdf",
		"Ollama Fundamentals PDF",
		mcp.WithResourceDescription("PDF file containing the fundamentals of Ollama."),
		mcp.WithMIMEType("application/pdf"),
	)

	serverMcp.AddResource(resourcePDFOllama, readFileHandler)

	if mode == "http" {
		httpServer := server.NewStreamableHTTPServer(serverMcp)
		log.Printf("Starting MCP server on HTTP at port :8080...")
		log.Printf("MCP endpoint will be available at: http://localhost:8080/mcp")
		if err := httpServer.Start(":8080"); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	} else {
		fmt.Fprintln(os.Stderr, "Starting MCP server on stdio...")
		if err := server.ServeStdio(serverMcp); err != nil {
			fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
		}
	}
}
