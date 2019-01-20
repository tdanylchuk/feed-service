package models

import (
	"fmt"
	"github.com/pkg/errors"
)

type ActionRequest struct {
	Follow   *string `json:"follow,omitempty"`
	Unfollow *string `json:"unfollow,omitempty"`
}

func (request *ActionRequest) Validate(actor string) error {
	follow := request.Follow
	if follow != nil {
		return validateTargetValue(follow, actor, "follow")
	}
	unfollow := request.Unfollow
	if unfollow != nil {
		return validateTargetValue(unfollow, actor, "unfollow")
	}
	return errors.New(fmt.Sprintf("One of the actions should present. Eligible actions: 'follow','unfollow'"))
}

func validateTargetValue(target *string, actor string, description string) error {
	if len(*target) == 0 {
		return errors.New(fmt.Sprintf("%s target cannot be empty.", description))
	}
	if *target == actor {
		return errors.New(fmt.Sprintf("Actor[%s] cannot %s himself.", actor, description))
	}
	return nil
}
