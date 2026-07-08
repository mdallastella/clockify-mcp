package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"clockify-mcp/internal/clockify"
)

// RegisterProjectTools registers all project-related MCP tools.
func RegisterProjectTools(s *server.MCPServer, client *clockify.Client) {

	// list_projects
	s.AddTool(
		mcp.NewTool("list_projects",
			mcp.WithDescription("Get all projects in a Clockify workspace."),
			mcp.WithString("workspace_id", mcp.Required(), mcp.Description("Workspace ID.")),
			mcp.WithNumber("page", mcp.Description("Page number (1-based).")),
			mcp.WithNumber("page_size", mcp.Description("Number of projects per page (max 5000).")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			workspaceID, err := req.RequireString("workspace_id")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			page := req.GetInt("page", 1)
			pageSize := req.GetInt("page_size", 50)
			projects, err := client.ListProjects(workspaceID, page, pageSize)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("failed to list projects: %v", err)), nil
			}
			data, _ := json.MarshalIndent(projects, "", "  ")
			return mcp.NewToolResultText(string(data)), nil
		},
	)

	// create_project
	s.AddTool(
		mcp.NewTool("create_project",
			mcp.WithDescription("Create a new project in a Clockify workspace."),
			mcp.WithString("workspace_id", mcp.Required(), mcp.Description("Workspace ID.")),
			mcp.WithString("name", mcp.Required(), mcp.Description("Project name (2–250 characters).")),
			mcp.WithString("client_id", mcp.Description("Client ID to associate with this project.")),
			mcp.WithString("color", mcp.Description("Hex color code (e.g. #FF5733).")),
			mcp.WithBoolean("billable", mcp.Description("Whether the project is billable.")),
			mcp.WithBoolean("is_public", mcp.Description("Whether the project is public.")),
			mcp.WithString("note", mcp.Description("Project notes (max 16384 characters).")),
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
			r := clockify.ProjectRequest{
				Name:     name,
				ClientID: req.GetString("client_id", ""),
				Color:    req.GetString("color", ""),
				Note:     req.GetString("note", ""),
			}
			if args := req.GetArguments(); args != nil {
				if v, ok := args["billable"].(bool); ok {
					r.Billable = &v
				}
				if v, ok := args["is_public"].(bool); ok {
					r.IsPublic = &v
				}
			}
			project, err := client.CreateProject(workspaceID, r)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("failed to create project: %v", err)), nil
			}
			data, _ := json.MarshalIndent(project, "", "  ")
			return mcp.NewToolResultText(string(data)), nil
		},
	)

	// get_project
	s.AddTool(
		mcp.NewTool("get_project",
			mcp.WithDescription("Get a Clockify project by ID."),
			mcp.WithString("workspace_id", mcp.Required(), mcp.Description("Workspace ID.")),
			mcp.WithString("project_id", mcp.Required(), mcp.Description("Project ID.")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			workspaceID, err := req.RequireString("workspace_id")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			projectID, err := req.RequireString("project_id")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			project, err := client.GetProject(workspaceID, projectID)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("failed to get project: %v", err)), nil
			}
			data, _ := json.MarshalIndent(project, "", "  ")
			return mcp.NewToolResultText(string(data)), nil
		},
	)

	// update_project
	s.AddTool(
		mcp.NewTool("update_project",
			mcp.WithDescription("Update an existing Clockify project."),
			mcp.WithString("workspace_id", mcp.Required(), mcp.Description("Workspace ID.")),
			mcp.WithString("project_id", mcp.Required(), mcp.Description("Project ID.")),
			mcp.WithString("name", mcp.Required(), mcp.Description("Project name (2–250 characters).")),
			mcp.WithString("client_id", mcp.Description("Client ID.")),
			mcp.WithString("color", mcp.Description("Hex color code.")),
			mcp.WithBoolean("billable", mcp.Description("Whether the project is billable.")),
			mcp.WithBoolean("is_public", mcp.Description("Whether the project is public.")),
			mcp.WithString("note", mcp.Description("Project notes.")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			workspaceID, err := req.RequireString("workspace_id")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			projectID, err := req.RequireString("project_id")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			name, err := req.RequireString("name")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			r := clockify.ProjectRequest{
				Name:     name,
				ClientID: req.GetString("client_id", ""),
				Color:    req.GetString("color", ""),
				Note:     req.GetString("note", ""),
			}
			if args := req.GetArguments(); args != nil {
				if v, ok := args["billable"].(bool); ok {
					r.Billable = &v
				}
				if v, ok := args["is_public"].(bool); ok {
					r.IsPublic = &v
				}
			}
			project, err := client.UpdateProject(workspaceID, projectID, r)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("failed to update project: %v", err)), nil
			}
			data, _ := json.MarshalIndent(project, "", "  ")
			return mcp.NewToolResultText(string(data)), nil
		},
	)

	// delete_project
	s.AddTool(
		mcp.NewTool("delete_project",
			mcp.WithDescription("Delete a project from a Clockify workspace."),
			mcp.WithString("workspace_id", mcp.Required(), mcp.Description("Workspace ID.")),
			mcp.WithString("project_id", mcp.Required(), mcp.Description("Project ID.")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			workspaceID, err := req.RequireString("workspace_id")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			projectID, err := req.RequireString("project_id")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			if err := client.DeleteProject(workspaceID, projectID); err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("failed to delete project: %v", err)), nil
			}
			return mcp.NewToolResultText("Project deleted successfully."), nil
		},
	)

	// update_project_estimate
	s.AddTool(
		mcp.NewTool("update_project_estimate",
			mcp.WithDescription("Update the time estimate for a project."),
			mcp.WithString("workspace_id", mcp.Required(), mcp.Description("Workspace ID.")),
			mcp.WithString("project_id", mcp.Required(), mcp.Description("Project ID.")),
			mcp.WithString("estimate", mcp.Required(), mcp.Description("Estimate duration in ISO 8601 format (e.g. PT8H).")),
			mcp.WithString("type", mcp.Description("Estimate type: AUTO or MANUAL.")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			workspaceID, err := req.RequireString("workspace_id")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			projectID, err := req.RequireString("project_id")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			estimate, err := req.RequireString("estimate")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			r := clockify.EstimateRequest{
				Estimate: estimate,
				Type:     req.GetString("type", ""),
			}
			project, err := client.UpdateProjectEstimate(workspaceID, projectID, r)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("failed to update project estimate: %v", err)), nil
			}
			data, _ := json.MarshalIndent(project, "", "  ")
			return mcp.NewToolResultText(string(data)), nil
		},
	)

	// add_users_to_project
	s.AddTool(
		mcp.NewTool("add_users_to_project",
			mcp.WithDescription("Assign or remove users from a project. Provide user IDs as a JSON array string."),
			mcp.WithString("workspace_id", mcp.Required(), mcp.Description("Workspace ID.")),
			mcp.WithString("project_id", mcp.Required(), mcp.Description("Project ID.")),
			mcp.WithString("user_ids", mcp.Description(`JSON array of user IDs to add, e.g. ["id1","id2"].`)),
			mcp.WithString("group_ids", mcp.Description(`JSON array of group IDs to add, e.g. ["id1","id2"].`)),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			workspaceID, err := req.RequireString("workspace_id")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			projectID, err := req.RequireString("project_id")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			r := clockify.AddUsersToProjectRequest{}
			if raw := req.GetString("user_ids", ""); raw != "" {
				if err := json.Unmarshal([]byte(raw), &r.UserIDs); err != nil {
					return mcp.NewToolResultError(fmt.Sprintf("user_ids must be a JSON array: %v", err)), nil
				}
			}
			if raw := req.GetString("group_ids", ""); raw != "" {
				if err := json.Unmarshal([]byte(raw), &r.GroupIDs); err != nil {
					return mcp.NewToolResultError(fmt.Sprintf("group_ids must be a JSON array: %v", err)), nil
				}
			}
			project, err := client.AddUsersToProject(workspaceID, projectID, r)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("failed to add users to project: %v", err)), nil
			}
			data, _ := json.MarshalIndent(project, "", "  ")
			return mcp.NewToolResultText(string(data)), nil
		},
	)

	// create_project_from_template
	s.AddTool(
		mcp.NewTool("create_project_from_template",
			mcp.WithDescription("Create a new project from an existing project template."),
			mcp.WithString("workspace_id", mcp.Required(), mcp.Description("Workspace ID.")),
			mcp.WithString("name", mcp.Required(), mcp.Description("New project name.")),
			mcp.WithString("template_id", mcp.Required(), mcp.Description("Template project ID.")),
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
			templateID, err := req.RequireString("template_id")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			r := clockify.ProjectTemplateRequest{Name: name, TemplateID: templateID}
			project, err := client.CreateProjectFromTemplate(workspaceID, r)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("failed to create project from template: %v", err)), nil
			}
			data, _ := json.MarshalIndent(project, "", "  ")
			return mcp.NewToolResultText(string(data)), nil
		},
	)

	// update_project_template_flag
	s.AddTool(
		mcp.NewTool("update_project_template_flag",
			mcp.WithDescription("Set or unset a project as a template."),
			mcp.WithString("workspace_id", mcp.Required(), mcp.Description("Workspace ID.")),
			mcp.WithString("project_id", mcp.Required(), mcp.Description("Project ID.")),
			mcp.WithBoolean("is_template", mcp.Required(), mcp.Description("Whether the project should be a template.")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			workspaceID, err := req.RequireString("workspace_id")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			projectID, err := req.RequireString("project_id")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			isTemplate, err := req.RequireBool("is_template")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			project, err := client.UpdateProjectTemplateFlag(workspaceID, projectID, isTemplate)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("failed to update project template flag: %v", err)), nil
			}
			data, _ := json.MarshalIndent(project, "", "  ")
			return mcp.NewToolResultText(string(data)), nil
		},
	)

	// set_project_user_hourly_rate
	s.AddTool(
		mcp.NewTool("set_project_user_hourly_rate",
			mcp.WithDescription("Set the hourly rate for a specific user on a project."),
			mcp.WithString("workspace_id", mcp.Required(), mcp.Description("Workspace ID.")),
			mcp.WithString("project_id", mcp.Required(), mcp.Description("Project ID.")),
			mcp.WithString("user_id", mcp.Required(), mcp.Description("User ID.")),
			mcp.WithNumber("amount", mcp.Required(), mcp.Description("Rate amount in smallest currency unit (e.g. 10500 = $105.00).")),
			mcp.WithString("currency", mcp.Description("ISO currency code (e.g. USD).")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			workspaceID, err := req.RequireString("workspace_id")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			projectID, err := req.RequireString("project_id")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			userID, err := req.RequireString("user_id")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			amount, err := req.RequireInt("amount")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			r := clockify.RateRequest{Amount: amount, Currency: req.GetString("currency", "")}
			if err := client.SetProjectUserHourlyRate(workspaceID, projectID, userID, r); err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("failed to set hourly rate: %v", err)), nil
			}
			return mcp.NewToolResultText("Hourly rate updated successfully."), nil
		},
	)

	// set_project_user_cost_rate
	s.AddTool(
		mcp.NewTool("set_project_user_cost_rate",
			mcp.WithDescription("Set the cost rate for a specific user on a project."),
			mcp.WithString("workspace_id", mcp.Required(), mcp.Description("Workspace ID.")),
			mcp.WithString("project_id", mcp.Required(), mcp.Description("Project ID.")),
			mcp.WithString("user_id", mcp.Required(), mcp.Description("User ID.")),
			mcp.WithNumber("amount", mcp.Required(), mcp.Description("Rate amount in smallest currency unit (e.g. 10500 = $105.00).")),
			mcp.WithString("currency", mcp.Description("ISO currency code (e.g. USD).")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			workspaceID, err := req.RequireString("workspace_id")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			projectID, err := req.RequireString("project_id")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			userID, err := req.RequireString("user_id")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			amount, err := req.RequireInt("amount")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			r := clockify.RateRequest{Amount: amount, Currency: req.GetString("currency", "")}
			if err := client.SetProjectUserCostRate(workspaceID, projectID, userID, r); err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("failed to set cost rate: %v", err)), nil
			}
			return mcp.NewToolResultText("Cost rate updated successfully."), nil
		},
	)
}
