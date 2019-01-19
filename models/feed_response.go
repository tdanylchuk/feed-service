package models

type FeedsResponse struct {
	Feed    []Feed `json:"my_feed,omitempty"`
	NextUrl string `json:"next_url"`
}
