package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/kevensen/go-random-number-mcp/internal/random"
	"github.com/mark3labs/mcp-go/server"
)

const (
	serverName    = "go-random-number-mcp"
	serverVersion = "0.1.0"
)

func main() {
	listenAddr := flag.String("addr", "127.0.0.1", "Listen address")
	listenPort := flag.Int("port", 6767, "Listen port")
	flag.Parse()

	mcpServer := random.NewMCPServer(serverName, serverVersion)

	streamServer := server.NewStreamableHTTPServer(mcpServer)
	addr := fmt.Sprintf("%s:%d", *listenAddr, *listenPort)
	slog.Info("MCP server listening", slog.String("url", "http://"+addr+"/mcp"))
	if err := streamServer.Start(addr); err != nil {
		slog.Error("unable to start MCP streaming server", slog.Any("error", err))
		os.Exit(1)
	}
}
