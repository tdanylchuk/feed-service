package models

import (
	"time"
)

type Feed struct {
	Actor    string    `json:"actor,omitempty"`
	Verb     string    `json:"verb,omitempty"`
	Object   string    `json:"object,omitempty"`
	Target   string    `json:"target,omitempty"`
	Datetime time.Time `json:"datetime,omitempty"`
}
