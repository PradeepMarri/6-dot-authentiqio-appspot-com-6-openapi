package main

import (
	"github.com/authentiq-api/mcp-server/config"
	"github.com/authentiq-api/mcp-server/models"
	tools_key "github.com/authentiq-api/mcp-server/tools/key"
	tools_scope "github.com/authentiq-api/mcp-server/tools/scope"
)

func GetAll(cfg *config.APIConfig) []models.Tool {
	return []models.Tool{
		tools_key.CreateKey_revoke_nosecretTool(cfg),
		tools_key.CreateKey_revokeTool(cfg),
		tools_key.CreateKey_retrieveTool(cfg),
		tools_key.CreateHead_key_pkTool(cfg),
		tools_scope.CreateSign_confirmTool(cfg),
		tools_scope.CreateSign_deleteTool(cfg),
		tools_scope.CreateSign_retrieveTool(cfg),
		tools_scope.CreateSign_retrieve_headTool(cfg),
	}
}
