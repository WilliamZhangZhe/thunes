// Package errwrap is a package wrap original error.
package errwrap

import (
	"errors"
	"fmt"
	"runtime"
)

var (
	// ErrNotInit not init
	ErrNotInit = errors.New("not inited")

	// ErrInit 初始化错误
	ErrInit = errors.New("init error")
)

// WrapWithStack wrap error witch stack info
// max bufSize is 1M
func WithStack(err error, bufSize int64) error {
	if err == nil {
		return nil
	}

	if bufSize > 1024*1024 {
		bufSize = 1024 * 1024
	}

	buf := make([]byte, bufSize)
	runtime.Stack(buf, false)

	return withStack{
		wrapped: wrapped{
			error: err,
		},
		Stack: buf,
	}
}

// WithContext 包装err，添加context描述；
// 在fmt.Printx时，将输出 ctx：err
func WithContext(err error, ctx string) error {
	if err == nil {
		return err
	}

	return wrapped{
		error: err,
		Ctx:   ctx,
	}
}

// WithCode 包装err，并添加code信息
func WithCode(code int64, err error) error {
	if err == nil {
		return nil
	}

	return withCode{
		wrapped: wrapped{
			error: err,
		},
		code: code,
	}
}

///////////////////////////////////////////////////////////////////////////////////////////////////

type wrapped struct {
	error

	Ctx string
}

// Unwrap implement of errors.UnWrapper
func (W wrapped) Unwrap() error {
	return W.error
}

// Error implement of interface errors.error
func (W wrapped) Error() string {
	if W.error == nil {
		return ""
	}

	if len(W.Ctx) > 0 {
		return fmt.Sprintf("%s: %s", W.Ctx, W.error)
	}

	return W.error.Error()
}

///////////////////////////////////////////////////////////////////////////////////////////////////

// errwrap api error handler
type withStack struct {
	wrapped

	Stack runtimeStack
}

func (E withStack) Error() string {
	return fmt.Sprintf("%s\n%s", E.wrapped.Error(), E.Stack)
}

//////////////////////////////////////////////////////////////////////////////////////////////////

type runtimeStack []byte

// String implement of interface Stringer
func (S runtimeStack) String() string {
	return string(S)
}

//////////////////////////////////////////////////////////////////////////////////////////////////

type withCode struct {
	wrapped

	code int64
}

func (E withCode) Error() string {
	if E.error == nil {
		return ""
	}

	return fmt.Sprintf("[%d]: %s", E.code, E.wrapped.Error())
}
