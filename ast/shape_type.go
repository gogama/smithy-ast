package ast

type ShapeType string

const (
	BigDecimalType ShapeType = "bigDecimal"
	BigIntegerType ShapeType = "bigInteger"
	BlobType       ShapeType = "blob"
	BooleanType    ShapeType = "boolean"
	ByteType       ShapeType = "byte"
	DoubleType     ShapeType = "double"
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

	ServiceType   ShapeType = "service"
	ResourceType  ShapeType = "resource"
	OperationType ShapeType = "operation"

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

var ShapeTypes = map[ShapeType]bool{
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
	ListType:       true,
	SetType:        true,
	MapType:        true,
	StructureType:  true,
	UnionType:      true,
	ServiceType:    true,
	ResourceType:   true,
	OperationType:  true,
	ApplyType:      true,
}
