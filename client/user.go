package mattermost

import (
	"encoding/json"
	"net/http"
)

const (
	Me                             = "me"
	UserNotifyAll                  = "all"
	UserNotifyHere                 = "here"
	UserNotifyMention              = "mention"
	UserNotifyNone                 = "none"
	DesktopNotifyProp              = "desktop"
	DesktopSoundNotifyProp         = "desktop_sound"
	MarkUnreadNotifyProp           = "mark_unread"
	PushNotifyProp                 = "push"
	PushStatusNotifyProp           = "push_status"
	EmailNotifyProp                = "email"
	ChannelMentionsNotifyProp      = "channel"
	CommentsNotifyProp             = "comments"
	MentionKeysNotifyProp          = "mention_keys"
	CommentsNotifyNever            = "never"
	CommentsNotifyRoot             = "root"
	CommentsNotifyAny              = "any"
	CommentsNotifyCRT              = "crt"
	FirstNameNotifyProp            = "first_name"
	AutoResponderActiveNotifyProp  = "auto_responder_active"
	AutoResponderMessageNotifyProp = "auto_responder_message"
	DesktopThreadsNotifyProp       = "desktop_threads"
	PushThreadsNotifyProp          = "push_threads"
	EmailThreadsNotifyProp         = "email_threads"

	DefaultLocale        = "en"
	UserAuthServiceEmail = "email"

	UserEmailMaxLength    = 128
	UserNicknameMaxRunes  = 64
	UserPositionMaxRunes  = 128
	UserFirstNameMaxRunes = 64
	UserLastNameMaxRunes  = 64
	UserAuthDataMaxLength = 128
	UserNameMaxLength     = 64
	UserNameMinLength     = 1
	UserPasswordMaxLength = 72
	UserLocaleMaxLength   = 5
	UserTimezoneMaxRunes  = 256
	UserRolesMaxLength    = 256
)

//msgp:tuple User

// User contains the details about the user.
// This struct's serializer methods are auto-generated. If a new field is added/removed,
// please run make gen-serialized.
type User struct {
	Id                     string    `json:"id"`
	CreateAt               int64     `json:"create_at,omitempty"`
	UpdateAt               int64     `json:"update_at,omitempty"`
	DeleteAt               int64     `json:"delete_at"`
	Username               string    `json:"username"`
	Password               string    `json:"password,omitempty"`
	AuthData               *string   `json:"auth_data,omitempty"`
	AuthService            string    `json:"auth_service"`
	Email                  string    `json:"email"`
	EmailVerified          bool      `json:"email_verified,omitempty"`
	Nickname               string    `json:"nickname"`
	FirstName              string    `json:"first_name"`
	LastName               string    `json:"last_name"`
	Position               string    `json:"position"`
	Roles                  string    `json:"roles"`
	AllowMarketing         bool      `json:"allow_marketing,omitempty"`
	Props                  StringMap `json:"props,omitempty"`
	NotifyProps            StringMap `json:"notify_props,omitempty"`
	LastPasswordUpdate     int64     `json:"last_password_update,omitempty"`
	LastPictureUpdate      int64     `json:"last_picture_update,omitempty"`
	FailedAttempts         int       `json:"failed_attempts,omitempty"`
	Locale                 string    `json:"locale"`
	Timezone               StringMap `json:"timezone"`
	MfaActive              bool      `json:"mfa_active,omitempty"`
	MfaSecret              string    `json:"mfa_secret,omitempty"`
	RemoteId               *string   `json:"remote_id,omitempty"`
	LastActivityAt         int64     `json:"last_activity_at,omitempty"`
	IsBot                  bool      `json:"is_bot,omitempty"`
	BotDescription         string    `json:"bot_description,omitempty"`
	BotLastIconUpdate      int64     `json:"bot_last_icon_update,omitempty"`
	TermsOfServiceId       string    `json:"terms_of_service_id,omitempty"`
	TermsOfServiceCreateAt int64     `json:"terms_of_service_create_at,omitempty"`
	DisableWelcomeEmail    bool      `json:"disable_welcome_email"`
}

// UserMap is a map from a userId to a user object.
// It is used to generate methods which can be used for fast serialization/de-serialization.
type UserMap map[string]*User

func (c *Client4) GetMe(etag string) (*User, *Response, error) {
	r, err := c.DoAPIGet(c.userRoute(Me), etag)
	if err != nil {
		return nil, BuildResponse(r), err
	}
	defer closeBody(r)
	var u User
	if r.StatusCode == http.StatusNotModified {
		return &u, BuildResponse(r), nil
	}
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		return nil, nil, NewAppError("GetMe", "api.unmarshal_error", nil, "", http.StatusInternalServerError).Wrap(err)
	}
	return &u, BuildResponse(r), nil
}

// GetUser returns a user based on the provided user id string.
func (c *Client4) GetUser(userId, etag string) (*User, *Response, error) {
	r, err := c.DoAPIGet(c.userRoute(userId), etag)
	if err != nil {
		return nil, BuildResponse(r), err
	}
	defer closeBody(r)
	var u User
	if r.StatusCode == http.StatusNotModified {
		return &u, BuildResponse(r), nil
	}
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		return nil, nil, NewAppError("GetUser", "api.unmarshal_error", nil, "", http.StatusInternalServerError).Wrap(err)
	}
	return &u, BuildResponse(r), nil
}

// GetUserByUsername returns a user based on the provided user name string.
func (c *Client4) GetUserByUsername(userName, etag string) (*User, *Response, error) {
	r, err := c.DoAPIGet(c.userByUsernameRoute(userName), etag)
	if err != nil {
		return nil, BuildResponse(r), err
	}
	defer closeBody(r)
	var u User
	if r.StatusCode == http.StatusNotModified {
		return &u, BuildResponse(r), nil
	}
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		return nil, nil, NewAppError("GetUserByUsername", "api.unmarshal_error", nil, "", http.StatusInternalServerError).Wrap(err)
	}
	return &u, BuildResponse(r), nil
}

func (c *Client4) UpdateThreadFollowForUser(userId, teamId, threadId string, state bool) (*Response, error) {
	var err error
	var r *http.Response
	if state {
		r, err = c.DoAPIPut(c.userThreadRoute(userId, teamId, threadId)+"/following", "")
	} else {
		r, err = c.DoAPIDelete(c.userThreadRoute(userId, teamId, threadId) + "/following")
	}
	if err != nil {
		return BuildResponse(r), err
	}
	defer closeBody(r)

	return BuildResponse(r), nil
}

func (c *Client4) UpdateThreadFollowAllUsersInChannel(channelID, postID string, state bool) error {
	channel, _, err := c.GetChannel(channelID, "")
	if err != nil {
		return err
	}
	users, _, err := c.GetChannelMembers(channelID, 0, 200, "")
	if err != nil {
		return err
	}
	for _, user := range users {
		_, _ = c.UpdateThreadFollowForUser(user.UserId, channel.TeamId, postID, state)
	}
	return nil
}
