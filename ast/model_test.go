package ast

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestModel(t *testing.T) {
	testCases := []struct {
		name  string
		json  string
		model Model
		err   error
	}{
		{
			name: "missing version",
			json: `{}`,
			err:  jsonError("missing version", 0),
		},
		{
			name:  "only version",
			json:  `{"version":"foo"}`,
			model: Model{Version: StringNode{Value: "foo"}},
		},
		{
			name: "version and metadata",
			json: `{"version":"1.0","metadata":{"foo":["bar",{"baz":"qux"}]}}`,
			model: Model{
				Version: StringNode{Value: "1.0"},
				Metadata: map[string]InterfaceNode{
					"foo": {Value: []interface{}{"bar", map[string]interface{}{"baz": "qux"}}},
				},
			},
		},
		{
			name: "version and shapes",
			json: `{"version":"1.1","shapes":{"foo#Bar":{"type":"string"}}}`,
			model: Model{
				Version: StringNode{Value: "1.1"},
				Shapes: map[AbsShapeID]Shape{
					"foo#Bar": {
						Type: StringType,
					},
				},
			},
		},
	}

	validateRead := func(t *testing.T, expectedErr, err error, expectedModel, model Model) {
		if expectedErr == nil {
			assert.NoError(t, err)
			assert.Equal(t, expectedModel, model)
		} else {
			assert.EqualError(t, err, expectedErr.Error())
			assert.ErrorIs(t, err, expectedErr)
		}
	}

	// Iterative test cases.
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Run("Decode", func(t *testing.T) {
				var m Model
				dec := json.NewDecoder(strings.NewReader(testCase.json))

				err := m.Decode(dec)

				validateRead(t, testCase.err, err, testCase.model, m)
			})

			t.Run("UnmarshalJSON", func(t *testing.T) {
				var m Model

				err := json.Unmarshal([]byte(testCase.json), &m)

				validateRead(t, testCase.err, err, testCase.model, m)
			})

			t.Run("ReadModel", func(t *testing.T) {
				r := strings.NewReader(testCase.json)

				m, err := ReadModel(r)

				validateRead(t, testCase.err, err, testCase.model, m)
			})

			t.Run("WriteModel", func(t *testing.T) {
				if testCase.err != nil {
					t.Skip()
				}

				w := bytes.Buffer{}

				err := WriteModel(testCase.model, &w)

				assert.NoError(t, err)
				assert.Equal(t, testCase.json, strings.TrimRight(w.String(), "\n"))
			})
		})
	}

	// Static test cases.
	t.Run("WriteModel.Error", func(t *testing.T) {
		// TODO. Cursory test case using a mock writer that errors out.
	})
}
