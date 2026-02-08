package random

import (
	"context"
	"crypto/rand"
	"fmt"
	"log/slog"
	"math"
	"math/big"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// randomInt64InRange returns a cryptographically secure random integer in the
// inclusive range [min, max].
func randomInt64InRange(min, max int64) (int64, error) {
	minBig := big.NewInt(min)
	maxBig := big.NewInt(max)
	if minBig.Cmp(maxBig) > 0 {
		return 0, fmt.Errorf("min cannot be greater than max")
	}

	rangeSize := new(big.Int).Sub(maxBig, minBig)
	rangeSize.Add(rangeSize, big.NewInt(1))
	value, err := rand.Int(rand.Reader, rangeSize)
	if err != nil {
		return 0, err
	}

	value.Add(value, minBig)
	return value.Int64(), nil
}

type randomIntResponse struct {
	Value int64 `json:"value"`
}

type randomIntArgs struct {
	Min *int64 `json:"min,omitempty"`
	Max *int64 `json:"max,omitempty"`
}

// NewMCPServer builds the MCP server with the random_int tool registered.
func NewMCPServer(name, version string) *server.MCPServer {
	mcpServer := server.NewMCPServer(
		name,
		version,
		server.WithInstructions("Use the random_int tool to get a cryptographically secure random integer."),
	)

	tool := mcp.NewTool(
		"random_int",
		mcp.WithDescription("Returns a cryptographically secure random integer. Optional arguments: min, max."),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithInputSchema[randomIntArgs](),
		mcp.WithOutputSchema[randomIntResponse](),
	)

	mcpServer.AddTool(tool, randomIntHandler)

	return mcpServer
}

func randomIntHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var args randomIntArgs
	if err := request.BindArguments(&args); err != nil {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				mcp.TextContent{Type: "text", Text: fmt.Sprintf("random_int failed: %v", err)},
			},
		}, nil
	}

	min := int64(0)
	max := int64(math.MaxInt64)
	if args.Min != nil {
		min = *args.Min
	}
	if args.Max != nil {
		max = *args.Max
	}
	slog.InfoContext(ctx, "randomIntHandler", slog.Int64("min", min), slog.Int64("max", max))
	value, err := randomInt64InRange(min, max)
	if err != nil {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				mcp.TextContent{Type: "text", Text: fmt.Sprintf("random_int failed: %v", err)},
			},
		}, nil
	}
	slog.InfoContext(ctx, "randomIntHandler", slog.Int64("result", value))

	response := randomIntResponse{Value: value}
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{Type: "text", Text: fmt.Sprintf("%d", value)},
		},
		StructuredContent: response,
	}, nil
}
