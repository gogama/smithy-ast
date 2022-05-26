package ast

import (
	"bytes"
	"encoding/json"
	"math/big"
	"strconv"
)

type Location struct {
	Path   string
	Offset int
	Row    int
	Col    int
}

func (loc *Location) IsEmpty() bool {
	return *loc == Location{}
}

type Node interface {
	Location() Location
	SetLocation(loc Location)
	Decode(dec *json.Decoder) error
}

type node struct {
	loc Location
}

func (n *node) Location() Location {
	return n.loc
}

func (n *node) SetLocation(loc Location) {
	n.loc = loc
}

func unmarshalJSON(data []byte, n Node) error {
	dec := json.NewDecoder(bytes.NewReader(data))
	return n.Decode(dec)
}

type InterfaceNode struct {
	node
	Value interface{}
}

func (n *InterfaceNode) Decode(dec *json.Decoder) error {
	return dec.Decode(&n.Value)
}

func (n *InterfaceNode) UnmarshalJSON(data []byte) error {
	return unmarshalJSON(data, n)
}

func (n InterfaceNode) MarshalJSON() ([]byte, error) {
	return json.Marshal(n.Value)
}

type StringNode struct {
	node
	Value string
}

func (n *StringNode) Decode(dec *json.Decoder) error {
	offset := dec.InputOffset()
	t, err := dec.Token()
	if isNonSyntaxError(err) {
		return err
	}
	if s, ok := t.(string); ok {
		n.Value = s
		return nil
	}
	return jsonError("expected string", offset)
}

func (n *StringNode) UnmarshalJSON(data []byte) error {
	return unmarshalJSON(data, n)
}

func (n StringNode) MarshalJSON() ([]byte, error) {
	return json.Marshal(n.Value)
}

type BoolNode struct {
	node
	Value bool
}

func (n *BoolNode) Decode(dec *json.Decoder) error {
	offset := dec.InputOffset()
	t, err := dec.Token()
	if isNonSyntaxError(err) {
		return err
	}
	if b, ok := t.(bool); ok {
		n.Value = b
		return nil
	}
	return jsonError("expected boolean", offset)
}

func (n *BoolNode) UnmarshalJSON(data []byte) error {
	return unmarshalJSON(data, n)
}

func (n BoolNode) MarshalJSON() ([]byte, error) {
	return json.Marshal(n.Value)
}

type Int32Node struct {
	node
	Value int32
}

func (n *Int32Node) Decode(dec *json.Decoder) error {
	return decodeNumber(dec, func(s string) error {
		i, err := strconv.ParseInt(s, 10, 32)
		if err != nil {
			return err
		}
		n.Value = int32(i)
		return nil
	})
}

func (n *Int32Node) UnmarshalJSON(data []byte) error {
	return unmarshalJSON(data, n)
}

func (n Int32Node) MarshalJSON() ([]byte, error) {
	return json.Marshal(n.Value)
}

type Int64Node struct {
	node
	Value int64
}

func (n *Int64Node) Decode(dec *json.Decoder) error {
	return decodeNumber(dec, func(s string) error {
		i, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return err
		}
		n.Value = i
		return nil
	})
}

func (n *Int64Node) UnmarshalJSON(data []byte) error {
	return unmarshalJSON(data, n)
}

func (n Int64Node) MarshalJSON() ([]byte, error) {
	return json.Marshal(n.Value)
}

type BigFloatNode struct {
	node
	Value big.Float
}

func (n *BigFloatNode) Decode(dec *json.Decoder) error {
	return decodeNumber(dec, func(s string) error {
		return n.Value.UnmarshalText([]byte(s))
	})
}

func (n *BigFloatNode) UnmarshalJSON(data []byte) error {
	return unmarshalJSON(data, n)
}

func (n BigFloatNode) MarshalJSON() ([]byte, error) {
	return json.Marshal(n.Value)
}
