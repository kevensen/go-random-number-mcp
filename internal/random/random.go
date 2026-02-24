package random

import (
	"context"
	"crypto/rand"
	"fmt"
	"log/slog"
	"math"
	"math/big"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type randomIntResponse struct {
	Value int64 `json:"value"`
}

type randomIntArgs struct {
	Min        *int64 `json:"min,omitempty"`
	Max        *int64 `json:"max,omitempty"`
	IncludeMin *bool  `json:"includeMin,omitempty"`
	IncludeMax *bool  `json:"includeMax,omitempty"`
}

type randomFloatResponse struct {
	Value float64 `json:"value"`
}

type randomFloatArgs struct {
	Min        *float64 `json:"min,omitempty"`
	Max        *float64 `json:"max,omitempty"`
	IncludeMin *bool    `json:"includeMin,omitempty"`
	IncludeMax *bool    `json:"includeMax,omitempty"`
}

type randomASCIIResponse struct {
	Value string `json:"value"`
}

type randomASCIIArgs struct {
	Length int `json:"length"`
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
		mcp.WithDescription("Returns a cryptographically secure random integer. Optional arguments: min, max, includeMin, includeMax."),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithInputSchema[randomIntArgs](),
		mcp.WithOutputSchema[randomIntResponse](),
	)

	mcpServer.AddTool(tool, randomIntHandler)

	floatTool := mcp.NewTool(
		"random_float",
		mcp.WithDescription("Returns a cryptographically secure random floating-point number. Optional arguments: min, max, includeMin, includeMax."),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithInputSchema[randomFloatArgs](),
		mcp.WithOutputSchema[randomFloatResponse](),
	)

	mcpServer.AddTool(floatTool, randomFloatHandler)

	stringTool := mcp.NewTool(
		"random_ascii",
		mcp.WithDescription("Returns a cryptographically secure random ASCII string. Required argument: length."),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithInputSchema[randomASCIIArgs](),
		mcp.WithOutputSchema[randomASCIIResponse](),
	)

	mcpServer.AddTool(stringTool, randomASCIIHandler)

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
	includeMin := true
	includeMax := true
	if args.Min != nil {
		min = *args.Min
	}
	if args.Max != nil {
		max = *args.Max
	}
	if args.IncludeMin != nil {
		includeMin = *args.IncludeMin
	}
	if args.IncludeMax != nil {
		includeMax = *args.IncludeMax
	}

	adjustedMin := min
	adjustedMax := max
	if args.Min != nil && !includeMin {
		if min == math.MaxInt64 {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.TextContent{Type: "text", Text: "random_int failed: min cannot be excluded when min is MaxInt64"},
				},
			}, nil
		}
		adjustedMin = min + 1
	}
	if args.Max != nil && !includeMax {
		if max == math.MinInt64 {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.TextContent{Type: "text", Text: "random_int failed: max cannot be excluded when max is MinInt64"},
				},
			}, nil
		}
		adjustedMax = max - 1
	}

	slog.InfoContext(ctx, "randomIntHandler", slog.Int64("min", min), slog.Int64("max", max), slog.Bool("includeMin", includeMin), slog.Bool("includeMax", includeMax))
	value, err := randomInt64InRange(adjustedMin, adjustedMax)
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

func randomFloatHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var args randomFloatArgs
	if err := request.BindArguments(&args); err != nil {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				mcp.TextContent{Type: "text", Text: fmt.Sprintf("random_float failed: %v", err)},
			},
		}, nil
	}

	min := 0.0
	max := math.MaxFloat64
	includeMin := true
	includeMax := true
	if args.Min != nil {
		min = *args.Min
	}
	if args.Max != nil {
		max = *args.Max
	}
	if args.IncludeMin != nil {
		includeMin = *args.IncludeMin
	}
	if args.IncludeMax != nil {
		includeMax = *args.IncludeMax
	}

	value, err := randomFloat64InRange(min, max, includeMin, includeMax, args.Min != nil, args.Max != nil)
	if err != nil {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				mcp.TextContent{Type: "text", Text: fmt.Sprintf("random_float failed: %v", err)},
			},
		}, nil
	}

	response := randomFloatResponse{Value: value}
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{Type: "text", Text: fmt.Sprintf("%g", value)},
		},
		StructuredContent: response,
	}, nil
}

func randomASCIIHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var args randomASCIIArgs
	if err := request.BindArguments(&args); err != nil {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				mcp.TextContent{Type: "text", Text: fmt.Sprintf("random_ascii failed: %v", err)},
			},
		}, nil
	}

	value, err := randomASCIIString(args.Length)
	if err != nil {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				mcp.TextContent{Type: "text", Text: fmt.Sprintf("random_ascii failed: %v", err)},
			},
		}, nil
	}

	response := randomASCIIResponse{Value: value}
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{Type: "text", Text: value},
		},
		StructuredContent: response,
	}, nil
}

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

func randomFloat64InRange(min, max float64, includeMin, includeMax, hasMin, hasMax bool) (float64, error) {
	if math.IsNaN(min) || math.IsNaN(max) {
		return 0, fmt.Errorf("min and max must not be NaN")
	}
	if math.IsInf(min, 0) || math.IsInf(max, 0) {
		return 0, fmt.Errorf("min and max must be finite")
	}
	if min > max {
		return 0, fmt.Errorf("min cannot be greater than max")
	}
	if min == max {
		if includeMin && includeMax {
			return min, nil
		}
		return 0, fmt.Errorf("range is empty when min equals max and is excluded")
	}

	adjustedMin := min
	adjustedMax := max
	if hasMin && !includeMin {
		adjustedMin = math.Nextafter(min, math.Inf(1))
	}
	if hasMax && !includeMax {
		adjustedMax = math.Nextafter(max, math.Inf(-1))
	}
	if adjustedMin > adjustedMax {
		return 0, fmt.Errorf("range is empty after applying exclusivity")
	}

	unit, err := cryptoRandFloat64()
	if err != nil {
		return 0, err
	}

	return adjustedMin + unit*(adjustedMax-adjustedMin), nil
}

func cryptoRandFloat64() (float64, error) {
	const maxUint53 = 1 << 53
	value, err := rand.Int(rand.Reader, big.NewInt(maxUint53))
	if err != nil {
		return 0, err
	}
	return float64(value.Int64()) / float64(maxUint53), nil
}

// randomASCIIString returns a cryptographically secure random string of printable ASCII characters.
// Length must be greater than zero.
func randomASCIIString(length int) (string, error) {
	if length <= 0 {
		return "", &ZeroLengthError{}
	}

	const asciiStart = 32
	const asciiEnd = 126
	const asciiRange = asciiEnd - asciiStart + 1

	var builder strings.Builder
	builder.Grow(length)
	max := big.NewInt(asciiRange)
	for i := 0; i < length; i++ {
		value, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", err
		}
		builder.WriteByte(byte(asciiStart + value.Int64()))
	}

	return builder.String(), nil
}
