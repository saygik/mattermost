// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package mattermost

import (
	"sync"
)

const (
	PostSystemMessagePrefix        = "system_"
	PostTypeDefault                = ""
	PostTypeSlackAttachment        = "slack_attachment"
	PostTypeSystemGeneric          = "system_generic"
	PostTypeJoinLeave              = "system_join_leave" // Deprecated, use PostJoinChannel or PostLeaveChannel instead
	PostTypeJoinChannel            = "system_join_channel"
	PostTypeGuestJoinChannel       = "system_guest_join_channel"
	PostTypeLeaveChannel           = "system_leave_channel"
	PostTypeJoinTeam               = "system_join_team"
	PostTypeLeaveTeam              = "system_leave_team"
	PostTypeAutoResponder          = "system_auto_responder"
	PostTypeAddRemove              = "system_add_remove" // Deprecated, use PostAddToChannel or PostRemoveFromChannel instead
	PostTypeAddToChannel           = "system_add_to_channel"
	PostTypeAddGuestToChannel      = "system_add_guest_to_chan"
	PostTypeRemoveFromChannel      = "system_remove_from_channel"
	PostTypeMoveChannel            = "system_move_channel"
	PostTypeAddToTeam              = "system_add_to_team"
	PostTypeRemoveFromTeam         = "system_remove_from_team"
	PostTypeHeaderChange           = "system_header_change"
	PostTypeDisplaynameChange      = "system_displayname_change"
	PostTypeConvertChannel         = "system_convert_channel"
	PostTypePurposeChange          = "system_purpose_change"
	PostTypeChannelDeleted         = "system_channel_deleted"
	PostTypeChannelRestored        = "system_channel_restored"
	PostTypeEphemeral              = "system_ephemeral"
	PostTypeChangeChannelPrivacy   = "system_change_chan_privacy"
	PostTypeAddBotTeamsChannels    = "add_bot_teams_channels"
	PostTypeSystemWarnMetricStatus = "warn_metric_status"
	PostTypeMe                     = "me"
	PostCustomTypePrefix           = "custom_"
	PostTypeReminder               = "reminder"

	PostFileidsMaxRunes   = 300
	PostFilenamesMaxRunes = 4000
	PostHashtagsMaxRunes  = 1000
	PostMessageMaxRunesV1 = 4000
	PostMessageMaxBytesV2 = 65535                     // Maximum size of a TEXT column in MySQL
	PostMessageMaxRunesV2 = PostMessageMaxBytesV2 / 4 // Assume a worst-case representation
	PostPropsMaxRunes     = 800000
	PostPropsMaxUserRunes = PostPropsMaxRunes - 40000 // Leave some room for system / pre-save modifications

	PropsAddChannelMember = "add_channel_member"

	PostPropsAddedUserId       = "addedUserId"
	PostPropsDeleteBy          = "deleteBy"
	PostPropsOverrideIconURL   = "override_icon_url"
	PostPropsOverrideIconEmoji = "override_icon_emoji"

	PostPropsMentionHighlightDisabled = "mentionHighlightDisabled"
	PostPropsGroupHighlightDisabled   = "disable_group_highlight"

	PostPropsPreviewedPost = "previewed_post"
)

const (
	ModifierMessages string = "messages"
	ModifierFiles    string = "files"
)

type StringInterface map[string]any
type StringArray []string

type Post struct {
	Id         string `json:"id"`
	CreateAt   int64  `json:"create_at"`
	UpdateAt   int64  `json:"update_at"`
	EditAt     int64  `json:"edit_at"`
	DeleteAt   int64  `json:"delete_at"`
	IsPinned   bool   `json:"is_pinned"`
	UserId     string `json:"user_id"`
	ChannelId  string `json:"channel_id"`
	RootId     string `json:"root_id"`
	OriginalId string `json:"original_id"`

	Message string `json:"message"`
	// MessageSource will contain the message as submitted by the user if Message has been modified
	// by Mattermost for presentation (e.g if an image proxy is being used). It should be used to
	// populate edit boxes if present.
	MessageSource string `json:"message_source,omitempty"`

	Type          string          `json:"type"`
	propsMu       sync.RWMutex    `db:"-"`       // Unexported mutex used to guard Post.Props.
	Props         StringInterface `json:"props"` // Deprecated: use GetProps()
	Hashtags      string          `json:"hashtags"`
	Filenames     StringArray     `json:"-"` // Deprecated, do not use this field any more
	FileIds       StringArray     `json:"file_ids,omitempty"`
	PendingPostId string          `json:"pending_post_id"`
	HasReactions  bool            `json:"has_reactions,omitempty"`
	RemoteId      *string         `json:"remote_id,omitempty"`

	// Transient data populated before sending a post to the client
	ReplyCount  int64 `json:"reply_count"`
	LastReplyAt int64 `json:"last_reply_at"`
	IsFollowing *bool `json:"is_following,omitempty"` // for root posts in collapsed thread mode indicates if the current user is following this thread
}
