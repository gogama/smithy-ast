package ast

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

// JSONError describes an error in the JSON representation of the Smithy
// AST.
type JSONError struct {
	msg    string // description of error
	Offset int64  // error occurred after reading Offset bytes
}

func isNonSyntaxError(err error) bool {
	if err == nil {
		return false
	}

	_, ok := err.(*json.SyntaxError)
	return ok
}

func jsonError(msg string, offset int64) error {
	return &JSONError{prefix + msg, offset}
}

func (err *JSONError) Error() string {
	return err.msg + " at offset " + strconv.FormatInt(err.Offset, 10)
}

func (err *JSONError) Is(other error) bool {
	if x, ok := other.(*JSONError); ok {
		return *err == *x
	}

	return false
}

type MergeError struct {
}

func unsupportedKeyError(name, key string, offset int64) error {
	return jsonError("unsupported key "+strconv.Quote(key)+" in "+name, offset)
}

func newError(text string) error {
	return errors.New(prefix + text)
}

func newErrorf(format string, a ...interface{}) error {
	return fmt.Errorf(prefix+format, a...)
}

const prefix = "ast: "
