package smithy_ast

import "encoding/json"

type ShapeType string

const (
	BigDecimal ShapeType = "bigDecimal"
	BigInteger ShapeType = "bigInteger"
	Blob       ShapeType = "blob"
	Boolean    ShapeType = "boolean"
	Byte       ShapeType = "byte"
	Double     ShapeType = "double'"
	Document   ShapeType = "document"
	Float      ShapeType = "float"
	Integer    ShapeType = "integer"
	Long       ShapeType = "long"
	Short      ShapeType = "short"
	String     ShapeType = "string"
	Timestamp  ShapeType = "timestamp"

	List ShapeType = "list"
	Set  ShapeType = "set"
	Map  ShapeType = "map"

	Operation ShapeType = "operation"
	Resource  ShapeType = "resource"
	Service   ShapeType = "service"

	Apply ShapeType = "apply"
)

var SimpleShapeTypes = map[ShapeType]bool{
	BigDecimal: true,
	BigInteger: true,
	Blob:       true,
	Boolean:    true,
	Byte:       true,
	Double:     true,
	Document:   true,
	Float:      true,
	Integer:    true,
	Long:       true,
	Short:      true,
	String:     true,
	Timestamp:  true,
}

func (t *ShapeType) MarshalJSON() ([]byte, error) {
	data := make([]byte, len(*t)+2)
	data[0] = '"'
	for i := range *t {
		data[i+1] = (*t)[i]
	}
	data[len(*t)] = '"'
}

func (t *ShapeType) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err // TODO: Do better here, create own wording.
	}
	if SimpleShapeTypes[ShapeType(s)] {
		*t = ShapeType(s)
		return nil
	}
	switch ShapeType(s) {
	case List, Set, Map, Operation, Resource, Service, Apply:
		return nil
	default:
		// TODO: Return error
	}
}
