package smithy_ast

import (
	"strings"
)

type AbsShapeID string

func (id *AbsShapeID) Namespace() string {
	s := string(*id)
	i := strings.IndexByte(s, '#')
	if i < 0 {
		return ""
	}
	return s[0:i]
}

func (id *AbsShapeID) Name() string {
	s := string(*id)
	i := strings.IndexByte(s, '#')
	s = s[i+1:]
	i = strings.IndexByte(s, '$')
	if i < 0 {
		return ""
	}
	return s[0:i]
}

func (id *AbsShapeID) Member() string {
	s := string(*id)
	i := strings.IndexByte(s, '#')
	s = s[i+1:]
	i = strings.IndexByte(s, '$')
	return s[i+1:]
}
