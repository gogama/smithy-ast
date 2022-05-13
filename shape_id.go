package smithy_ast

import (
	"bytes"
	"encoding/json"
	"strings"
)

type AbsShapeID string

func (id *AbsShapeID) Namespace() string {
	return namespace(string(*id))
}

func (id *AbsShapeID) Name() string {
	s := string(*id)
	i := strings.IndexByte(s, '#')
	return s[i+1:]
}

func (id *AbsShapeID) decode(dec *json.Decoder) error {
	o := dec.InputOffset()
	t, err := dec.Token()
	if err != nil {
		return err
	}

	var s string
	var ok bool
	if s, ok = t.(string); !ok {
		return modelError("expected string", o)
	}

	i := matchNamespaceHash(s)
	if i < 0 {
		return modelError("expected identifier [namespace] followed by '#' in absolute shape ID", o)
	}
	i = matchIdentifier(s[i:])
	if i < 0 {
		return modelError("expected absolute shape ID to end with identifier [shape name] after '#'", o)
	}

	*id = AbsShapeID(s)
	return nil
}

func (id *AbsShapeID) UnmarshalJSON(data []byte) error {
	dec := json.NewDecoder(bytes.NewReader(data))
	return id.decode(dec)
}

type AbsShapeOrMemberID string

func (id *AbsShapeOrMemberID) Namespace() string {
	return namespace(string(*id))
}

func (id *AbsShapeOrMemberID) Name() string {
	s := string(*id)
	i := strings.IndexByte(s, '#')
	s = s[i+1:]
	i = strings.IndexByte(s, '$')
	if i < 0 {
		return ""
	}
	return s[0:i]
}

func (id *AbsShapeOrMemberID) Member() string {
	s := string(*id)
	i := strings.IndexByte(s, '#')
	s = s[i+1:]
	i = strings.IndexByte(s, '$')
	return s[i+1:]
}

func (id *AbsShapeOrMemberID) decode(dec *json.Decoder) error {
	o := dec.InputOffset()
	t, err := dec.Token()
	if err != nil {
		return err
	}

	var s string
	var ok bool
	if s, ok = t.(string); !ok {
		return modelError("expected string", o)
	}

	i := matchNamespaceHash(s)
	if i < 0 {
		return modelError("expected identifier [namespace] followed by '#' in absolute shape ID", o)
	}
	// TODO: This part

	*id = AbsShapeOrMemberID(s)

}

func (id *AbsShapeOrMemberID) UnmarshalJSON(data []byte) error {
	dec := json.NewDecoder(bytes.NewReader(data))
	return id.decode(dec)
}

func namespace(s string) string {
	i := strings.IndexByte(s, '#')
	if i < 0 {
		return ""
	}
	return s[0:i]
}
