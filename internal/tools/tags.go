package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"clockify-mcp/internal/clockify"
)

// RegisterTagTools registers all tag-related MCP tools.
func RegisterTagTools(s *server.MCPServer, client *clockify.Client) {

	// list_tags
	s.AddTool(
		mcp.NewTool("list_tags",
			mcp.WithDescription("Get all tags in a Clockify workspace."),
			mcp.WithString("workspace_id", mcp.Required(), mcp.Description("Workspace ID.")),
			mcp.WithNumber("page", mcp.Description("Page number (1-based).")),
			mcp.WithNumber("page_size", mcp.Description("Number of tags per page.")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			workspaceID, err := req.RequireString("workspace_id")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			page := req.GetInt("page", 1)
			pageSize := req.GetInt("page_size", 50)
			tags, err := client.ListTags(workspaceID, page, pageSize)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("failed to list tags: %v", err)), nil
			}
			data, _ := json.MarshalIndent(tags, "", "  ")
			return mcp.NewToolResultText(string(data)), nil
		},
	)

	// create_tag
	s.AddTool(
		mcp.NewTool("create_tag",
			mcp.WithDescription("Create a new tag in a Clockify workspace."),
			mcp.WithString("workspace_id", mcp.Required(), mcp.Description("Workspace ID.")),
			mcp.WithString("name", mcp.Required(), mcp.Description("Tag name.")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			workspaceID, err := req.RequireString("workspace_id")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			name, err := req.RequireString("name")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			tag, err := client.CreateTag(workspaceID, clockify.TagRequest{Name: name})
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("failed to create tag: %v", err)), nil
			}
			data, _ := json.MarshalIndent(tag, "", "  ")
			return mcp.NewToolResultText(string(data)), nil
		},
	)

	// get_tag
	s.AddTool(
		mcp.NewTool("get_tag",
			mcp.WithDescription("Get a Clockify tag by ID."),
			mcp.WithString("workspace_id", mcp.Required(), mcp.Description("Workspace ID.")),
			mcp.WithString("tag_id", mcp.Required(), mcp.Description("Tag ID.")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			workspaceID, err := req.RequireString("workspace_id")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			tagID, err := req.RequireString("tag_id")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			tag, err := client.GetTag(workspaceID, tagID)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("failed to get tag: %v", err)), nil
			}
			data, _ := json.MarshalIndent(tag, "", "  ")
			return mcp.NewToolResultText(string(data)), nil
		},
	)

	// update_tag
	s.AddTool(
		mcp.NewTool("update_tag",
			mcp.WithDescription("Update an existing Clockify tag."),
			mcp.WithString("workspace_id", mcp.Required(), mcp.Description("Workspace ID.")),
			mcp.WithString("tag_id", mcp.Required(), mcp.Description("Tag ID.")),
			mcp.WithString("name", mcp.Required(), mcp.Description("New tag name.")),
			mcp.WithBoolean("archived", mcp.Description("Whether the tag is archived.")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			workspaceID, err := req.RequireString("workspace_id")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			tagID, err := req.RequireString("tag_id")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			name, err := req.RequireString("name")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			r := clockify.TagRequest{Name: name}
			if args := req.GetArguments(); args != nil {
				if v, ok := args["archived"].(bool); ok {
					r.Archived = &v
				}
			}
			tag, err := client.UpdateTag(workspaceID, tagID, r)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("failed to update tag: %v", err)), nil
			}
			data, _ := json.MarshalIndent(tag, "", "  ")
			return mcp.NewToolResultText(string(data)), nil
		},
	)

	// delete_tag
	s.AddTool(
		mcp.NewTool("delete_tag",
			mcp.WithDescription("Delete a tag from a Clockify workspace."),
			mcp.WithString("workspace_id", mcp.Required(), mcp.Description("Workspace ID.")),
			mcp.WithString("tag_id", mcp.Required(), mcp.Description("Tag ID.")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			workspaceID, err := req.RequireString("workspace_id")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			tagID, err := req.RequireString("tag_id")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			if err := client.DeleteTag(workspaceID, tagID); err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("failed to delete tag: %v", err)), nil
			}
			return mcp.NewToolResultText("Tag deleted successfully."), nil
		},
	)
}
