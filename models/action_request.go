package models

import (
	"fmt"
	"github.com/pkg/errors"
)

type ActionRequest struct {
	Follow *string `json:"follow,omitempty"`
}

func (request *ActionRequest) Validate(actor string) error {
	if request.Follow != nil {
		if len(*request.Follow) == 0 {
			return errors.New("Follow target cannot be empty.")
		}
		if *request.Follow == actor {
			return errors.New(fmt.Sprintf("Actor[%s] cannot follow himself.", actor))
		}
		return nil
	}
	return errors.New(fmt.Sprintf("One of the actions should present. Eligible actions: 'follow'"))

}
