package ast

import (
	"bytes"
	"encoding/json"
	"math/big"
	"reflect"
)

const (
	EnumTraitID        AbsShapeID = "smithy.api#enum"
	IDRefTraitID       AbsShapeID = "smithy.api#idRef"
	LengthTraitID      AbsShapeID = "smithy.api#length"
	PatternTraitID     AbsShapeID = "smithy.api#pattern"
	PrivateTraitID     AbsShapeID = "smithy.api#private"
	RangeTraitID       AbsShapeID = "smithy.api#range"
	RequiredTraitID    AbsShapeID = "smithy.api#required"
	UniqueItemsTraitID AbsShapeID = "smithy.api#uniqueItems"

	DeprecatedTraitID            AbsShapeID = "smithy.api#deprecated"
	DocumentationTraitID         AbsShapeID = "smithy.api#documentation"
	ExamplesTraitID              AbsShapeID = "smithy.api#examples"
	ExternalDocumentationTraitID AbsShapeID = "smithy.api#externalDocumentation"
	InternalTraitID              AbsShapeID = "smithy.api#internal"
	RecommendedTraitID           AbsShapeID = "smithy.api#recommended"
	SensitiveTraitID             AbsShapeID = "smithy.api#sensitive"
	SinceTraitID                 AbsShapeID = "smithy.api#since"
	TagsTraitID                  AbsShapeID = "smithy.api#tags"
	TitleTraitID                 AbsShapeID = "smithy.api#title"
	UnstableTraitID              AbsShapeID = "smithy.api#unstable"

	BoxTraitID    AbsShapeID = "smithy.api#box"
	ErrorTraitID  AbsShapeID = "smithy.api#error"
	InputTraitID  AbsShapeID = "smithy.api#input"
	OutputTraitID AbsShapeID = "smithy.api#output"
	SparseTraitID AbsShapeID = "smithy.api#sparse"

	ProtocolDefinitionTraitID AbsShapeID = "smithy.api#protocolDefinition"
	JSONNameTraitID           AbsShapeID = "smithy.api#jsonName"
	MediaTypeTraitID          AbsShapeID = "smithy.api#mediaType"
	TimestampFormatTraitID    AbsShapeID = "smithy.api#timestampFormat"

	AuthDefinitionTraitID AbsShapeID = "smithy.api#authDefinition"
	HTTPBasicAuthTraitID  AbsShapeID = "smithy.api#httpBasicAuth"
	HTTPDigestAuthTraitID AbsShapeID = "smithy.api#httpDigestAuth"
	HTTPBearerAuthTraitID AbsShapeID = "smithy.api#httpBearerAuth"
	HTTPAPIKeyAuthTraitID AbsShapeID = "smithy.api#httpApiKeyAuth"
	OptionalAuthTraitID   AbsShapeID = "smithy.api#optionalAuth"
	AuthTraitID           AbsShapeID = "smithy.api#auth"

	IdempotencyTokenTraitID     AbsShapeID = "smithy.api#idempotencyToken"
	IdempotentTraitID           AbsShapeID = "smithy.api#idempotent"
	ReadOnlyTraitID             AbsShapeID = "smithy.api#readonly"
	RetryableTraitID            AbsShapeID = "smithy.api#retryable"
	PaginatedTraitID            AbsShapeID = "smithy.api#paginated"
	HTTPChecksumRequiredTraitID AbsShapeID = "smithy.api#httpChecksumRequired"

	NoReplaceTraitID          AbsShapeID = "smithy.api#noReplace"
	ReferencesTraitID         AbsShapeID = "smithy.api#references"
	ResourceIdentifierTraitID AbsShapeID = "smithy.api#resourceIdentifier"

	StreamingTraitID      AbsShapeID = "smithy.api#streaming"
	RequiresLengthTraitID AbsShapeID = "smithy.api#requiresLength"

	HTTPTraitID                AbsShapeID = "smithy.api#http"
	HTTPErrorTraitID           AbsShapeID = "smithy.api#httpError"
	HTTPHeaderTraitID          AbsShapeID = "smithy.api#httpHeader"
	HTTPLabelTraitID           AbsShapeID = "smithy.api#httpLabel"
	HTTPPayloadTraitID         AbsShapeID = "smithy.api#httpPayload"
	HTTPPrefixedHeadersTraitID AbsShapeID = "smithy.api#httpPrefixedHeaders"
	HTTPQueryTraitID           AbsShapeID = "smithy.api#httpQuery"
	HTTPQueryParamsTraitID     AbsShapeID = "smithy.api#httpQueryParams"
	HTTPResponseCodeTraitID    AbsShapeID = "smithy.api#httpResponseCode" + ""
	CORSTraitID                AbsShapeID = "smithy.api#cors"

	XMLAttributeTraitID AbsShapeID = "smithy.api#xmlAttribute"
	XMLFlattenedTraitID AbsShapeID = "smithy.api#xmlFlattened"
	XMLNameTraitID      AbsShapeID = "smithy.api#xmlName"
	XMLNamespaceTraitID AbsShapeID = "smithy.api#xmlNamespace"

	EndpointTraitID  AbsShapeID = "smithy.api#endpoint"
	HostLabelTraitID AbsShapeID = "smithy.api#hostLabel"
)

type Traits map[AbsShapeID]interface{}

func (t *Traits) UnmarshalJSON(data []byte) error {
	dec := json.NewDecoder(bytes.NewReader(data))
	return t.decode(dec)
}

func (t *Traits) decode(dec *json.Decoder) error {
	t2 := make(Traits)
	err := decodeObject(dec, "traits map", func(dec2 *json.Decoder, key string, keyOffset int64) error {
		// Determine the type of the value to decode.
		var val interface{}
		if tp, ok := builtinTraits[AbsShapeID(key)]; ok {
			val = reflect.New(tp)
		}

		// Decode the value.
		err2 := dec.Decode(&val)
		if err2 != nil {
			return err2
		}
		t2[AbsShapeID(key)] = val
		return nil
	})
	if err != nil {
		return err
	}
	*t = t2
	return nil
}

type AnnotationTrait struct {
}

type EnumTrait []EnumTraitItem

type EnumTraitItem struct {
	Value         string
	Name          string
	Documentation string
	Tags          []string
	Deprecated    bool
}

type IDRefTrait struct {
	FailWhenMissing bool
	Selector        string
	ErrorMessage    string
}

type LengthTrait struct {
	Min int64
	Max int64
}

type RangeTrait struct {
	Min big.Rat
	Max big.Rat
}

type DeprecatedTrait struct {
	Message string
	Since   string
}

type ExamplesTrait []ExamplesTraitItem

type ExamplesTraitItem struct {
	Title         string
	Documentation string
	Input         map[string]interface{}
	Output        map[string]interface{}
	Error         *ExamplesTraitError
}

type ExamplesTraitError struct {
	ShapeID AbsShapeID
	Content map[string]interface{}
}

type RecommendedTrait struct {
	Reason string
}

type ErrorTrait string

func (e *ErrorTrait) UnmarshalJSON(data []byte) error {
	// Get a decoder on the data.
	dec := json.NewDecoder(bytes.NewReader(data))

	tok, _ := dec.Token()
	var s string
	var ok bool
	if s, ok = tok.(string); !ok {
		return newError("expected string value for trait " + string(ErrorTraitID))
	}

	switch s {
	case "client", "server":
		return nil
	default:
		return newErrorf("expected value for trait %s to be %q or %q (it is %q)", ErrorTraitID, "client", "server", s)
	}
}

type ProtocolDefinitionTrait struct {
	Traits                  []AbsShapeID `json:"traits"`
	NoInlineDocumentSupport bool         `json:"noInlineDocumentSupport"`
}

type AuthDefinitionTrait struct {
	Traits []AbsShapeID `json:"traits"`
}

type PaginatedTrait struct {
	InputToken  string `json:"inputToken"`
	OutputToken string `json:"outputToken"`
	Items       string `json:"items"`
	PageSize    string `json:"pageSize"`
}

type ReferencesTrait struct {
	Service  AbsShapeID        `json:"service"`
	Resource AbsShapeID        `json:"resource"`
	IDs      map[string]string `json:"ids"`
	Rel      string            `json:"rel"`
}

type HTTPTrait struct {
	Method string `json:"method"`
	URI    string `json:"uri"`
	Code   uint16 `json:"code"`
}

type CORSTrait struct {
	Origin                   string   `json:"origin"`
	MaxAge                   int32    `json:"maxAge"`
	AdditionalAllowedHeaders []string `json:"additionalAllowedHeaders"`
	AdditionalExposedHeaders []string `json:"additionalExposedHeaders"`
}

type XMLNamespaceTrait struct {
	URI    string `json:"uri"`
	Prefix string `json:"prefix"`
}

type EndpointTrait struct {
	HostPrefix string `json:"hostPrefix"`
}

var builtinTraits = map[AbsShapeID]reflect.Type{
	EnumTraitID:        reflect.TypeOf(EnumTrait{}),
	IDRefTraitID:       reflect.TypeOf(IDRefTrait{}),
	LengthTraitID:      reflect.TypeOf(LengthTrait{}),
	PatternTraitID:     reflect.TypeOf(""),
	PrivateTraitID:     reflect.TypeOf(AnnotationTrait{}),
	RangeTraitID:       reflect.TypeOf(RangeTrait{}),
	RequiredTraitID:    reflect.TypeOf(AnnotationTrait{}),
	UniqueItemsTraitID: reflect.TypeOf(AnnotationTrait{}),

	DeprecatedTraitID:            reflect.TypeOf(DeprecatedTrait{}),
	DocumentationTraitID:         reflect.TypeOf(""),
	ExamplesTraitID:              reflect.TypeOf(ExamplesTrait{}),
	ExternalDocumentationTraitID: reflect.TypeOf(map[string]string{}),
	InternalTraitID:              reflect.TypeOf(AnnotationTrait{}),
	RecommendedTraitID:           reflect.TypeOf(RecommendedTrait{}),
	SensitiveTraitID:             reflect.TypeOf(AnnotationTrait{}),
	SinceTraitID:                 reflect.TypeOf(""),
	TagsTraitID:                  reflect.TypeOf([]string{}),
	TitleTraitID:                 reflect.TypeOf(""),
	UnstableTraitID:              reflect.TypeOf(AnnotationTrait{}),

	BoxTraitID:    reflect.TypeOf(AnnotationTrait{}),
	ErrorTraitID:  reflect.TypeOf(AnnotationTrait{}),
	InputTraitID:  reflect.TypeOf(AnnotationTrait{}),
	OutputTraitID: reflect.TypeOf(AnnotationTrait{}),
	SparseTraitID: reflect.TypeOf(AnnotationTrait{}),

	ProtocolDefinitionTraitID: reflect.TypeOf(ProtocolDefinitionTrait{}),
	JSONNameTraitID:           reflect.TypeOf(""),
	MediaTypeTraitID:          reflect.TypeOf(""),
	TimestampFormatTraitID:    reflect.TypeOf(""),

	AuthDefinitionTraitID: reflect.TypeOf(AuthDefinitionTrait{}),
	HTTPBasicAuthTraitID:  reflect.TypeOf(AnnotationTrait{}),
	HTTPDigestAuthTraitID: reflect.TypeOf(AnnotationTrait{}),
	HTTPBearerAuthTraitID: reflect.TypeOf(AnnotationTrait{}),
	HTTPAPIKeyAuthTraitID: reflect.TypeOf(AnnotationTrait{}),
	OptionalAuthTraitID:   reflect.TypeOf(AnnotationTrait{}),
	AuthTraitID:           reflect.TypeOf([]string{}),

	IdempotencyTokenTraitID:     reflect.TypeOf(AnnotationTrait{}),
	IdempotentTraitID:           reflect.TypeOf(AnnotationTrait{}),
	ReadOnlyTraitID:             reflect.TypeOf(AnnotationTrait{}),
	RetryableTraitID:            reflect.TypeOf(AnnotationTrait{}),
	PaginatedTraitID:            reflect.TypeOf(PaginatedTrait{}),
	HTTPChecksumRequiredTraitID: reflect.TypeOf(AnnotationTrait{}),

	NoReplaceTraitID:          reflect.TypeOf(AnnotationTrait{}),
	ReferencesTraitID:         reflect.TypeOf(ReferencesTrait{}),
	ResourceIdentifierTraitID: reflect.TypeOf(""),

	StreamingTraitID:      reflect.TypeOf(AnnotationTrait{}),
	RequiresLengthTraitID: reflect.TypeOf(AnnotationTrait{}),

	HTTPTraitID:                reflect.TypeOf(HTTPTrait{}),
	HTTPErrorTraitID:           reflect.TypeOf(uint16(0)),
	HTTPHeaderTraitID:          reflect.TypeOf(""),
	HTTPLabelTraitID:           reflect.TypeOf(AnnotationTrait{}),
	HTTPPayloadTraitID:         reflect.TypeOf(AnnotationTrait{}),
	HTTPPrefixedHeadersTraitID: reflect.TypeOf(""),
	HTTPQueryTraitID:           reflect.TypeOf(""),
	HTTPQueryParamsTraitID:     reflect.TypeOf(AnnotationTrait{}),
	HTTPResponseCodeTraitID:    reflect.TypeOf(AnnotationTrait{}),
	CORSTraitID:                reflect.TypeOf(CORSTrait{}),

	XMLAttributeTraitID: reflect.TypeOf(AnnotationTrait{}),
	XMLFlattenedTraitID: reflect.TypeOf(AnnotationTrait{}),
	XMLNameTraitID:      reflect.TypeOf(""),
	XMLNamespaceTraitID: reflect.TypeOf(XMLNamespaceTrait{}),

	EndpointTraitID:  reflect.TypeOf(EndpointTrait{}),
	HostLabelTraitID: reflect.TypeOf(AnnotationTrait{}),
}
