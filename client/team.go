package mattermost

import (
	"encoding/json"
	"net/http"
)

type Team struct {
	Id                  string  `json:"id"`
	CreateAt            int64   `json:"create_at"`
	UpdateAt            int64   `json:"update_at"`
	DeleteAt            int64   `json:"delete_at"`
	DisplayName         string  `json:"display_name"`
	Name                string  `json:"name"`
	Description         string  `json:"description"`
	Email               string  `json:"email"`
	Type                string  `json:"type"`
	CompanyName         string  `json:"company_name"`
	AllowedDomains      string  `json:"allowed_domains"`
	InviteId            string  `json:"invite_id"`
	AllowOpenInvite     bool    `json:"allow_open_invite"`
	LastTeamIconUpdate  int64   `json:"last_team_icon_update,omitempty"`
	SchemeId            *string `json:"scheme_id"`
	GroupConstrained    *bool   `json:"group_constrained"`
	PolicyID            *string `json:"policy_id"`
	CloudLimitsArchived bool    `json:"cloud_limits_archived"`
}

// GetTeamByName returns a team based on the provided team name string.
func (c *Client4) GetTeamByName(name, etag string) (*Team, *Response, error) {
	r, err := c.DoAPIGet(c.teamByNameRoute(name), etag)
	if err != nil {
		return nil, BuildResponse(r), err
	}
	defer closeBody(r)
	var t Team
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		return nil, nil, NewAppError("GetTeamByName", "api.unmarshal_error", nil, "", http.StatusInternalServerError).Wrap(err)
	}
	return &t, BuildResponse(r), nil
}
