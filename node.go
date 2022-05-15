package smithy_ast

import (
	"bytes"
	"encoding/json"
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
	json.Unmarshaler
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

type StringNode struct {
	node
	Value string
}

func (n *StringNode) Decode(dec *json.Decoder) error {
	offset := dec.InputOffset()
	t, err := dec.Token()
	if err != nil {
		return err
	}
	if s, ok := t.(string); ok {
		n.Value = s
		return nil
	}
	return modelError("expected string", offset)
}

func (n *StringNode) MarshalJSON() ([]byte, error) {
	var b bytes.Buffer
	enc := json.NewEncoder(&b)
	err := enc.Encode(n.Value)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func (n *StringNode) UnmarshalJSON(data []byte) error {
	return unmarshalJSON(data, n)
}
