package client

import (
	"errors"

	"github.com/ogen-go/ogen/validate"
)

func IsNotFoundError(err error) bool {
	var statusErr *validate.UnexpectedStatusCodeError
	if errors.As(err, &statusErr) && statusErr.StatusCode == 404 {
		return true
	}
	return false
}
