package clockify

import (
	"fmt"
	"net/url"
	"strconv"
)

// TimeInterval represents a time entry's interval.
type TimeInterval struct {
	Start    string `json:"start"`
	End      string `json:"end,omitempty"`
	Duration string `json:"duration,omitempty"`
}

// TimeEntryDTO represents a Clockify time entry.
type TimeEntryDTO struct {
	ID           string       `json:"id"`
	WorkspaceID  string       `json:"workspaceId"`
	UserID       string       `json:"userId"`
	ProjectID    string       `json:"projectId,omitempty"`
	TaskID       string       `json:"taskId,omitempty"`
	TagIDs       []string     `json:"tagIds,omitempty"`
	Billable     bool         `json:"billable"`
	Description  string       `json:"description,omitempty"`
	IsLocked     bool         `json:"isLocked"`
	Type         string       `json:"type,omitempty"`
	TimeInterval TimeInterval `json:"timeInterval"`
}

// CreateTimeEntryRequest is the payload for creating a time entry.
type CreateTimeEntryRequest struct {
	Start       string   `json:"start,omitempty"`
	End         string   `json:"end,omitempty"`
	Description string   `json:"description,omitempty"`
	ProjectID   string   `json:"projectId,omitempty"`
	TaskID      string   `json:"taskId,omitempty"`
	TagIDs      []string `json:"tagIds,omitempty"`
	Billable    *bool    `json:"billable,omitempty"`
	Type        string   `json:"type,omitempty"`
}

// UpdateTimeEntryRequest is the payload for updating a time entry.
type UpdateTimeEntryRequest struct {
	Start       string   `json:"start"`
	End         string   `json:"end,omitempty"`
	Description string   `json:"description,omitempty"`
	ProjectID   string   `json:"projectId,omitempty"`
	TaskID      string   `json:"taskId,omitempty"`
	TagIDs      []string `json:"tagIds,omitempty"`
	Billable    *bool    `json:"billable,omitempty"`
	Type        string   `json:"type,omitempty"`
}

// UpdateInvoicedRequest marks time entries as invoiced.
type UpdateInvoicedRequest struct {
	TimeEntryIDs []string `json:"timeEntryIds"`
	Invoiced     bool     `json:"invoiced"`
}

// BulkUpdateTimeEntriesRequest replaces multiple time entries.
type BulkUpdateTimeEntriesRequest struct {
	TimeEntries []UpdateTimeEntryRequest `json:"timeEntries"`
}

// CreateTimeEntry creates a new time entry in a workspace.
func (c *Client) CreateTimeEntry(workspaceID string, req CreateTimeEntryRequest) (*TimeEntryDTO, error) {
	var result TimeEntryDTO
	if err := c.post(fmt.Sprintf("/workspaces/%s/time-entries", workspaceID), req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetTimeEntry returns a single time entry by ID.
func (c *Client) GetTimeEntry(workspaceID, timeEntryID string) (*TimeEntryDTO, error) {
	var result TimeEntryDTO
	if err := c.get(fmt.Sprintf("/workspaces/%s/time-entries/%s", workspaceID, timeEntryID), nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdateTimeEntry updates an existing time entry.
func (c *Client) UpdateTimeEntry(workspaceID, timeEntryID string, req UpdateTimeEntryRequest) (*TimeEntryDTO, error) {
	var result TimeEntryDTO
	if err := c.put(fmt.Sprintf("/workspaces/%s/time-entries/%s", workspaceID, timeEntryID), req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// DeleteTimeEntry deletes a time entry by ID.
func (c *Client) DeleteTimeEntry(workspaceID, timeEntryID string) error {
	return c.delete(fmt.Sprintf("/workspaces/%s/time-entries/%s", workspaceID, timeEntryID), nil)
}

// GetInProgressTimeEntries returns all currently running time entries in a workspace.
func (c *Client) GetInProgressTimeEntries(workspaceID string) ([]TimeEntryDTO, error) {
	var result []TimeEntryDTO
	if err := c.get(fmt.Sprintf("/workspaces/%s/time-entries/status/in-progress", workspaceID), nil, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// ListTimeEntries returns time entries for a specific user in a workspace.
func (c *Client) ListTimeEntries(workspaceID, userID string, page, pageSize int) ([]TimeEntryDTO, error) {
	q := url.Values{}
	if page > 0 {
		q.Set("page", strconv.Itoa(page))
	}
	if pageSize > 0 {
		q.Set("page-size", strconv.Itoa(pageSize))
	}
	var result []TimeEntryDTO
	if err := c.get(fmt.Sprintf("/workspaces/%s/user/%s/time-entries", workspaceID, userID), q, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// StopTimer stops the currently running timer for a user.
func (c *Client) StopTimer(workspaceID, userID string, end string) (*TimeEntryDTO, error) {
	body := map[string]string{"end": end}
	var result TimeEntryDTO
	if err := c.patch(fmt.Sprintf("/workspaces/%s/user/%s/time-entries", workspaceID, userID), body, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// CreateTimeEntryForUser creates a time entry on behalf of another user.
func (c *Client) CreateTimeEntryForUser(workspaceID, userID string, req CreateTimeEntryRequest) (*TimeEntryDTO, error) {
	var result TimeEntryDTO
	if err := c.post(fmt.Sprintf("/workspaces/%s/user/%s/time-entries", workspaceID, userID), req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// DuplicateTimeEntry duplicates an existing time entry.
func (c *Client) DuplicateTimeEntry(workspaceID, userID, timeEntryID string) (*TimeEntryDTO, error) {
	var result TimeEntryDTO
	if err := c.post(fmt.Sprintf("/workspaces/%s/user/%s/time-entries/%s/duplicate", workspaceID, userID, timeEntryID), nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// DeleteUserTimeEntries deletes all time entries for a user in a workspace.
func (c *Client) DeleteUserTimeEntries(workspaceID, userID string) error {
	return c.delete(fmt.Sprintf("/workspaces/%s/user/%s/time-entries", workspaceID, userID), nil)
}

// MarkTimeEntriesInvoiced marks time entries as invoiced or not invoiced.
func (c *Client) MarkTimeEntriesInvoiced(workspaceID string, req UpdateInvoicedRequest) error {
	return c.patch(fmt.Sprintf("/workspaces/%s/time-entries/invoiced", workspaceID), req, nil)
}
