package pkg

import "github.com/pkg/errors"

func MaskErr(err error) error { // used to get stack trace from error
	return errors.New(err.Error())
}
