package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/authentiq-api/mcp-server/config"
	"github.com/authentiq-api/mcp-server/models"
	"github.com/mark3labs/mcp-go/mcp"
)

func Key_revoke_nosecretHandler(cfg *config.APIConfig) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, ok := request.Params.Arguments.(map[string]any)
		if !ok {
			return mcp.NewToolResultError("Invalid arguments object"), nil
		}
		queryParams := make([]string, 0)
		if val, ok := args["email"]; ok {
			queryParams = append(queryParams, fmt.Sprintf("email=%v", val))
		}
		if val, ok := args["phone"]; ok {
			queryParams = append(queryParams, fmt.Sprintf("phone=%v", val))
		}
		if val, ok := args["code"]; ok {
			queryParams = append(queryParams, fmt.Sprintf("code=%v", val))
		}
		queryString := ""
		if len(queryParams) > 0 {
			queryString = "?" + strings.Join(queryParams, "&")
		}
		url := fmt.Sprintf("%s/key%s", cfg.BaseURL, queryString)
		req, err := http.NewRequest("DELETE", url, nil)
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

func CreateKey_revoke_nosecretTool(cfg *config.APIConfig) models.Tool {
	tool := mcp.NewTool("delete_key",
		mcp.WithDescription("Revoke an Authentiq ID using email & phone.

If called with `email` and `phone` only, a verification code 
will be sent by email. Do a second call adding `code` to 
complete the revocation.
"),
		mcp.WithString("email", mcp.Required(), mcp.Description("primary email associated to Key (ID)")),
		mcp.WithString("phone", mcp.Required(), mcp.Description("primary phone number, international representation")),
		mcp.WithString("code", mcp.Description("verification code sent by email")),
	)

	return models.Tool{
		Definition: tool,
		Handler:    Key_revoke_nosecretHandler(cfg),
	}
}
