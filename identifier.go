package smithy_ast

import "regexp"

type Identifier string

const identifierPattern = "[A-Za-z_][A-Za-z0-9_]*"

var identifierRegexp = regexp.MustCompile("^" + identifierPattern)

func matchIdentifier(s string) int {
	idx := identifierRegexp.FindStringIndex(s)
	if idx == nil {
		return -1
	}
	if idx[1] != len(s) {
		return -1
	}
	return 0
}

func matchNamespaceHash(s string) int {
	idx := identifierRegexp.FindStringIndex(s)
	if idx == nil {
		return -1
	}
	if idx[1] < len(s) && s[idx[1]] == '#' {
		return idx[1] + 1
	}
	return -1
}
