package models

type FriendsFeedsResponse struct {
	Feed    []FeedResponse `json:"friends_feed,omitempty"`
	NextUrl string         `json:"next_url"`
}
