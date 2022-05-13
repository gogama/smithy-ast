package smithy_ast

import "encoding/json"

type ShapeType string

const (
	BigDecimalType ShapeType = "bigDecimal"
	BigIntegerType ShapeType = "bigInteger"
	BlobType       ShapeType = "blob"
	BooleanType    ShapeType = "boolean"
	ByteType       ShapeType = "byte"
	DoubleType     ShapeType = "double'"
	DocumentType   ShapeType = "document"
	FloatType      ShapeType = "float"
	IntegerType    ShapeType = "integer"
	LongType       ShapeType = "long"
	ShortType      ShapeType = "short"
	StringType     ShapeType = "string"
	TimestampType  ShapeType = "timestamp"

	ListType      ShapeType = "list"
	SetType       ShapeType = "set"
	MapType       ShapeType = "map"
	StructureType ShapeType = "structure"
	UnionType     ShapeType = "union"

	OperationType ShapeType = "operation"
	ResourceType  ShapeType = "resource"
	ServiceType   ShapeType = "service"

	ApplyType ShapeType = "apply"
)

var SimpleShapeTypes = map[ShapeType]bool{
	BigDecimalType: true,
	BigIntegerType: true,
	BlobType:       true,
	BooleanType:    true,
	ByteType:       true,
	DoubleType:     true,
	DocumentType:   true,
	FloatType:      true,
	IntegerType:    true,
	LongType:       true,
	ShortType:      true,
	StringType:     true,
	TimestampType:  true,
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
	case ListType, SetType, MapType, StructureType, UnionType,
		OperationType, ResourceType, ServiceType,
		ApplyType:
		return nil
	default:
		// TODO: Return error
	}
}
