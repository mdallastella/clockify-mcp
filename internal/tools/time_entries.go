package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"clockify-mcp/internal/clockify"
)

// RegisterTimeEntryTools registers all time-entry-related MCP tools.
func RegisterTimeEntryTools(s *server.MCPServer, client *clockify.Client) {

	// create_time_entry
	s.AddTool(
		mcp.NewTool("create_time_entry",
			mcp.WithDescription("Create a new time entry in a workspace. Omit 'end' to start a running timer."),
			mcp.WithString("workspace_id", mcp.Required(), mcp.Description("Workspace ID.")),
			mcp.WithString("start", mcp.Description("Start time in ISO 8601 format (e.g. 2024-01-15T09:00:00Z). Omit to start a timer now.")),
			mcp.WithString("end", mcp.Description("End time in ISO 8601 format. Omit to keep the timer running.")),
			mcp.WithString("description", mcp.Description("Description of the time entry.")),
			mcp.WithString("project_id", mcp.Description("Project ID.")),
			mcp.WithString("task_id", mcp.Description("Task ID.")),
			mcp.WithBoolean("billable", mcp.Description("Whether the time entry is billable.")),
			mcp.WithString("type", mcp.Description("Type: REGULAR or BREAK.")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			workspaceID, err := req.RequireString("workspace_id")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			r := clockify.CreateTimeEntryRequest{
				Start:       req.GetString("start", ""),
				End:         req.GetString("end", ""),
				Description: req.GetString("description", ""),
				ProjectID:   req.GetString("project_id", ""),
				TaskID:      req.GetString("task_id", ""),
				Type:        req.GetString("type", ""),
			}
			if args := req.GetArguments(); args != nil {
				if v, ok := args["billable"].(bool); ok {
					r.Billable = &v
				}
			}
			entry, err := client.CreateTimeEntry(workspaceID, r)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("failed to create time entry: %v", err)), nil
			}
			data, _ := json.MarshalIndent(entry, "", "  ")
			return mcp.NewToolResultText(string(data)), nil
		},
	)

	// get_time_entry
	s.AddTool(
		mcp.NewTool("get_time_entry",
			mcp.WithDescription("Get a specific time entry from a workspace by ID."),
			mcp.WithString("workspace_id", mcp.Required(), mcp.Description("Workspace ID.")),
			mcp.WithString("time_entry_id", mcp.Required(), mcp.Description("Time entry ID.")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			workspaceID, err := req.RequireString("workspace_id")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			timeEntryID, err := req.RequireString("time_entry_id")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			entry, err := client.GetTimeEntry(workspaceID, timeEntryID)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("failed to get time entry: %v", err)), nil
			}
			data, _ := json.MarshalIndent(entry, "", "  ")
			return mcp.NewToolResultText(string(data)), nil
		},
	)

	// update_time_entry
	s.AddTool(
		mcp.NewTool("update_time_entry",
			mcp.WithDescription("Update an existing time entry. The 'start' field is required."),
			mcp.WithString("workspace_id", mcp.Required(), mcp.Description("Workspace ID.")),
			mcp.WithString("time_entry_id", mcp.Required(), mcp.Description("Time entry ID.")),
			mcp.WithString("start", mcp.Required(), mcp.Description("Start time in ISO 8601 format.")),
			mcp.WithString("end", mcp.Description("End time in ISO 8601 format.")),
			mcp.WithString("description", mcp.Description("Description of the time entry.")),
			mcp.WithString("project_id", mcp.Description("Project ID.")),
			mcp.WithString("task_id", mcp.Description("Task ID.")),
			mcp.WithBoolean("billable", mcp.Description("Whether the time entry is billable.")),
			mcp.WithString("type", mcp.Description("Type: REGULAR or BREAK.")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			workspaceID, err := req.RequireString("workspace_id")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			timeEntryID, err := req.RequireString("time_entry_id")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			start, err := req.RequireString("start")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			r := clockify.UpdateTimeEntryRequest{
				Start:       start,
				End:         req.GetString("end", ""),
				Description: req.GetString("description", ""),
				ProjectID:   req.GetString("project_id", ""),
				TaskID:      req.GetString("task_id", ""),
				Type:        req.GetString("type", ""),
			}
			if args := req.GetArguments(); args != nil {
				if v, ok := args["billable"].(bool); ok {
					r.Billable = &v
				}
			}
			entry, err := client.UpdateTimeEntry(workspaceID, timeEntryID, r)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("failed to update time entry: %v", err)), nil
			}
			data, _ := json.MarshalIndent(entry, "", "  ")
			return mcp.NewToolResultText(string(data)), nil
		},
	)

	// delete_time_entry
	s.AddTool(
		mcp.NewTool("delete_time_entry",
			mcp.WithDescription("Delete a time entry from a workspace."),
			mcp.WithString("workspace_id", mcp.Required(), mcp.Description("Workspace ID.")),
			mcp.WithString("time_entry_id", mcp.Required(), mcp.Description("Time entry ID.")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			workspaceID, err := req.RequireString("workspace_id")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			timeEntryID, err := req.RequireString("time_entry_id")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			if err := client.DeleteTimeEntry(workspaceID, timeEntryID); err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("failed to delete time entry: %v", err)), nil
			}
			return mcp.NewToolResultText("Time entry deleted successfully."), nil
		},
	)

	// get_in_progress_time_entries
	s.AddTool(
		mcp.NewTool("get_in_progress_time_entries",
			mcp.WithDescription("Get all currently running (in-progress) time entries on a workspace."),
			mcp.WithString("workspace_id", mcp.Required(), mcp.Description("Workspace ID.")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			workspaceID, err := req.RequireString("workspace_id")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			entries, err := client.GetInProgressTimeEntries(workspaceID)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("failed to get in-progress entries: %v", err)), nil
			}
			data, _ := json.MarshalIndent(entries, "", "  ")
			return mcp.NewToolResultText(string(data)), nil
		},
	)

	// list_time_entries
	s.AddTool(
		mcp.NewTool("list_time_entries",
			mcp.WithDescription("Get time entries for a specific user in a workspace."),
			mcp.WithString("workspace_id", mcp.Required(), mcp.Description("Workspace ID.")),
			mcp.WithString("user_id", mcp.Required(), mcp.Description("User ID.")),
			mcp.WithNumber("page", mcp.Description("Page number (1-based).")),
			mcp.WithNumber("page_size", mcp.Description("Number of entries per page (max 200).")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			workspaceID, err := req.RequireString("workspace_id")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			userID, err := req.RequireString("user_id")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			page := req.GetInt("page", 1)
			pageSize := req.GetInt("page_size", 50)
			entries, err := client.ListTimeEntries(workspaceID, userID, page, pageSize)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("failed to list time entries: %v", err)), nil
			}
			data, _ := json.MarshalIndent(entries, "", "  ")
			return mcp.NewToolResultText(string(data)), nil
		},
	)

	// stop_timer
	s.AddTool(
		mcp.NewTool("stop_timer",
			mcp.WithDescription("Stop the currently running timer for a user."),
			mcp.WithString("workspace_id", mcp.Required(), mcp.Description("Workspace ID.")),
			mcp.WithString("user_id", mcp.Required(), mcp.Description("User ID.")),
			mcp.WithString("end", mcp.Required(), mcp.Description("End time in ISO 8601 format (e.g. 2024-01-15T17:00:00Z).")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			workspaceID, err := req.RequireString("workspace_id")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			userID, err := req.RequireString("user_id")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			end, err := req.RequireString("end")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			entry, err := client.StopTimer(workspaceID, userID, end)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("failed to stop timer: %v", err)), nil
			}
			data, _ := json.MarshalIndent(entry, "", "  ")
			return mcp.NewToolResultText(string(data)), nil
		},
	)

	// create_time_entry_for_user
	s.AddTool(
		mcp.NewTool("create_time_entry_for_user",
			mcp.WithDescription("Create a time entry on behalf of another user in a workspace."),
			mcp.WithString("workspace_id", mcp.Required(), mcp.Description("Workspace ID.")),
			mcp.WithString("user_id", mcp.Required(), mcp.Description("User ID of the target user.")),
			mcp.WithString("start", mcp.Description("Start time in ISO 8601 format.")),
			mcp.WithString("end", mcp.Description("End time in ISO 8601 format.")),
			mcp.WithString("description", mcp.Description("Description of the time entry.")),
			mcp.WithString("project_id", mcp.Description("Project ID.")),
			mcp.WithString("task_id", mcp.Description("Task ID.")),
			mcp.WithBoolean("billable", mcp.Description("Whether the time entry is billable.")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			workspaceID, err := req.RequireString("workspace_id")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			userID, err := req.RequireString("user_id")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			r := clockify.CreateTimeEntryRequest{
				Start:       req.GetString("start", ""),
				End:         req.GetString("end", ""),
				Description: req.GetString("description", ""),
				ProjectID:   req.GetString("project_id", ""),
				TaskID:      req.GetString("task_id", ""),
			}
			if args := req.GetArguments(); args != nil {
				if v, ok := args["billable"].(bool); ok {
					r.Billable = &v
				}
			}
			entry, err := client.CreateTimeEntryForUser(workspaceID, userID, r)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("failed to create time entry: %v", err)), nil
			}
			data, _ := json.MarshalIndent(entry, "", "  ")
			return mcp.NewToolResultText(string(data)), nil
		},
	)

	// duplicate_time_entry
	s.AddTool(
		mcp.NewTool("duplicate_time_entry",
			mcp.WithDescription("Duplicate an existing time entry."),
			mcp.WithString("workspace_id", mcp.Required(), mcp.Description("Workspace ID.")),
			mcp.WithString("user_id", mcp.Required(), mcp.Description("User ID.")),
			mcp.WithString("time_entry_id", mcp.Required(), mcp.Description("Time entry ID to duplicate.")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			workspaceID, err := req.RequireString("workspace_id")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			userID, err := req.RequireString("user_id")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			timeEntryID, err := req.RequireString("time_entry_id")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			entry, err := client.DuplicateTimeEntry(workspaceID, userID, timeEntryID)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("failed to duplicate time entry: %v", err)), nil
			}
			data, _ := json.MarshalIndent(entry, "", "  ")
			return mcp.NewToolResultText(string(data)), nil
		},
	)

	// delete_user_time_entries
	s.AddTool(
		mcp.NewTool("delete_user_time_entries",
			mcp.WithDescription("Delete ALL time entries for a user in a workspace. This action is irreversible."),
			mcp.WithString("workspace_id", mcp.Required(), mcp.Description("Workspace ID.")),
			mcp.WithString("user_id", mcp.Required(), mcp.Description("User ID whose time entries will be deleted.")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			workspaceID, err := req.RequireString("workspace_id")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			userID, err := req.RequireString("user_id")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			if err := client.DeleteUserTimeEntries(workspaceID, userID); err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("failed to delete user time entries: %v", err)), nil
			}
			return mcp.NewToolResultText("All user time entries deleted successfully."), nil
		},
	)

	// mark_time_entries_invoiced
	s.AddTool(
		mcp.NewTool("mark_time_entries_invoiced",
			mcp.WithDescription("Mark a set of time entries as invoiced or not invoiced."),
			mcp.WithString("workspace_id", mcp.Required(), mcp.Description("Workspace ID.")),
			mcp.WithBoolean("invoiced", mcp.Required(), mcp.Description("Set to true to mark as invoiced, false to un-mark.")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			workspaceID, err := req.RequireString("workspace_id")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			invoiced, err := req.RequireBool("invoiced")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			// time_entry_ids is passed as a JSON array string
			timeEntryIDsRaw := req.GetString("time_entry_ids", "")
			var timeEntryIDs []string
			if timeEntryIDsRaw != "" {
				if err := json.Unmarshal([]byte(timeEntryIDsRaw), &timeEntryIDs); err != nil {
					return mcp.NewToolResultError(fmt.Sprintf("time_entry_ids must be a JSON array of strings: %v", err)), nil
				}
			}
			r := clockify.UpdateInvoicedRequest{
				TimeEntryIDs: timeEntryIDs,
				Invoiced:     invoiced,
			}
			if err := client.MarkTimeEntriesInvoiced(workspaceID, r); err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("failed to mark time entries invoiced: %v", err)), nil
			}
			return mcp.NewToolResultText("Time entries invoiced status updated successfully."), nil
		},
	)
}
