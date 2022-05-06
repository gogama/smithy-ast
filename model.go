package smithy_ast

import "io"

type Model struct {
	Version  string                       `json:"version"`
	Metadata map[string]interface{}       `json:"metadata"`
	Shapes   map[AbsShapeOrMemberID]Shape `json:"shapes"`
}

func ReadModel(r io.Reader) (Model, error) {

}

func WriteModel(m *Model, w io.Writer) error {

}
