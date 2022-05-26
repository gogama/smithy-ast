package ast

import (
	"encoding/json"
	"io"
)

// Model represents the root of a Smithy model's abstract syntax tree in
// JSON AST format. A Model is a Node.
type Model struct {
	node
	Version  StringNode               `json:"version"`
	Metadata map[string]InterfaceNode `json:"metadata,omitempty"`
	Shapes   map[AbsShapeID]Shape     `json:"shapes,omitempty"`
}

func (m *Model) Decode(dec *json.Decoder) error {
	offset := dec.InputOffset()
	version := false

	err := decodeObject(dec, "model", func(dec2 *json.Decoder, key string, keyOffset int64) error {
		switch key {
		case "version":
			version = true
			return m.Version.Decode(dec2)
		case "metadata":
			m.Metadata = make(map[string]InterfaceNode)
			return decodeToMap(dec2, "metadata", m.Metadata)
		case "shapes":
			m.Shapes = make(map[AbsShapeID]Shape)
			return decodeToMap(dec2, "shapes", m.Shapes)
		default:
			return unsupportedKeyError("model", key, keyOffset)
		}
	})

	if err != nil {
		return err
	}

	if !version {
		return jsonError("missing version", offset)
	}

	return nil
}

func (m *Model) UnmarshalJSON(data []byte) error {
	return unmarshalJSON(data, m)
}

// ReadModel reads a Model from an io.Reader. The reader must "contain"
// valid JSON which is a valid JSON AST.
//
// If an error occurs because of a problem with the input JSON, the
// returned error has type *JSONError. Other errors may also be
// returned, e.g. for input/output errors with the reader.
func ReadModel(r io.Reader) (m Model, err error) {
	dec := json.NewDecoder(r)
	err = m.Decode(dec)
	return
}

// WriteModel writes a Model to an io.Writer.
func WriteModel(m Model, w io.Writer) error {
	enc := json.NewEncoder(w)
	return enc.Encode(&m)
}

// MergeModels merges the given models together following the Smithy
// model merging rules: https://awslabs.github.io/smithy/1.0/spec/core/model.html#merging-models.
//
// If the merge succeeds without conflicts, the returned error is nil.
// Otherwise, the returned error has type MergeConflictsError and
// contains all conflicts discovered during the merge process as
// sub-errors. The returned model always represents a best effort at
// merging: even if conflicts were encountered, it will contain all
// elements that were able to be successfully merged.
func MergeModels(m ...Model) (Model, error) {
	if len(m) == 0 {
		panic(newErrorf("no models to merge"))
	}

	var r Model
	var err MergeConflictsError

	for i := range m {
		err = append(err, mergeVersions(&r, &m[i])...)
		err = append(err, mergeMetadata(&r, &m[i])...)
		err = append(err, mergeShapes(&r, &m[i])...)
	}

	if len(err) > 0 {
		return r, err
	}

	return r, nil
}

func mergeVersions(dst, src *Model) []MergeConflictError {
	return nil // TODO
}

func mergeMetadata(dst, src *Model) []MergeConflictError {
	return nil // TODO
}

func mergeShapes(dst, src *Model) []MergeConflictError {
	return nil // TODO
}
