package clockify

import (
	"fmt"
	"net/url"
	"strconv"
)

// TaskDTO represents a Clockify task.
type TaskDTO struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	ProjectID   string   `json:"projectId"`
	Status      string   `json:"status,omitempty"`
	AssigneeIDs []string `json:"assigneeIds,omitempty"`
	Billable    bool     `json:"billable"`
	Estimate    string   `json:"estimate,omitempty"`
	Duration    string   `json:"duration,omitempty"`
}

// TaskRequest is the payload for creating or updating a task.
type TaskRequest struct {
	Name        string   `json:"name"`
	AssigneeIDs []string `json:"assigneeIds,omitempty"`
	Estimate    string   `json:"estimate,omitempty"`
	Billable    *bool    `json:"billable,omitempty"`
	Status      string   `json:"status,omitempty"`
}

// ListTasks returns all tasks for a project in a workspace.
func (c *Client) ListTasks(workspaceID, projectID string, page, pageSize int) ([]TaskDTO, error) {
	q := url.Values{}
	if page > 0 {
		q.Set("page", strconv.Itoa(page))
	}
	if pageSize > 0 {
		q.Set("page-size", strconv.Itoa(pageSize))
	}
	var result []TaskDTO
	if err := c.get(fmt.Sprintf("/workspaces/%s/projects/%s/tasks", workspaceID, projectID), q, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// CreateTask creates a new task on a project.
func (c *Client) CreateTask(workspaceID, projectID string, req TaskRequest) (*TaskDTO, error) {
	var result TaskDTO
	if err := c.post(fmt.Sprintf("/workspaces/%s/projects/%s/tasks", workspaceID, projectID), req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetTask returns a task by ID.
func (c *Client) GetTask(workspaceID, projectID, taskID string) (*TaskDTO, error) {
	var result TaskDTO
	if err := c.get(fmt.Sprintf("/workspaces/%s/projects/%s/tasks/%s", workspaceID, projectID, taskID), nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdateTask updates a task on a project.
func (c *Client) UpdateTask(workspaceID, projectID, taskID string, req TaskRequest) (*TaskDTO, error) {
	var result TaskDTO
	if err := c.put(fmt.Sprintf("/workspaces/%s/projects/%s/tasks/%s", workspaceID, projectID, taskID), req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// DeleteTask deletes a task from a project.
func (c *Client) DeleteTask(workspaceID, projectID, taskID string) error {
	return c.delete(fmt.Sprintf("/workspaces/%s/projects/%s/tasks/%s", workspaceID, projectID, taskID), nil)
}

// SetTaskHourlyRate sets the hourly rate for a task.
func (c *Client) SetTaskHourlyRate(workspaceID, projectID, taskID string, req RateRequest) error {
	return c.put(fmt.Sprintf("/workspaces/%s/projects/%s/tasks/%s/hourly-rate", workspaceID, projectID, taskID), req, nil)
}

// SetTaskCostRate sets the cost rate for a task.
func (c *Client) SetTaskCostRate(workspaceID, projectID, taskID string, req RateRequest) error {
	return c.put(fmt.Sprintf("/workspaces/%s/projects/%s/tasks/%s/cost-rate", workspaceID, projectID, taskID), req, nil)
}
