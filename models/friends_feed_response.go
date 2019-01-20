package models

type FriendsFeedsResponse struct {
	Feed    []Feed `json:"friends_feed,omitempty"`
	NextUrl string `json:"next_url"`
}
