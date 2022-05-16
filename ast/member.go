package ast

import "encoding/json"

type Member struct {
	node
	Target AbsShapeID
	Traits Traits
}

func (m *Member) Decode(dec *json.Decoder) error {
	// TODO
	return nil
}

func (m *Member) UnmarshalJSON(data []byte) error {
	// TODO
	return nil
}
