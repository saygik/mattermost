package mattermost

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

const (
	ChannelTypeOpen    ChannelType = "O"
	ChannelTypePrivate ChannelType = "P"
	ChannelTypeDirect  ChannelType = "D"
	ChannelTypeGroup   ChannelType = "G"

	ChannelGroupMaxUsers       = 8
	ChannelGroupMinUsers       = 3
	DefaultChannelName         = "town-square"
	ChannelDisplayNameMaxRunes = 64
	ChannelNameMinLength       = 1
	ChannelNameMaxLength       = 64
	ChannelHeaderMaxRunes      = 1024
	ChannelPurposeMaxRunes     = 250
	ChannelCacheSize           = 25000

	ChannelSortByUsername = "username"
	ChannelSortByStatus   = "status"
)

type ChannelType string

type Channel struct {
	Id                string                 `json:"id"`
	CreateAt          int64                  `json:"create_at"`
	UpdateAt          int64                  `json:"update_at"`
	DeleteAt          int64                  `json:"delete_at"`
	TeamId            string                 `json:"team_id"`
	Type              ChannelType            `json:"type"`
	DisplayName       string                 `json:"display_name"`
	Name              string                 `json:"name"`
	Header            string                 `json:"header"`
	Purpose           string                 `json:"purpose"`
	LastPostAt        int64                  `json:"last_post_at"`
	TotalMsgCount     int64                  `json:"total_msg_count"`
	ExtraUpdateAt     int64                  `json:"extra_update_at"`
	CreatorId         string                 `json:"creator_id"`
	SchemeId          *string                `json:"scheme_id"`
	Props             map[string]interface{} `json:"props"`
	GroupConstrained  *bool                  `json:"group_constrained"`
	Shared            *bool                  `json:"shared"`
	TotalMsgCountRoot int64                  `json:"total_msg_count_root"`
	PolicyID          *string                `json:"policy_id"`
	LastRootPostAt    int64                  `json:"last_root_post_at"`
}
type ChannelMember struct {
	ChannelId        string    `json:"channel_id"`
	UserId           string    `json:"user_id"`
	Roles            string    `json:"roles"`
	LastViewedAt     int64     `json:"last_viewed_at"`
	MsgCount         int64     `json:"msg_count"`
	MentionCount     int64     `json:"mention_count"`
	MentionCountRoot int64     `json:"mention_count_root"`
	MsgCountRoot     int64     `json:"msg_count_root"`
	NotifyProps      StringMap `json:"notify_props"`
	LastUpdateAt     int64     `json:"last_update_at"`
	SchemeGuest      bool      `json:"scheme_guest"`
	SchemeUser       bool      `json:"scheme_user"`
	SchemeAdmin      bool      `json:"scheme_admin"`
	ExplicitRoles    string    `json:"explicit_roles"`
}
type ChannelMembers []ChannelMember

// CreateDirectChannel creates a direct message channel based on the two user
// ids provided.
func (c *Client4) CreateDirectChannel(userId1, userId2 string) (*Channel, *Response, error) {
	requestBody := []string{userId1, userId2}
	r, err := c.DoAPIPost(c.channelsRoute()+"/direct", ArrayToJSON(requestBody))
	if err != nil {
		return nil, BuildResponse(r), err
	}
	defer closeBody(r)

	var ch *Channel
	err = json.NewDecoder(r.Body).Decode(&ch)
	if err != nil {
		return nil, BuildResponse(r), NewAppError("CreateDirectChannel", "api.marshal_error", nil, "", http.StatusInternalServerError).Wrap(err)
	}
	return ch, BuildResponse(r), nil
}

// GetChannelMember gets a channel member.
func (c *Client4) GetChannelMember(channelId, userId, etag string) (*ChannelMember, *Response, error) {
	r, err := c.DoAPIGet(c.channelMemberRoute(channelId, userId), etag)
	if err != nil {
		return nil, BuildResponse(r), err
	}
	defer closeBody(r)

	var ch *ChannelMember
	err = json.NewDecoder(r.Body).Decode(&ch)
	if err != nil {
		return nil, BuildResponse(r), NewAppError("GetChannelMember", "api.marshal_error", nil, "", http.StatusInternalServerError).Wrap(err)
	}
	return ch, BuildResponse(r), nil
}

// GetChannelMembers gets a page of channel members specific to a channel.
func (c *Client4) GetChannelMembers(channelId string, page, perPage int, etag string) (ChannelMembers, *Response, error) {
	query := fmt.Sprintf("?page=%v&per_page=%v", page, perPage)
	r, err := c.DoAPIGet(c.channelMembersRoute(channelId)+query, etag)
	if err != nil {
		return nil, BuildResponse(r), err
	}
	defer closeBody(r)

	var ch ChannelMembers
	err = json.NewDecoder(r.Body).Decode(&ch)
	if err != nil {
		return nil, BuildResponse(r), NewAppError("GetChannelMembers", "api.marshal_error", nil, "", http.StatusInternalServerError).Wrap(err)
	}
	return ch, BuildResponse(r), nil
}

func PrepareChannelId(c *Client4, mattermostChannel string) (string, error) {

	if strings.HasPrefix(mattermostChannel, "@") {
		userFrom, _, err := c.GetMe("")
		if err != nil {
			return "", err
		}
		userTo, _, err := c.GetUserByUsername(strings.TrimLeft(mattermostChannel, "@"), "")
		if err != nil {
			return "", err
		}
		channel, _, err := c.CreateDirectChannel(userFrom.Id, userTo.Id)
		if err != nil {
			return "", err
		}
		return channel.Id, nil
	} else {
		return mattermostChannel, nil
	}
}

// GetChannel returns a channel based on the provided channel id string.
func (c *Client4) GetChannel(channelId, etag string) (*Channel, *Response, error) {
	r, err := c.DoAPIGet(c.channelRoute(channelId), etag)
	if err != nil {
		return nil, BuildResponse(r), err
	}
	defer closeBody(r)

	var ch *Channel
	err = json.NewDecoder(r.Body).Decode(&ch)
	if err != nil {
		return nil, BuildResponse(r), NewAppError("GetChannel", "api.marshal_error", nil, "", http.StatusInternalServerError).Wrap(err)
	}
	return ch, BuildResponse(r), nil
}
