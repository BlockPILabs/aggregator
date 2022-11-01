package aggregator

import (
	"fmt"
)

type Error struct {
	Code    int
	Message string
}

func (err *Error) Error() string {
	return fmt.Sprintf("[%d]%s", err.Code, err.Message)
}

func NewError(code int, msg string) *Error {
	return &Error{Code: code, Message: msg}
}

var (
	ErrServerError    = NewError(-32000, "server error")
	ErrInvalidRequest = NewError(-32600, "invalid request")
	ErrInvalidMethod  = NewError(-32601, "invalid method")
	ErrInvalidChain   = NewError(-32601, "invalid chain")
	ErrMustReturn     = NewError(10000, "must return")
	ErrDenyRequest    = NewError(-32602, "deny request")
)
