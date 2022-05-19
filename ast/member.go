package ast

import "encoding/json"

type Member struct {
	node
	Target AbsShapeIDNode `json:"target"`
	Traits Traits         `json:"traits,omitempty"`
}

func (m *Member) Decode(dec *json.Decoder) error {
	return decodeObject(dec, "member", func(dec2 *json.Decoder, key string, keyOffset int64) error {
		switch key {
		case "target":
			return m.Target.Decode(dec2)
		case "traits":
			return m.Traits.decode(dec2)
		default:
			return unsupportedKeyError("member", key, keyOffset)
		}
	})
}

func (m *Member) UnmarshalJSON(data []byte) error {
	return unmarshalJSON(data, m)
}

func (m *Member) ResolveTraits(model Model) (Traits, error) {
	targetShape, ok := model.Shapes[m.Target.Value]
	if !ok {
		return nil, newErrorf("model does not contain member target shape %s", m.Target.Value)
	}

	if len(m.Traits) == 0 {
		return targetShape.Traits, nil
	} else if len(targetShape.Traits) == 0 {
		return m.Traits, nil
	} else {
		merged := make(Traits, len(m.Traits)+len(targetShape.Traits))
		for k, v := range targetShape.Traits {
			merged[k] = v
		}
		for k, v := range m.Traits {
			merged[k] = v
		}
		return merged, nil
	}
}
