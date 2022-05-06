package smithy_ast

type AbsShapeID string

func (id *AbsShapeID) Namespace() string {

}

func (id *AbsShapeID) Name() Identifier {

}

func (id *AbsShapeID) MarshalJSON() ([]byte, error) {

}

func (id *AbsShapeID) UnmarshalJSON([]byte) error {

}

type AbsShapeOrMemberID string

func (id *AbsShapeID) Namespace() string {

}

func (id *AbsShapeID) Name() string {

}

func (id *AbsShapeID) Member() string {

}

func (id *AbsShapeOrMemberID) MarshalJSON() ([]byte, error) {

}

func (id *AbsShapeOrMemberID) UnmarshalJSON([]byte) error {

}
