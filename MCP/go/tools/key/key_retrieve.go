package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/authentiq-api/mcp-server/config"
	"github.com/authentiq-api/mcp-server/models"
	"github.com/mark3labs/mcp-go/mcp"
)

func Key_retrieveHandler(cfg *config.APIConfig) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, ok := request.Params.Arguments.(map[string]any)
		if !ok {
			return mcp.NewToolResultError("Invalid arguments object"), nil
		}
		PKVal, ok := args["PK"]
		if !ok {
			return mcp.NewToolResultError("Missing required path parameter: PK"), nil
		}
		PK, ok := PKVal.(string)
		if !ok {
			return mcp.NewToolResultError("Invalid path parameter: PK"), nil
		}
		url := fmt.Sprintf("%s/key/%s", cfg.BaseURL, PK)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Failed to create request", err), nil
		}
		// No authentication required for this endpoint
		req.Header.Set("Accept", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Request failed", err), nil
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Failed to read response body", err), nil
		}

		if resp.StatusCode >= 400 {
			return mcp.NewToolResultError(fmt.Sprintf("API error: %s", body)), nil
		}
		// Use properly typed response
		var result map[string]interface{}
		if err := json.Unmarshal(body, &result); err != nil {
			// Fallback to raw text if unmarshaling fails
			return mcp.NewToolResultText(string(body)), nil
		}

		prettyJSON, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Failed to format JSON", err), nil
		}

		return mcp.NewToolResultText(string(prettyJSON)), nil
	}
}

func CreateKey_retrieveTool(cfg *config.APIConfig) models.Tool {
	tool := mcp.NewTool("get_key_PK",
		mcp.WithDescription("Get public details of an Authentiq ID.
"),
		mcp.WithString("PK", mcp.Required(), mcp.Description("Public Signing Key - Authentiq ID (43 chars)")),
	)

	return models.Tool{
		Definition: tool,
		Handler:    Key_retrieveHandler(cfg),
	}
}
