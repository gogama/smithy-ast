package ast

import (
	"errors"
	"fmt"
	"strconv"
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

func unsupportedKeyError(name, key string, offset int64) error {
	return modelError("unsupported key "+strconv.Quote(key)+" in "+name, offset)
}

func newError(text string) error {
	return errors.New(prefix + text)
}

func newErrorf(format string, a ...interface{}) error {
	return fmt.Errorf(prefix+format, a...)
}

const prefix = "smithy_ast: "
