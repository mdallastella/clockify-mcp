package clockify

import "fmt"

// WorkspaceDTO represents a Clockify workspace.
type WorkspaceDTO struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	// Additional fields returned by the API
	ImageURL string `json:"imageUrl,omitempty"`
}

// ListWorkspaces returns all workspaces for the authenticated user.
func (c *Client) ListWorkspaces() ([]WorkspaceDTO, error) {
	var result []WorkspaceDTO
	if err := c.get("/workspaces", nil, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// GetWorkspace returns a single workspace by ID.
func (c *Client) GetWorkspace(workspaceID string) (*WorkspaceDTO, error) {
	var result WorkspaceDTO
	if err := c.get(fmt.Sprintf("/workspaces/%s", workspaceID), nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
