package mattermost

import "fmt"

func (c *Client4) postsRoute() string {
	return "/posts"
}
func (c *Client4) usersRoute() string {
	return "/users"
}
func (c *Client4) channelMembersRoute(channelId string) string {
	return fmt.Sprintf(c.channelRoute(channelId) + "/members")
}
func (c *Client4) channelRoute(channelId string) string {
	return fmt.Sprintf(c.channelsRoute()+"/%v", channelId)
}
func (c *Client4) channelsRoute() string {
	return "/channels"
}
func (c *Client4) channelMemberRoute(channelId, userId string) string {
	return fmt.Sprintf(c.channelMembersRoute(channelId)+"/%v", userId)
}
func (c *Client4) userRoute(userId string) string {
	return fmt.Sprintf(c.usersRoute()+"/%v", userId)
}
func (c *Client4) teamRoute(teamId string) string {
	return fmt.Sprintf(c.teamsRoute()+"/%v", teamId)
}
func (c *Client4) teamsRoute() string {
	return "/teams"
}
func (c *Client4) teamByNameRoute(teamName string) string {
	return fmt.Sprintf(c.teamsRoute()+"/name/%v", teamName)
}
func (c *Client4) userByUsernameRoute(userName string) string {
	return fmt.Sprintf(c.usersRoute()+"/username/%v", userName)
}

func (c *Client4) userByEmailRoute(email string) string {
	return fmt.Sprintf(c.usersRoute()+"/email/%v", email)
}

func (c *Client4) userThreadsRoute(userID, teamID string) string {
	return c.userRoute(userID) + c.teamRoute(teamID) + "/threads"
}

func (c *Client4) userThreadRoute(userId, teamId, threadId string) string {
	return c.userThreadsRoute(userId, teamId) + "/" + threadId
}
