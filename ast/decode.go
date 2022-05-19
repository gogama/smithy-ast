package ast

import (
	"encoding/json"
	"reflect"
	"strconv"
	"strings"
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
		return jsonError("expected '{' to start "+name, offset)
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
			return jsonError("expected string key within "+name, offset)
		}

		// Check for duplication.
		if seen[key] {
			return jsonError("duplicate key "+strconv.Quote(key)+" within "+name, offset)
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
		return jsonError("expected '}' to end "+name, dec.InputOffset())
	}

	// Object parsed successfully.
	return nil
}

// decodeToMap decodes a JSON object into a map. The map key type must
// be a kind of string, and the value type must be type whose pointer
// type implements Node. The map must be non-nil.
func decodeToMap(dec *json.Decoder, name string, target interface{}) error {
	v := reflect.ValueOf(target)
	t := v.Type()
	kt := t.Key()
	if kt.Kind() != reflect.String {
		panic(newErrorf("map key type must be a kind of string within %s map but %s is not", name, kt.Name()))
	}
	vt := t.Elem()
	var n Node
	nt := reflect.TypeOf(&n).Elem()
	if !reflect.PtrTo(vt).Implements(nt) {
		panic(newErrorf("map value type must implement Node within %s map but %s does not", name, vt.Name()))
	}

	return decodeObject(dec, name, func(dec2 *json.Decoder, key string, _ int64) error {
		vv := reflect.New(vt)
		n2 := vv.Interface().(Node)
		err := n2.Decode(dec2)
		if err != nil {
			return err
		}
		v.SetMapIndex(reflect.ValueOf(key).Convert(kt), vv.Elem())
		return nil
	})
}

// decodeToStructPtr decodes a JSON object into a pointer to a struct.
// Each struct field must have a "json" tags to specify the JSON key
// corresponding to the field. Each struct field must either implement
// Node or be a map whose keys are strings and whose values implement
// Node.
func decodeToStructPtr(dec *json.Decoder, name string, target interface{}) error {
	v := reflect.ValueOf(target)
	t := v.Type()

	if t.Kind() != reflect.Pointer {
		panic(newError("pointer to struct required"))
	}

	v = v.Elem()
	t = v.Type()

	if t.Kind() != reflect.Struct {
		panic(newError("pointer to struct required"))
	}

	m := t.NumField()
	fields := make(map[string]reflect.StructField, m)
	for i := 0; i < m; i++ {
		f := t.Field(i)
		key := f.Tag.Get("json")
		x := strings.IndexByte(key, ',')
		if x >= 0 {
			key = key[0:x]
		}
		if key == "" {
			panic(newErrorf("field %s [%i] in struct %s has no usable json tag", f.Name, i, t.Name()))
		}
		if _, ok := fields[key]; ok {
			panic(newErrorf("field %s [%i] in struct %s duplicates JSON key %q", f.Name, i, t.Name(), key))
		}
		fields[key] = f
	}

	var n Node
	nt := reflect.TypeOf(&n)

	return decodeObject(dec, name, func(dec2 *json.Decoder, key string, keyOffset int64) error {
		f, ok := fields[key]
		if !ok {
			return unsupportedKeyError(name, key, keyOffset)
		}

		fv := v.FieldByIndex(f.Index)
		ft := fv.Type()
		if ft.Implements(nt) {
			return fv.Interface().(Node).Decode(dec2)
		} else if ft.Kind() == reflect.Map {
			fv.Set(reflect.MakeMap(ft))
			return decodeToMap(dec2, name+`["`+key+`"]`, fv.Interface())
		} else if ft.Kind() == reflect.Slice {
			return decodeToSlicePtr(dec2, name+`["`+key+`"]`, fv.Addr().Interface())
		} else {
			panic(newErrorf("field %s in struct %s has invalid type", f.Name, t.Name()))
		}
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
		return jsonError("expected '[' to start "+name, offset)
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
		return jsonError("expected ']' to end "+name, dec.InputOffset())
	}

	// Array parsed successfully.
	return nil
}

// decodeToSlice decodes a JSON array into a pointer to a slice. The
// slice element type must implement Node.
func decodeToSlicePtr(dec *json.Decoder, name string, target interface{}) error {
	v := reflect.ValueOf(target)
	t := v.Type()

	if t.Kind() != reflect.Pointer {
		panic(newError("pointer to slice required"))
	}

	v2 := v.Elem()
	t2 := v.Type()

	if t2.Kind() != reflect.Slice {
		panic(newError("pointer to slice required"))
	}

	var n Node
	nt := reflect.TypeOf(&n)
	et := t2.Elem()
	if !et.Implements(nt) {
		panic(newError("slice element type must implement Node"))
	}

	return decodeArray(dec, name, func(dec2 *json.Decoder, index int) error {
		ev := reflect.New(et)
		err := ev.Interface().(Node).Decode(dec2)
		if err != nil {
			return err
		}
		v.Set(reflect.Append(v2, ev))
		return nil
	})
}

type numberDecoder func(s string) error

// decodeNumber decodes the JSON number at the current position in the
// decoder. It expects and reads a syntactically valid JSON number from
// the decoder as a string, then passes the string to a callback
// function to decode it. If the decoded number needs to be persisted,
// this is the responsibility of the callback.
func decodeNumber(dec *json.Decoder, numDec numberDecoder) error {
	offset := dec.InputOffset()
	dec.UseNumber()
	t, err := dec.Token()
	if err != nil {
		return err
	}
	var n json.Number
	var ok bool
	if n, ok = t.(json.Number); !ok {
		return jsonError("expected number", offset)
	}
	return numDec(string(n))
}
