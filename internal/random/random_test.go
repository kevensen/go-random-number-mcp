package random

import (
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
			desc: "valid request with min excluded",
			request: mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Arguments: map[string]any{
						"min":        int64(3),
						"max":        int64(7),
						"includeMin": false,
					},
				},
			},
			min: 4,
			max: 7,
		},
		{
			desc: "valid request with max excluded",
			request: mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Arguments: map[string]any{
						"min":        int64(3),
						"max":        int64(7),
						"includeMax": false,
					},
				},
			},
			min: 3,
			max: 6,
		},
		{
			desc: "valid request with min and max excluded",
			request: mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Arguments: map[string]any{
						"min":        int64(3),
						"max":        int64(7),
						"includeMin": false,
						"includeMax": false,
					},
				},
			},
			min: 4,
			max: 6,
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
		{
			desc: "invalid request with min excluded at max boundary",
			request: mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Arguments: map[string]any{
						"min":        int64(math.MaxInt64),
						"includeMin": false,
					},
				},
			},
			min:     math.MaxInt64,
			max:     math.MaxInt64,
			wantErr: true,
		},
		{
			desc: "invalid request with max excluded at min boundary",
			request: mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Arguments: map[string]any{
						"max":        int64(math.MinInt64),
						"includeMax": false,
					},
				},
			},
			min:     math.MinInt64,
			max:     math.MinInt64,
			wantErr: true,
		},
	}

	ctx := t.Context()
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
	if _, ok := tools["random_float"]; !ok {
		t.Fatalf("NewMCPServer() missing random_float tool")
	}
}

func TestRandomFloatHandler(t *testing.T) {
	testCases := []struct {
		desc        string
		request     mcp.CallToolRequest
		min         float64
		max         float64
		includeMin  bool
		includeMax  bool
		minProvided bool
		maxProvided bool
		wantErr     bool
	}{
		{
			desc:        "valid request with no args",
			request:     mcp.CallToolRequest{},
			min:         0,
			max:         math.MaxFloat64,
			includeMin:  true,
			includeMax:  true,
			minProvided: false,
			maxProvided: false,
		},
		{
			desc: "valid request with min only",
			request: mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Arguments: map[string]any{
						"min": 1.5,
					},
				},
			},
			min:         1.5,
			max:         math.MaxFloat64,
			includeMin:  true,
			includeMax:  true,
			minProvided: true,
			maxProvided: false,
		},
		{
			desc: "valid request with max only",
			request: mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Arguments: map[string]any{
						"max": 9.5,
					},
				},
			},
			min:         0,
			max:         9.5,
			includeMin:  true,
			includeMax:  true,
			minProvided: false,
			maxProvided: true,
		},
		{
			desc: "valid request with min and max",
			request: mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Arguments: map[string]any{
						"min": 2.5,
						"max": 7.5,
					},
				},
			},
			min:         2.5,
			max:         7.5,
			includeMin:  true,
			includeMax:  true,
			minProvided: true,
			maxProvided: true,
		},
		{
			desc: "valid request with min excluded",
			request: mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Arguments: map[string]any{
						"min":        2.5,
						"max":        7.5,
						"includeMin": false,
					},
				},
			},
			min:         2.5,
			max:         7.5,
			includeMin:  false,
			includeMax:  true,
			minProvided: true,
			maxProvided: true,
		},
		{
			desc: "valid request with max excluded",
			request: mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Arguments: map[string]any{
						"min":        2.5,
						"max":        7.5,
						"includeMax": false,
					},
				},
			},
			min:         2.5,
			max:         7.5,
			includeMin:  true,
			includeMax:  false,
			minProvided: true,
			maxProvided: true,
		},
		{
			desc: "valid request with min and max excluded",
			request: mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Arguments: map[string]any{
						"min":        2.5,
						"max":        7.5,
						"includeMin": false,
						"includeMax": false,
					},
				},
			},
			min:         2.5,
			max:         7.5,
			includeMin:  false,
			includeMax:  false,
			minProvided: true,
			maxProvided: true,
		},
		{
			desc: "invalid request with min greater than max",
			request: mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Arguments: map[string]any{
						"min": 9.5,
						"max": 2.5,
					},
				},
			},
			min:         9.5,
			max:         2.5,
			includeMin:  true,
			includeMax:  true,
			minProvided: true,
			maxProvided: true,
			wantErr:     true,
		},
		{
			desc: "invalid request with equal bounds and excluded min",
			request: mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Arguments: map[string]any{
						"min":        3.5,
						"max":        3.5,
						"includeMin": false,
					},
				},
			},
			min:         3.5,
			max:         3.5,
			includeMin:  false,
			includeMax:  true,
			minProvided: true,
			maxProvided: true,
			wantErr:     true,
		},
		{
			desc: "invalid request with equal bounds and excluded max",
			request: mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Arguments: map[string]any{
						"min":        3.5,
						"max":        3.5,
						"includeMax": false,
					},
				},
			},
			min:         3.5,
			max:         3.5,
			includeMin:  true,
			includeMax:  false,
			minProvided: true,
			maxProvided: true,
			wantErr:     true,
		},
		{
			desc: "valid request with equal bounds and both included",
			request: mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Arguments: map[string]any{
						"min": 3.5,
						"max": 3.5,
					},
				},
			},
			min:         3.5,
			max:         3.5,
			includeMin:  true,
			includeMax:  true,
			minProvided: true,
			maxProvided: true,
		},
	}

	ctx := t.Context()
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			result, err := randomFloatHandler(ctx, tc.request)
			if err != nil {
				t.Fatalf("randomFloatHandler() error = %v", err)
			}
			if result == nil || len(result.Content) == 0 {
				t.Fatalf("randomFloatHandler() result is nil or empty")
			}
			if tc.wantErr {
				if !result.IsError {
					t.Fatalf("randomFloatHandler() expected error, got success")
				}
				return
			}
			if result.IsError {
				t.Fatalf("randomFloatHandler() returned error content: %+v", result.Content[0])
			}

			textContent, ok := result.Content[0].(mcp.TextContent)
			if !ok {
				t.Fatalf("randomFloatHandler() content type = %T, want TextContent", result.Content[0])
			}

			valueFromText, err := strconv.ParseFloat(textContent.Text, 64)
			if err != nil {
				t.Fatalf("randomFloatHandler() invalid text content: %v", err)
			}
			if valueFromText < tc.min || valueFromText > tc.max {
				t.Fatalf("randomFloatHandler() text value out of range: %f", valueFromText)
			}
			if tc.minProvided && !tc.includeMin && valueFromText <= tc.min {
				t.Fatalf("randomFloatHandler() expected value > min, got %f", valueFromText)
			}
			if tc.maxProvided && !tc.includeMax && valueFromText >= tc.max {
				t.Fatalf("randomFloatHandler() expected value < max, got %f", valueFromText)
			}

			structured, ok := result.StructuredContent.(randomFloatResponse)
			if !ok {
				t.Fatalf("randomFloatHandler() structured content type = %T, want randomFloatResponse", result.StructuredContent)
			}
			if structured.Value != valueFromText {
				t.Fatalf("randomFloatHandler() structured value %f != text value %f", structured.Value, valueFromText)
			}
		})
	}
}
