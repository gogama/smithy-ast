package smithy_ast

type Shape struct {
	node
	Type             ShapeType
	Traits           Traits
	ListOrSet        *Member
	MapKey           *Member
	MapValue         *Member
	Operation        *OperationShape
	Resource         *ResourceShape
	Service          *ServiceShape
	StructureOrUnion map[string]Member
}

func (s *Shape) MarshalJSON() ([]byte, error) {

}

func (s *Shape) UnmarshalJSON([]byte) error {

}

type OperationShape struct {
	Input  *AbsShapeID  `json:"input,omitempty"`
	Output *AbsShapeID  `json:"output"`
	Errors []AbsShapeID `json:"errors"`
}

type ResourceShape struct {
	Identifiers          map[string]AbsShapeID `json:"identifiers"`
	Create               *AbsShapeID           `json:"create"`
	Put                  *AbsShapeID           `json:"put"`
	Read                 *AbsShapeID           `json:"read"`
	Update               *AbsShapeID           `json:"update"`
	Delete               *AbsShapeID           `json:"delete"`
	List                 *AbsShapeID           `json:"list"`
	Operations           []AbsShapeID          `json:"operations"`
	CollectionOperations []AbsShapeID          `json:"collectionOperations"`
	Resources            []AbsShapeID          `json:"resources"`
}

type ServiceShape struct {
	Version    string                    `json:"version"`
	Operations []AbsShapeID              `json:"operations"`
	Resources  []AbsShapeID              `json:"resources"`
	Errors     []AbsShapeID              `json:"errors"`
	Rename     map[AbsShapeID]Identifier `json:"rename"`
}
