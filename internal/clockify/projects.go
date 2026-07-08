package clockify

import (
	"fmt"
	"net/url"
	"strconv"
)

// ProjectDTO represents a Clockify project.
type ProjectDTO struct {
	ID          string `json:"id"`
	WorkspaceID string `json:"workspaceId"`
	Name        string `json:"name"`
	ClientID    string `json:"clientId,omitempty"`
	ClientName  string `json:"clientName,omitempty"`
	Color       string `json:"color,omitempty"`
	Billable    bool   `json:"billable"`
	Archived    bool   `json:"archived"`
	Public      bool   `json:"public"`
	Template    bool   `json:"template"`
	Note        string `json:"note,omitempty"`
	Duration    string `json:"duration,omitempty"`
}

// ProjectRequest is the payload for creating or updating a project.
type ProjectRequest struct {
	Name     string `json:"name"`
	ClientID string `json:"clientId,omitempty"`
	Color    string `json:"color,omitempty"`
	Billable *bool  `json:"billable,omitempty"`
	IsPublic *bool  `json:"isPublic,omitempty"`
	Note     string `json:"note,omitempty"`
}

// EstimateRequest is used to update a project's estimate.
type EstimateRequest struct {
	Estimate string `json:"estimate"`
	Type     string `json:"type,omitempty"`
}

// MembershipRequest represents a project membership entry.
type MembershipRequest struct {
	UserID      string `json:"userId"`
	HourlyRate  *int   `json:"hourlyRate,omitempty"`
	MembershipType string `json:"membershipType,omitempty"`
	MembershipStatus string `json:"membershipStatus,omitempty"`
}

// UpdateMembershipsRequest is the payload for updating project memberships.
type UpdateMembershipsRequest struct {
	Memberships []MembershipRequest `json:"memberships"`
}

// AddUsersToProjectRequest adds/removes users from a project.
type AddUsersToProjectRequest struct {
	UserIDs []string `json:"userIds,omitempty"`
	GroupIDs []string `json:"groupIds,omitempty"`
}

// ProjectTemplateRequest creates a project from a template.
type ProjectTemplateRequest struct {
	Name       string `json:"name"`
	TemplateID string `json:"templateId"`
}

// RateRequest sets an hourly or cost rate.
type RateRequest struct {
	Amount   int    `json:"amount"`
	Currency string `json:"currency,omitempty"`
}

// ListProjects returns all projects in a workspace.
func (c *Client) ListProjects(workspaceID string, page, pageSize int) ([]ProjectDTO, error) {
	q := url.Values{}
	if page > 0 {
		q.Set("page", strconv.Itoa(page))
	}
	if pageSize > 0 {
		q.Set("page-size", strconv.Itoa(pageSize))
	}
	var result []ProjectDTO
	if err := c.get(fmt.Sprintf("/workspaces/%s/projects", workspaceID), q, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// CreateProject creates a new project in a workspace.
func (c *Client) CreateProject(workspaceID string, req ProjectRequest) (*ProjectDTO, error) {
	var result ProjectDTO
	if err := c.post(fmt.Sprintf("/workspaces/%s/projects", workspaceID), req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetProject returns a single project by ID.
func (c *Client) GetProject(workspaceID, projectID string) (*ProjectDTO, error) {
	var result ProjectDTO
	if err := c.get(fmt.Sprintf("/workspaces/%s/projects/%s", workspaceID, projectID), nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdateProject updates an existing project.
func (c *Client) UpdateProject(workspaceID, projectID string, req ProjectRequest) (*ProjectDTO, error) {
	var result ProjectDTO
	if err := c.put(fmt.Sprintf("/workspaces/%s/projects/%s", workspaceID, projectID), req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// DeleteProject deletes a project by ID.
func (c *Client) DeleteProject(workspaceID, projectID string) error {
	return c.delete(fmt.Sprintf("/workspaces/%s/projects/%s", workspaceID, projectID), nil)
}

// UpdateProjectEstimate updates the estimate for a project.
func (c *Client) UpdateProjectEstimate(workspaceID, projectID string, req EstimateRequest) (*ProjectDTO, error) {
	var result ProjectDTO
	if err := c.patch(fmt.Sprintf("/workspaces/%s/projects/%s/estimate", workspaceID, projectID), req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdateProjectMemberships updates the memberships of a project.
func (c *Client) UpdateProjectMemberships(workspaceID, projectID string, req UpdateMembershipsRequest) (*ProjectDTO, error) {
	var result ProjectDTO
	if err := c.patch(fmt.Sprintf("/workspaces/%s/projects/%s/memberships", workspaceID, projectID), req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// AddUsersToProject assigns or removes users from a project.
func (c *Client) AddUsersToProject(workspaceID, projectID string, req AddUsersToProjectRequest) (*ProjectDTO, error) {
	var result ProjectDTO
	if err := c.post(fmt.Sprintf("/workspaces/%s/projects/%s/memberships", workspaceID, projectID), req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// CreateProjectFromTemplate creates a project from a template.
func (c *Client) CreateProjectFromTemplate(workspaceID string, req ProjectTemplateRequest) (*ProjectDTO, error) {
	var result ProjectDTO
	if err := c.post(fmt.Sprintf("/workspaces/%s/projects/from-template", workspaceID), req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdateProjectTemplateFlagRequest is used to set/unset a project as a template.
type UpdateProjectTemplateFlagRequest struct {
	IsTemplate bool `json:"isTemplate"`
}

// UpdateProjectTemplateFlag sets or unsets the template flag on a project.
func (c *Client) UpdateProjectTemplateFlag(workspaceID, projectID string, isTemplate bool) (*ProjectDTO, error) {
	req := UpdateProjectTemplateFlagRequest{IsTemplate: isTemplate}
	var result ProjectDTO
	if err := c.patch(fmt.Sprintf("/workspaces/%s/projects/%s/template", workspaceID, projectID), req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// SetProjectUserHourlyRate sets the hourly rate for a user on a project.
func (c *Client) SetProjectUserHourlyRate(workspaceID, projectID, userID string, req RateRequest) error {
	return c.put(fmt.Sprintf("/workspaces/%s/projects/%s/users/%s/hourly-rate", workspaceID, projectID, userID), req, nil)
}

// SetProjectUserCostRate sets the cost rate for a user on a project.
func (c *Client) SetProjectUserCostRate(workspaceID, projectID, userID string, req RateRequest) error {
	return c.put(fmt.Sprintf("/workspaces/%s/projects/%s/users/%s/cost-rate", workspaceID, projectID, userID), req, nil)
}
