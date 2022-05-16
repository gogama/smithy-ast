package ast

import (
	"encoding/json"
	"reflect"
	"strconv"
)

type valueDecoder func(dec *json.Decoder, key string, keyOffset int64) error

// decodeObject decodes the JSON object at the current position in
// the decoder. It extracts keys and then calls a callback function
// to decode the value. If the key and value need to be persisted, this
// is the responsibility of the callback.
func decodeObject(dec *json.Decoder, name string, valDec valueDecoder) error {
	var tok json.Token
	var offset int64

	// Expect an open brace starting a JSON object.
	offset = dec.InputOffset()
	tok, _ = dec.Token()
	var delim json.Delim
	var ok bool
	if delim, ok = tok.(json.Delim); !ok || delim != '{' {
		return modelError("expected '{' to start "+name, offset)
	}

	// Record keys already seen.
	seen := make(map[string]bool)

	// Find all key, value pairs in the object.
	for dec.More() {
		// Get the key.
		offset = dec.InputOffset()
		tok, _ = dec.Token()
		var key string
		if key, ok = tok.(string); !ok {
			return modelError("expected string key within "+name, offset)
		}

		// Check for duplication.
		if seen[key] {
			return modelError("duplicate key "+strconv.Quote(key)+" within "+name, offset)
		}

		// Decode the value.
		err := valDec(dec, key, dec.InputOffset())
		if err != nil {
			return err
		}
	}

	// Expect a closing brace ending the JSON object.
	tok, _ = dec.Token()
	if delim, ok = tok.(json.Delim); !ok || delim != '}' {
		return modelError("expected '}' to end "+name, dec.InputOffset())
	}

	// Object parsed successfully.
	return nil
}

// decodeToMap decodes a JSON object into a map. The map key type must
// be a kind of string, and the value type must be a kind of Node. The
// map must be non-nil.
func decodeToMap(dec *json.Decoder, name string, target interface{}) error {
	v := reflect.ValueOf(target)
	t := v.Type()
	kt := t.Key()
	vt := t.Elem()

	return decodeObject(dec, name, func(dec2 *json.Decoder, key string, _ int64) error {
		kv := reflect.Zero(kt)
		kv.Set(reflect.ValueOf(key))
		vv := reflect.Zero(vt)
		n := vv.Interface().(Node)
		err := n.Decode(dec2)
		if err != nil {
			return err
		}
		v.SetMapIndex(reflect.ValueOf(key), vv)
		return nil
	})
}

type elementDecoder func(dec *json.Decoder, index int) error

// decodeArray decodes the JSON array at the current position in the
// decoder. For each element in the array, it calls a callback function
// to decode that element. If the decoded element needs to be persisted,
// this is the responsibility of the callback.
func decodeArray(dec *json.Decoder, name string, elemDec elementDecoder) error {
	var tok json.Token
	var offset int64

	// Expect an open bracket starting a JSON object.
	offset = dec.InputOffset()
	tok, _ = dec.Token()
	var delim json.Delim
	var ok bool
	if delim, ok = tok.(json.Delim); !ok || delim != '[' {
		return modelError("expected '[' to start "+name, offset)
	}

	// Decode each element in the array.
	for index := 0; dec.More(); index++ {
		err := elemDec(dec, index)
		if err != nil {
			return err
		}
	}

	// Expect a closing bracket ending the JSON array.
	tok, _ = dec.Token()
	if delim, ok = tok.(json.Delim); !ok || delim != ']' {
		return modelError("expected ']' to end "+name, dec.InputOffset())
	}

	// Array parsed successfully.
	return nil
}
