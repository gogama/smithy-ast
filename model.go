package smithy_ast

import (
	"encoding/json"
	"io"
)

type Model struct {
	node
	Version  *StringNode              `json:"version,omitempty"`
	Metadata map[string]InterfaceNode `json:"metadata,omitempty"`
	Shapes   map[AbsShapeID]Shape     `json:"shapes,omitempty"`
}

func (m *Model) Decode(dec *json.Decoder) error {
	return decodeObject(dec, "model", func(dec2 *json.Decoder, key string, keyOffset int64) error {
		switch key {
		case "version":
			return m.Version.Decode(dec2)
		case "metadata":
			m.Metadata = make(map[string]InterfaceNode)
			return decodeToMap(dec2, "metadata", m.Metadata)
		case "shapes":
			m.Shapes = make(map[AbsShapeID]Shape)
			return decodeToMap(dec2, "shapes", m.Shapes)
		default:
			return unsupportedKeyError("model", key, keyOffset)
		}
	})
}

func (m *Model) UnmarshalJSON(data []byte) error {
	return unmarshalJSON(data, m)
}

// TODO: These guys should accept a Config that lets you pass in extra trait node type mappings.
func ReadModel(r io.Reader) (m Model, err error) {
	dec := json.NewDecoder(r)
	err = m.Decode(dec)
	return
}

func WriteModel(m Model, w io.Writer) error {
	enc := json.NewEncoder(w)
	return enc.Encode(m)
}

func MergeModels(m ...Model) (Model, error) {
	// TODO.
	// https://awslabs.github.io/smithy/1.0/spec/core/model.html#merging-models
	return Model{}, nil
}
