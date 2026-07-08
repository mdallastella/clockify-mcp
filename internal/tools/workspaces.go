package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"clockify-mcp/internal/clockify"
)

// RegisterWorkspaceTools registers all workspace-related MCP tools.
func RegisterWorkspaceTools(s *server.MCPServer, client *clockify.Client) {
	// list_workspaces
	s.AddTool(
		mcp.NewTool("list_workspaces",
			mcp.WithDescription("Get all Clockify workspaces for the authenticated user."),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			workspaces, err := client.ListWorkspaces()
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("failed to list workspaces: %v", err)), nil
			}
			data, err := json.MarshalIndent(workspaces, "", "  ")
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("failed to encode response: %v", err)), nil
			}
			return mcp.NewToolResultText(string(data)), nil
		},
	)

	// get_workspace
	s.AddTool(
		mcp.NewTool("get_workspace",
			mcp.WithDescription("Get a single Clockify workspace by ID."),
			mcp.WithString("workspace_id",
				mcp.Required(),
				mcp.Description("The workspace ID."),
			),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			workspaceID, err := req.RequireString("workspace_id")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			workspace, err := client.GetWorkspace(workspaceID)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("failed to get workspace: %v", err)), nil
			}
			data, err := json.MarshalIndent(workspace, "", "  ")
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("failed to encode response: %v", err)), nil
			}
			return mcp.NewToolResultText(string(data)), nil
		},
	)
}
