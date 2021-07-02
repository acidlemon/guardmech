package usecase

import (
	"context"
	"fmt"
)

type Context = context.Context
type ErrorType int

const (
	SystemError ErrorType = iota
	SecurityError
	VerificationError
	AuthError
)

type Error struct {
	kind    ErrorType
	message string
	wrap    error
}

func (e Error) Error() string {
	if e.wrap != nil {
		return fmt.Sprintf("%s: %s", e.message, e.wrap.Error())
	}
	return fmt.Sprintf("%s", e.message)
}

func (e Error) Type() ErrorType {
	return e.kind
}

func (e Error) Detail() error {
	return e.wrap
}

func systemError(msg string, err error) error {
	return Error{
		kind:    SystemError,
		message: msg,
		wrap:    err,
	}
}

func securityError(msg string, err error) error {
	return Error{
		kind:    SecurityError,
		message: msg,
		wrap:    err,
	}
}

func verificationError(msg string, err error) error {
	return Error{
		kind:    VerificationError,
		message: msg,
		wrap:    err,
	}
}

func authError(msg string, err error) error {
	return Error{
		kind:    AuthError,
		message: msg,
		wrap:    err,
	}
}
