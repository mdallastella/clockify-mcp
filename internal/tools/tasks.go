package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"clockify-mcp/internal/clockify"
)

// RegisterTaskTools registers all task-related MCP tools.
func RegisterTaskTools(s *server.MCPServer, client *clockify.Client) {

	// list_tasks
	s.AddTool(
		mcp.NewTool("list_tasks",
			mcp.WithDescription("Get all tasks for a project in a Clockify workspace."),
			mcp.WithString("workspace_id", mcp.Required(), mcp.Description("Workspace ID.")),
			mcp.WithString("project_id", mcp.Required(), mcp.Description("Project ID.")),
			mcp.WithNumber("page", mcp.Description("Page number (1-based).")),
			mcp.WithNumber("page_size", mcp.Description("Number of tasks per page.")),
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
			page := req.GetInt("page", 1)
			pageSize := req.GetInt("page_size", 50)
			tasks, err := client.ListTasks(workspaceID, projectID, page, pageSize)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("failed to list tasks: %v", err)), nil
			}
			data, _ := json.MarshalIndent(tasks, "", "  ")
			return mcp.NewToolResultText(string(data)), nil
		},
	)

	// create_task
	s.AddTool(
		mcp.NewTool("create_task",
			mcp.WithDescription("Create a new task on a project in a Clockify workspace."),
			mcp.WithString("workspace_id", mcp.Required(), mcp.Description("Workspace ID.")),
			mcp.WithString("project_id", mcp.Required(), mcp.Description("Project ID.")),
			mcp.WithString("name", mcp.Required(), mcp.Description("Task name.")),
			mcp.WithString("estimate", mcp.Description("Time estimate in ISO 8601 duration format (e.g. PT1H30M).")),
			mcp.WithBoolean("billable", mcp.Description("Whether the task is billable.")),
			mcp.WithString("status", mcp.Description("Task status: ACTIVE or DONE.")),
			mcp.WithString("assignee_ids", mcp.Description(`JSON array of user IDs to assign, e.g. ["id1","id2"].`)),
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
			r := clockify.TaskRequest{
				Name:     name,
				Estimate: req.GetString("estimate", ""),
				Status:   req.GetString("status", ""),
			}
			if args := req.GetArguments(); args != nil {
				if v, ok := args["billable"].(bool); ok {
					r.Billable = &v
				}
			}
			if raw := req.GetString("assignee_ids", ""); raw != "" {
				if err := json.Unmarshal([]byte(raw), &r.AssigneeIDs); err != nil {
					return mcp.NewToolResultError(fmt.Sprintf("assignee_ids must be a JSON array: %v", err)), nil
				}
			}
			task, err := client.CreateTask(workspaceID, projectID, r)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("failed to create task: %v", err)), nil
			}
			data, _ := json.MarshalIndent(task, "", "  ")
			return mcp.NewToolResultText(string(data)), nil
		},
	)

	// get_task
	s.AddTool(
		mcp.NewTool("get_task",
			mcp.WithDescription("Get a specific task by ID from a project."),
			mcp.WithString("workspace_id", mcp.Required(), mcp.Description("Workspace ID.")),
			mcp.WithString("project_id", mcp.Required(), mcp.Description("Project ID.")),
			mcp.WithString("task_id", mcp.Required(), mcp.Description("Task ID.")),
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
			taskID, err := req.RequireString("task_id")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			task, err := client.GetTask(workspaceID, projectID, taskID)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("failed to get task: %v", err)), nil
			}
			data, _ := json.MarshalIndent(task, "", "  ")
			return mcp.NewToolResultText(string(data)), nil
		},
	)

	// update_task
	s.AddTool(
		mcp.NewTool("update_task",
			mcp.WithDescription("Update an existing task on a project."),
			mcp.WithString("workspace_id", mcp.Required(), mcp.Description("Workspace ID.")),
			mcp.WithString("project_id", mcp.Required(), mcp.Description("Project ID.")),
			mcp.WithString("task_id", mcp.Required(), mcp.Description("Task ID.")),
			mcp.WithString("name", mcp.Required(), mcp.Description("Task name.")),
			mcp.WithString("estimate", mcp.Description("Time estimate in ISO 8601 duration format (e.g. PT1H30M).")),
			mcp.WithBoolean("billable", mcp.Description("Whether the task is billable.")),
			mcp.WithString("status", mcp.Description("Task status: ACTIVE or DONE.")),
			mcp.WithString("assignee_ids", mcp.Description(`JSON array of user IDs to assign, e.g. ["id1","id2"].`)),
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
			taskID, err := req.RequireString("task_id")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			name, err := req.RequireString("name")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			r := clockify.TaskRequest{
				Name:     name,
				Estimate: req.GetString("estimate", ""),
				Status:   req.GetString("status", ""),
			}
			if args := req.GetArguments(); args != nil {
				if v, ok := args["billable"].(bool); ok {
					r.Billable = &v
				}
			}
			if raw := req.GetString("assignee_ids", ""); raw != "" {
				if err := json.Unmarshal([]byte(raw), &r.AssigneeIDs); err != nil {
					return mcp.NewToolResultError(fmt.Sprintf("assignee_ids must be a JSON array: %v", err)), nil
				}
			}
			task, err := client.UpdateTask(workspaceID, projectID, taskID, r)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("failed to update task: %v", err)), nil
			}
			data, _ := json.MarshalIndent(task, "", "  ")
			return mcp.NewToolResultText(string(data)), nil
		},
	)

	// delete_task
	s.AddTool(
		mcp.NewTool("delete_task",
			mcp.WithDescription("Delete a task from a project."),
			mcp.WithString("workspace_id", mcp.Required(), mcp.Description("Workspace ID.")),
			mcp.WithString("project_id", mcp.Required(), mcp.Description("Project ID.")),
			mcp.WithString("task_id", mcp.Required(), mcp.Description("Task ID.")),
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
			taskID, err := req.RequireString("task_id")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			if err := client.DeleteTask(workspaceID, projectID, taskID); err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("failed to delete task: %v", err)), nil
			}
			return mcp.NewToolResultText("Task deleted successfully."), nil
		},
	)

	// set_task_hourly_rate
	s.AddTool(
		mcp.NewTool("set_task_hourly_rate",
			mcp.WithDescription("Set the hourly rate for a task."),
			mcp.WithString("workspace_id", mcp.Required(), mcp.Description("Workspace ID.")),
			mcp.WithString("project_id", mcp.Required(), mcp.Description("Project ID.")),
			mcp.WithString("task_id", mcp.Required(), mcp.Description("Task ID.")),
			mcp.WithNumber("amount", mcp.Required(), mcp.Description("Rate amount in smallest currency unit.")),
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
			taskID, err := req.RequireString("task_id")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			amount, err := req.RequireInt("amount")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			r := clockify.RateRequest{Amount: amount, Currency: req.GetString("currency", "")}
			if err := client.SetTaskHourlyRate(workspaceID, projectID, taskID, r); err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("failed to set task hourly rate: %v", err)), nil
			}
			return mcp.NewToolResultText("Task hourly rate updated successfully."), nil
		},
	)

	// set_task_cost_rate
	s.AddTool(
		mcp.NewTool("set_task_cost_rate",
			mcp.WithDescription("Set the cost rate for a task."),
			mcp.WithString("workspace_id", mcp.Required(), mcp.Description("Workspace ID.")),
			mcp.WithString("project_id", mcp.Required(), mcp.Description("Project ID.")),
			mcp.WithString("task_id", mcp.Required(), mcp.Description("Task ID.")),
			mcp.WithNumber("amount", mcp.Required(), mcp.Description("Rate amount in smallest currency unit.")),
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
			taskID, err := req.RequireString("task_id")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			amount, err := req.RequireInt("amount")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			r := clockify.RateRequest{Amount: amount, Currency: req.GetString("currency", "")}
			if err := client.SetTaskCostRate(workspaceID, projectID, taskID, r); err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("failed to set task cost rate: %v", err)), nil
			}
			return mcp.NewToolResultText("Task cost rate updated successfully."), nil
		},
	)
}
