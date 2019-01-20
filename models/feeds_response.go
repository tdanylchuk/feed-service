package models

type FeedsResponse struct {
	Feed    []FeedResponse `json:"my_feed,omitempty"`
	NextUrl string         `json:"next_url"`
}
