package smithy_ast

import (
	"bytes"
	"encoding/json"
	"strconv"
	"strings"
)

type AbsShapeID string

func (id *AbsShapeID) Namespace() string {
	s := string(*id)
	i := strings.IndexByte(s, '#')
	if i < 0 {
		return ""
	}
	return s[0:i]
}

func (id *AbsShapeID) Name() string {
	s := string(*id)
	i := strings.IndexByte(s, '#')
	s = s[i+1:]
	i = strings.IndexByte(s, '$')
	if i < 0 {
		return ""
	}
	return s[0:i]
}

func (id *AbsShapeID) Member() string {
	s := string(*id)
	i := strings.IndexByte(s, '#')
	s = s[i+1:]
	i = strings.IndexByte(s, '$')
	return s[i+1:]
}

type AbsShapeIDNode struct {
	node
	Value AbsShapeID
}

func (n *AbsShapeIDNode) Decode(dec *json.Decoder) error {
	offset := dec.InputOffset()
	t, err := dec.Token()
	if err != nil {
		return err
	}
	if s, ok := t.(string); ok {
		n.Value = AbsShapeID(s)
		return nil
	}
	return modelError("expected string [absolute shape ID]", offset)
}

func (n *AbsShapeIDNode) MarshalJSON() ([]byte, error) {
	var b bytes.Buffer
	enc := json.NewEncoder(&b)
	err := enc.Encode(n.Value)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func (n *AbsShapeIDNode) UnmarshalJSON(data []byte) error {
	return unmarshalJSON(data, n)
}

func decodeAbsShapeIDNodeTo(dec *json.Decoder, dst **AbsShapeIDNode) error {
	var id AbsShapeIDNode
	err := id.Decode(dec)
	if err != nil {
		return err
	}
	*dst = &id
	return nil
}

func decodeAbsShapeIDSliceTo(dec *json.Decoder, name string, dst *[]AbsShapeIDNode) error {
	ids := make([]AbsShapeIDNode, 0)
	err := decodeArray(dec, name, func(dec2 *json.Decoder, index int) error {
		var id AbsShapeIDNode
		err2 := id.Decode(dec2)
		if err2 == nil {
			ids = append(ids, id)
			return nil
		}
		if modelErr, ok := err2.(*ModelError); ok {
			modelErr.msg += " in " + name + " [index " + strconv.Itoa(index) + "]"
		}
		return err2
	})
	if err != nil {
		return err
	}
	*dst = ids
	return nil
}
