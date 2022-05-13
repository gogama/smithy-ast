package smithy_ast

import (
	"errors"
	"fmt"
)

// ModelError describes a Smithy AST model structure error.
type ModelError struct {
	msg    string // description of error
	Offset int64  // error occurred after reading Offset bytes
}

func modelError(msg string, offset int64) error {
	return &ModelError{prefix + msg, offset}
}

func (err *ModelError) Error() string { return err.msg }

func newError(text string) error {
	return errors.New(prefix + text)
}

func newErrorf(format string, a ...interface{}) error {
	return fmt.Errorf(prefix+format, a...)
}

type wrapError struct {
	text  string
	cause error
}

func (err *wrapError) Error() string {
	return prefix + err.text + " [" + err.cause.Error() + "]"
}

func (err *wrapError) Unwrap() error {
	return err.cause
}

const prefix = "smithy_ast: "
