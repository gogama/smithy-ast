package smithy_ast

type Member struct {
	Target AbsShapeID                 `json:"target"`
	Traits map[AbsShapeID]interface{} `json:"traits"`
}
