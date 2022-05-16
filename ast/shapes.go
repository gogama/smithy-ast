package ast

import (
	"bytes"
	"encoding/json"
	"strconv"
)

type Shape struct {
	node
	Type      ShapeType
	Traits    Traits
	Key       *Member
	Value     *Member
	Members   map[string]Member
	Service   *Service
	Resource  *Resource
	Operation *Operation
}

func (s *Shape) Decode(dec *json.Decoder) error {
	offset := dec.InputOffset()

	var t *ShapeType
	var traits Traits
	var buf shapeBuffer

	// Field decode phase. We have to buffer the decoded members
	// because ordering of object keys is not guaranteed in JSON, and
	// therefore the "type" member might be at the end.
	err := decodeObject(dec, "shape", func(dec2 *json.Decoder, key string, keyOffset int64) error {
		switch key {
		case "type":
			tok, err2 := dec2.Token()
			if err2 != nil {
				return err2
			}
			if s2, ok := tok.(string); ok {
				if !ShapeTypes[ShapeType(s2)] {
					return modelError("unrecognized shape type: "+strconv.Quote(s2), offset)
				}
				*t = ShapeType(s2)
			}
			return modelError("expected string [shape type]", offset)
		case "traits":
			err2 := traits.decode(dec2)
			if err2 != nil {
				return err2
			}
		default:
			f, ok := shapeFields[key]
			if !ok {
				return modelError("unrecognized shape key: "+strconv.Quote(key), keyOffset)
			}
			err2 := f.decodeFunc(dec2, &buf)
			if err2 != nil {
				return err2
			}
			buf.fields = append(buf.fields, f)
		}
		return nil
	})

	// Short-circuit if a field failed to decode.
	if err != nil {
		return err
	}

	// Validate that a shape type was received.
	if t == nil {
		return modelError("shape is missing type field", offset)
	}

	// Validate that all shape fields decoded are valid members of the
	// shape type specified.
	for i := range buf.fields {
		found := false
		for j := range buf.fields[i].types {
			if buf.fields[i].types[j] == *t {
				found = true
				break
			}
		}
		if !found {
			return modelError("shape of type "+string(*t)+" contains unsupported field "+strconv.Quote(buf.fields[i].name), offset)
		}
	}

	// Store the shape fields from the buffer onto the final shape.
	for i := range buf.fields {
		buf.fields[i].storeFunc(*t, &buf, s)
	}

	// Shape decoded successfully.
	return nil
}

func (s *Shape) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer

	_, _ = buf.WriteString(`{"type":`)
	p, _ := json.Marshal(s.Type)
	_, _ = buf.Write(p)

	if len(s.Traits) > 0 {
		_, _ = buf.WriteString(`,"traits":`)
		p, _ = json.Marshal(s.Traits)
		_, _ = buf.Write(p)
	}

	switch s.Type {
	case ListType, SetType:
		_, _ = buf.WriteString(`,"member":`)
		p, _ = json.Marshal(s.Value)
		_, _ = buf.Write(p)
	case MapType:
		_, _ = buf.WriteString(`,"key":`)
		p, _ = json.Marshal(s.Key)
		_, _ = buf.Write(p)
		_, _ = buf.WriteString(`,"value":`)
		p, _ = json.Marshal(s.Value)
		_, _ = buf.Write(p)
	case StructureType, UnionType:
		_, _ = buf.WriteString(`,"members":`)
		p, _ = json.Marshal(s.Members)
		_, _ = buf.Write(p)
	case ServiceType:
		_ = buf.WriteByte(',')
		p, _ = json.Marshal(s.Service)    // Marshals as a JSON object
		_, _ = buf.Write(p[1 : len(p)-1]) // Remove surrounding braces and just take key/value pairs.
	case ResourceType:
		_ = buf.WriteByte(',')
		p, _ = json.Marshal(s.Resource)   // Marshals as a JSON object
		_, _ = buf.Write(p[1 : len(p)-1]) // Remove surrounding braces and just take key/value pairs.
	case OperationType:
		_ = buf.WriteByte(',')
		p, _ = json.Marshal(s.Operation)  // Marshals as a JSON object
		_, _ = buf.Write(p[1 : len(p)-1]) // Remove surrounding braces and just take key/value pairs.
	}

	return buf.Bytes(), nil
}

func (s *Shape) UnmarshalJSON(data []byte) error {
	return unmarshalJSON(data, s)
}

func (s *Shape) service() *Service {
	if s.Service == nil {
		s.Service = &Service{}
	}
	return s.Service
}

func (s *Shape) resource() *Resource {
	if s.Resource == nil {
		s.Resource = &Resource{}
	}
	return s.Resource
}

func (s *Shape) operation() *Operation {
	if s.Operation == nil {
		s.Operation = &Operation{}
	}
	return s.Operation
}

type Service struct {
	Version    StringNode                `json:"version,omitempty"`
	Operations []AbsShapeIDNode          `json:"operations,omitempty"`
	Resources  []AbsShapeIDNode          `json:"resources,omitempty"`
	Errors     []AbsShapeIDNode          `json:"errors,omitempty"`
	Rename     map[AbsShapeID]StringNode `json:"rename,omitempty"`
}

type Resource struct {
	Identifiers          map[string]AbsShapeIDNode `json:"identifiers,omitempty"`
	Create               *AbsShapeIDNode           `json:"create,omitempty"`
	Put                  *AbsShapeIDNode           `json:"put,omitempty"`
	Read                 *AbsShapeIDNode           `json:"read,omitempty"`
	Update               *AbsShapeIDNode           `json:"update,omitempty"`
	Delete               *AbsShapeIDNode           `json:"delete,omitempty"`
	List                 *AbsShapeIDNode           `json:"list,omitempty"`
	Operations           []AbsShapeIDNode          `json:"operations,omitempty"`
	CollectionOperations []AbsShapeIDNode          `json:"collectionOperations,omitempty"`
	Resources            []AbsShapeIDNode          `json:"resources ,omitempty"`
}

type Operation struct {
	Input  *AbsShapeIDNode  `json:"input,omitempty"`
	Output *AbsShapeIDNode  `json:"output,omitempty"`
	Errors []AbsShapeIDNode `json:"errors,omitempty"`
}

type shapeField struct {
	name       string
	types      []ShapeType
	storeFunc  func(t ShapeType, src *shapeBuffer, dst *Shape)
	decodeFunc func(dec *json.Decoder, dst *shapeBuffer) error
}

var shapeFields = map[string]shapeField{
	"member": {
		name:  "member",
		types: []ShapeType{ListType, SetType},
		storeFunc: func(_ ShapeType, src *shapeBuffer, dst *Shape) {
			dst.Value = src.value
		},
		decodeFunc: func(dec *json.Decoder, dst *shapeBuffer) error {
			return dst.value.Decode(dec)
		},
	},
	"key": {
		name:  "key",
		types: []ShapeType{MapType},
		storeFunc: func(_ ShapeType, src *shapeBuffer, dst *Shape) {
			dst.Key = src.key
		},
		decodeFunc: func(dec *json.Decoder, dst *shapeBuffer) error {
			return dst.key.Decode(dec)
		},
	},
	"value": {
		name:  "value",
		types: []ShapeType{MapType},
		storeFunc: func(_ ShapeType, src *shapeBuffer, dst *Shape) {
			dst.Value = src.value
		},
		decodeFunc: func(dec *json.Decoder, dst *shapeBuffer) error {
			return dst.value.Decode(dec)
		},
	},
	"members": {
		name:  "members",
		types: []ShapeType{StructureType, UnionType},
		storeFunc: func(_ ShapeType, src *shapeBuffer, dst *Shape) {
			dst.Members = src.members
		},
		decodeFunc: func(dec *json.Decoder, dst *shapeBuffer) error {
			dst.members = make(map[string]Member)
			return decodeObject(dec, "structure/union members", func(dec2 *json.Decoder, key string, offset int64) error {
				var m Member
				err := m.Decode(dec2)
				if err != nil {
					return err
				}
				dst.members[key] = m
				return nil
			})
		},
	},
	"version": {
		name:  "version",
		types: []ShapeType{ServiceType},
		storeFunc: func(_ ShapeType, src *shapeBuffer, dst *Shape) {
			dst.service().Version = src.version
		},
		decodeFunc: func(dec *json.Decoder, dst *shapeBuffer) error {
			return dst.version.Decode(dec)
		},
	},
	"operations": {
		name:  "operations",
		types: []ShapeType{ServiceType, ResourceType},
		storeFunc: func(t ShapeType, src *shapeBuffer, dst *Shape) {
			var o *[]AbsShapeIDNode
			switch t {
			case ServiceType:
				o = &dst.service().Operations
			case ResourceType:
				o = &dst.resource().Operations
			}
			*o = src.operations
		},
		decodeFunc: func(dec *json.Decoder, dst *shapeBuffer) error {
			return decodeAbsShapeIDSliceTo(dec, "operations", &dst.operations)
		},
	},
	"resources": {
		name:  "resources",
		types: []ShapeType{ServiceType, ResourceType},
		storeFunc: func(t ShapeType, src *shapeBuffer, dst *Shape) {
			var o *[]AbsShapeIDNode
			switch t {
			case ServiceType:
				o = &dst.service().Resources
			case ResourceType:
				o = &dst.resource().Resources
			}
			*o = src.resources
		},
		decodeFunc: func(dec *json.Decoder, dst *shapeBuffer) error {
			return decodeAbsShapeIDSliceTo(dec, "resources", &dst.resources)
		},
	},
	"errors": {
		name:  "errors",
		types: []ShapeType{ServiceType, OperationType},
		storeFunc: func(t ShapeType, src *shapeBuffer, dst *Shape) {
			var o *[]AbsShapeIDNode
			switch t {
			case ServiceType:
				o = &dst.service().Errors
			case OperationType:
				o = &dst.operation().Errors
			}
			*o = src.errors
		},
		decodeFunc: func(dec *json.Decoder, dst *shapeBuffer) error {
			return decodeAbsShapeIDSliceTo(dec, "errors", &dst.errors)
		},
	},
	"rename": {
		name:  "rename",
		types: []ShapeType{ServiceType},
		storeFunc: func(_ ShapeType, src *shapeBuffer, dst *Shape) {
			dst.service().Rename = src.rename
		},
		decodeFunc: func(dec *json.Decoder, dst *shapeBuffer) error {
			return decodeToMap(dec, "rename", &dst.rename)
		},
	},
	"identifiers": {
		name:  "identifiers",
		types: []ShapeType{ResourceType},
		storeFunc: func(_ ShapeType, src *shapeBuffer, dst *Shape) {
			dst.resource().Identifiers = src.identifiers
		},
		decodeFunc: func(dec *json.Decoder, dst *shapeBuffer) error {
			return decodeToMap(dec, "identifiers", &dst.identifiers)
		},
	},
	"create": {
		name:  "create",
		types: []ShapeType{ResourceType},
		storeFunc: func(_ ShapeType, src *shapeBuffer, dst *Shape) {
			dst.resource().Create = src.create
		},
		decodeFunc: func(dec *json.Decoder, dst *shapeBuffer) error {
			return decodeAbsShapeIDNodeTo(dec, &dst.create)
		},
	},
	"put": {
		name:  "put",
		types: []ShapeType{ResourceType},
		storeFunc: func(_ ShapeType, src *shapeBuffer, dst *Shape) {
			dst.resource().Put = src.put
		},
		decodeFunc: func(dec *json.Decoder, dst *shapeBuffer) error {
			return decodeAbsShapeIDNodeTo(dec, &dst.put)
		},
	},
	"read": {
		name:  "read",
		types: []ShapeType{ResourceType},
		storeFunc: func(_ ShapeType, src *shapeBuffer, dst *Shape) {
			dst.resource().Read = src.read
		},
		decodeFunc: func(dec *json.Decoder, dst *shapeBuffer) error {
			return decodeAbsShapeIDNodeTo(dec, &dst.read)
		},
	},
	"update": {
		name:  "update",
		types: []ShapeType{ResourceType},
		storeFunc: func(_ ShapeType, src *shapeBuffer, dst *Shape) {
			dst.resource().Update = src.update
		},
		decodeFunc: func(dec *json.Decoder, dst *shapeBuffer) error {
			return decodeAbsShapeIDNodeTo(dec, &dst.update)
		},
	},
	"delete": {
		name:  "delete",
		types: []ShapeType{ResourceType},
		storeFunc: func(_ ShapeType, src *shapeBuffer, dst *Shape) {
			dst.resource().Delete = src.delete
		},
		decodeFunc: func(dec *json.Decoder, dst *shapeBuffer) error {
			return decodeAbsShapeIDNodeTo(dec, &dst.delete)
		},
	},
	"list": {
		name:  "list",
		types: []ShapeType{ResourceType},
		storeFunc: func(_ ShapeType, src *shapeBuffer, dst *Shape) {
			dst.resource().List = src.list
		},
		decodeFunc: func(dec *json.Decoder, dst *shapeBuffer) error {
			return decodeAbsShapeIDNodeTo(dec, &dst.list)
		},
	},
	"collectionOperations": {
		name:  "collectionOperations",
		types: []ShapeType{ResourceType},
		storeFunc: func(_ ShapeType, src *shapeBuffer, dst *Shape) {
			dst.resource().CollectionOperations = src.collectionOperations
		},
		decodeFunc: func(dec *json.Decoder, dst *shapeBuffer) error {
			return decodeAbsShapeIDSliceTo(dec, "collection operations", &dst.collectionOperations)
		},
	},
	"input": {
		name:  "input",
		types: []ShapeType{OperationType},
		storeFunc: func(_ ShapeType, src *shapeBuffer, dst *Shape) {
			dst.operation().Input = src.input
		},
		decodeFunc: func(dec *json.Decoder, dst *shapeBuffer) error {
			return decodeAbsShapeIDNodeTo(dec, &dst.input)
		},
	},
	"output": {
		name:  "output",
		types: []ShapeType{OperationType},
		storeFunc: func(_ ShapeType, src *shapeBuffer, dst *Shape) {
			dst.operation().Output = src.output
		},
		decodeFunc: func(dec *json.Decoder, dst *shapeBuffer) error {
			return decodeAbsShapeIDNodeTo(dec, &dst.output)
		},
	},
}

type shapeBuffer struct {
	// Some fields can be multiple shape types.
	// Want a map shapetype[bool], because then when we find out the
	// actual type, we can just check and ensure it's in there for
	// each field.
	fields []shapeField

	key     *Member
	value   *Member
	members map[string]Member

	version     StringNode
	operations  []AbsShapeIDNode
	resources   []AbsShapeIDNode
	errors      []AbsShapeIDNode
	rename      map[AbsShapeID]StringNode
	identifiers map[string]AbsShapeIDNode

	create               *AbsShapeIDNode
	put                  *AbsShapeIDNode
	read                 *AbsShapeIDNode
	update               *AbsShapeIDNode
	delete               *AbsShapeIDNode
	list                 *AbsShapeIDNode
	collectionOperations []AbsShapeIDNode

	input  *AbsShapeIDNode
	output *AbsShapeIDNode
}
