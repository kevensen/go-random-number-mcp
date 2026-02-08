package random

import (
	"context"
	"math"
	"strconv"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
)

func TestRandomIntHandler(t *testing.T) {
	testCases := []struct {
		desc    string
		request mcp.CallToolRequest
		min     int64
		max     int64
		wantErr bool
	}{
		{
			desc:    "valid request with no args",
			request: mcp.CallToolRequest{},
			min:     0,
			max:     math.MaxInt64,
		},
		{
			desc: "valid request with min only",
			request: mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Arguments: map[string]any{
						"min": int64(5),
					},
				},
			},
			min: 5,
			max: math.MaxInt64,
		},
		{
			desc: "valid request with max only",
			request: mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Arguments: map[string]any{
						"max": int64(10),
					},
				},
			},
			min: 0,
			max: 10,
		},
		{
			desc: "valid request with min and max",
			request: mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Arguments: map[string]any{
						"min": int64(3),
						"max": int64(7),
					},
				},
			},
			min: 3,
			max: 7,
		},
		{
			desc: "invalid request with min greater than max",
			request: mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Arguments: map[string]any{
						"min": int64(10),
						"max": int64(5),
					},
				},
			},
			min:     10,
			max:     5,
			wantErr: true,
		},
	}

	ctx := context.Background()
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			result, err := randomIntHandler(ctx, tc.request)
			if err != nil {
				t.Fatalf("randomIntHandler() error = %v", err)
			}
			if result == nil || len(result.Content) == 0 {
				t.Fatalf("randomIntHandler() result is nil or empty")
			}
			if tc.wantErr {
				if !result.IsError {
					t.Fatalf("randomIntHandler() expected error, got success")
				}
				return
			}
			if result.IsError {
				t.Fatalf("randomIntHandler() returned error content: %+v", result.Content[0])
			}

			textContent, ok := result.Content[0].(mcp.TextContent)
			if !ok {
				t.Fatalf("randomIntHandler() content type = %T, want TextContent", result.Content[0])
			}

			valueFromText, err := strconv.ParseInt(textContent.Text, 10, 64)
			if err != nil {
				t.Fatalf("randomIntHandler() invalid text content: %v", err)
			}
			if valueFromText < tc.min || valueFromText > tc.max {
				t.Fatalf("randomIntHandler() text value out of range: %d", valueFromText)
			}

			structured, ok := result.StructuredContent.(randomIntResponse)
			if !ok {
				t.Fatalf("randomIntHandler() structured content type = %T, want randomIntResponse", result.StructuredContent)
			}
			if structured.Value != valueFromText {
				t.Fatalf("randomIntHandler() structured value %d != text value %d", structured.Value, valueFromText)
			}
		})
	}
}

func TestNewMCPServerRegistersTool(t *testing.T) {
	server := NewMCPServer("test-server", "0.0.0")
	tools := server.ListTools()
	if tools == nil {
		t.Fatalf("NewMCPServer() tools list is nil")
	}
	if _, ok := tools["random_int"]; !ok {
		t.Fatalf("NewMCPServer() missing random_int tool")
	}
}
