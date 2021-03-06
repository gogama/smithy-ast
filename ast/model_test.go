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
		{
			name: "version, metadata, and shapes",
			json: `{"version":"1.2","metadata":{"number":123,"object":{"array":[]}},"shapes":{"test#List":{"type":"list","traits":{"smithy.api#length":{"max":5}},"member":{"target":"smithy.api#String","traits":{"smithy.api#required":{}}}}}}`,
			model: Model{
				Version: StringNode{Value: "1.2"},
				Metadata: map[string]InterfaceNode{
					"number": {Value: float64(123)},
					"object": {Value: map[string]interface{}{"array": []interface{}{}}},
				},
				Shapes: map[AbsShapeID]Shape{
					"test#List": {
						Type: ListType,
						Traits: Traits{
							LengthTraitID: &LengthTrait{
								Max: &Int64Node{Value: 5},
							},
						},
						Value: &Member{
							Target: AbsShapeIDNode{Value: "smithy.api#String"},
							Traits: Traits{
								RequiredTraitID: &AnnotationTrait{},
							},
						},
					},
				},
			},
		},
		{
			name: "error/unsupported key",
			json: `{"foo":"bar"}`,
			err:  jsonError(`unsupported key "foo" in model`, 1),
		},
		{
			name: "error/not an object",
			json: `   true   `,
			err:  jsonError("expected '{' to start model", 0),
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
