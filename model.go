package smithy_ast

import (
	"encoding/json"
	"io"
)

type Model struct {
	node
	Version  StringNode               `json:"version"`
	Metadata map[string]InterfaceNode `json:"metadata"`
	Shapes   map[AbsShapeID]Shape     `json:"shapes"`
}

func (m *Model) Decode(dec *json.Decoder) error {
	return dec.Decode(&m)
}

func (m *Model) UnmarshalJSON(data []byte) error {
	return unmarshalJSON(data, m)
}

// TODO:
func ReadModel(r io.Reader) (m Model, err error) {
	dec := json.NewDecoder(r)
	err = m.Decode(dec)
	return
}

func WriteModel(m Model, w io.Writer) error {
	enc := json.NewEncoder(w)
	return enc.Encode(m)
}
