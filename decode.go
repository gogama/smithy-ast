package smithy_ast

import "encoding/json"

type valueDecoder func(dec *json.Decoder, key string) error

func decodeObject(dec *json.Decoder, name string, valDec valueDecoder) error {
	var tok json.Token
	var err error

	// Expect an open brace starting a JSON object.
	tok, err = dec.Token()
	var delim json.Delim
	var ok bool
	if delim, ok = tok.(json.Delim); !ok || delim != '{' {
		return newError("expected { to start " + name)
	}

	// Find all key, value pairs in the object.
	for dec.More() {
		// Get the key.
		tok, err = dec.Token()
		var key string
		if key, ok = tok.(string); !ok {
			return newError("expected trait shape ID (key) within traits map")
		}

		// Decode the value.
		err = valDec(dec, key)
		if err != nil {
			return err
		}
	}

	// Expect a closing brace ending the JSON object.
	tok, err = dec.Token()
	if delim, ok = tok.(json.Delim); !ok || delim != '}' {
		return newError("expected } to end " + name)
	}

	// Object parsed successfully.
	return nil
}
