package clockify

import (
	"fmt"
	"net/url"
	"strconv"
)

// TagDTO represents a Clockify tag.
type TagDTO struct {
	ID          string `json:"id"`
	WorkspaceID string `json:"workspaceId"`
	Name        string `json:"name"`
	Archived    bool   `json:"archived"`
}

// TagRequest is the payload for creating or updating a tag.
type TagRequest struct {
	Name     string `json:"name"`
	Archived *bool  `json:"archived,omitempty"`
}

// ListTags returns all tags in a workspace.
func (c *Client) ListTags(workspaceID string, page, pageSize int) ([]TagDTO, error) {
	q := url.Values{}
	if page > 0 {
		q.Set("page", strconv.Itoa(page))
	}
	if pageSize > 0 {
		q.Set("page-size", strconv.Itoa(pageSize))
	}
	var result []TagDTO
	if err := c.get(fmt.Sprintf("/workspaces/%s/tags", workspaceID), q, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// CreateTag creates a new tag in a workspace.
func (c *Client) CreateTag(workspaceID string, req TagRequest) (*TagDTO, error) {
	var result TagDTO
	if err := c.post(fmt.Sprintf("/workspaces/%s/tags", workspaceID), req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetTag returns a tag by ID.
func (c *Client) GetTag(workspaceID, tagID string) (*TagDTO, error) {
	var result TagDTO
	if err := c.get(fmt.Sprintf("/workspaces/%s/tags/%s", workspaceID, tagID), nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdateTag updates an existing tag.
func (c *Client) UpdateTag(workspaceID, tagID string, req TagRequest) (*TagDTO, error) {
	var result TagDTO
	if err := c.put(fmt.Sprintf("/workspaces/%s/tags/%s", workspaceID, tagID), req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// DeleteTag deletes a tag by ID.
func (c *Client) DeleteTag(workspaceID, tagID string) error {
	return c.delete(fmt.Sprintf("/workspaces/%s/tags/%s", workspaceID, tagID), nil)
}
