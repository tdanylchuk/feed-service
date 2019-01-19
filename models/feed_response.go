package models

type FeedsResponse struct {
	Feed    []Feed `json:"my_feed"`
	NextUrl string `json:"next_url"`
}
