package ast

import (
	"bytes"
	"encoding/json"
	"reflect"
	"sort"
)

const (
	TraitTraitID       AbsShapeID = "smithy.api#trait"
	SuppressionTraitID AbsShapeID = "smithy.api#suppression"
	UnitTraitID        AbsShapeID = "smithy.api#unitShape"

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

type Traits map[AbsShapeID]Node

func (t *Traits) decode(dec *json.Decoder) error {
	t2 := make(Traits)
	err := decodeObject(dec, "traits map", func(dec2 *json.Decoder, key string, keyOffset int64) error {
		// Determine the type of the value to decode.
		var v reflect.Value
		if tp, ok := builtinTraits[AbsShapeID(key)]; ok {
			v = reflect.New(tp)
		} else {
			v = reflect.New(reflect.TypeOf(InterfaceNode{}))
		}

		// Decode the value.
		n := v.Interface().(Node)
		err2 := dec.Decode(n)
		if err2 != nil {
			return err2
		}
		t2[AbsShapeID(key)] = n
		return nil
	})
	if err != nil {
		return err
	}
	*t = t2
	return nil
}

func (t *Traits) UnmarshalJSON(data []byte) error {
	dec := json.NewDecoder(bytes.NewReader(data))
	return t.decode(dec)
}

func (t Traits) MarshalJSON() ([]byte, error) {
	keys := make([]string, 0, len(t))
	for key := range t {
		keys = append(keys, string(key))
	}
	sort.Strings(keys)
	var b bytes.Buffer
	b.WriteByte('{')
	enc := json.NewEncoder(&b)
	for _, key := range keys {
		err := enc.Encode(key)
		if err != nil {
			return nil, err
		}
		b.WriteByte(':')
		err = enc.Encode(t[AbsShapeID(key)])
		if err != nil {
			return nil, err
		}
	}
	b.WriteByte('}')
	return b.Bytes(), nil
}

type AnnotationTrait struct {
	node
}

func (n *AnnotationTrait) Decode(dec *json.Decoder) error {
	offset := dec.InputOffset()
	return decodeObject(dec, "annotation trait", func(_ *json.Decoder, _ string, _ int64) error {
		return jsonError("annotation trait must be an empty object", offset)
	})
}

func (n *AnnotationTrait) UnmarshalJSON(data []byte) error {
	return unmarshalJSON(data, n)
}

// TODO: document - https://awslabs.github.io/smithy/1.0/spec/core/model.html#traits
type TraitTrait struct {
	node
	Selector              *StringNode  `json:"selector,omitempty"`
	Conflicts             []StringNode `json:"conflicts,omitempty"`
	StructurallyExclusive *StringNode  `json:"structurallyExclusive,omitempty"`
}

func (n *TraitTrait) Decode(dec *json.Decoder) error {
	return decodeToStructPtr(dec, "trait trait", n)
}

func (n *TraitTrait) UnmarshalJSON(data []byte) error {
	return unmarshalJSON(data, n)
}

type SuppressionTrait struct {
	node
	Items []StringNode
}

func (n *SuppressionTrait) Decode(dec *json.Decoder) error {
	return decodeToSlicePtr(dec, "suppression trait", &n.Items)
}

func (n *SuppressionTrait) UnmarshalJSON(data []byte) error {
	return unmarshalJSON(data, n)
}

func (n SuppressionTrait) MarshalJSON() ([]byte, error) {
	return json.Marshal(n.Items)
}

type EnumTraitItem struct {
	node
	Value         StringNode   `json:"value"`
	Name          *StringNode  `json:"name,omitempty"`
	Documentation *StringNode  `json:"documentation,omitempty"`
	Tags          []StringNode `json:"tags,omitempty""`
	Deprecated    *BoolNode    `json:"deprecated,omitempty"`
}

func (n *EnumTraitItem) Decode(dec *json.Decoder) error {
	return decodeToStructPtr(dec, "enum trait item", n)
}

func (n *EnumTraitItem) UnmarshalJSON(data []byte) error {
	return unmarshalJSON(data, n)
}

type EnumTrait struct {
	node
	Items []EnumTraitItem
}

func (n *EnumTrait) Decode(dec *json.Decoder) error {
	return decodeToSlicePtr(dec, "enum trait", &n.Items)
}

func (n *EnumTrait) UnmarshalJSON(data []byte) error {
	return unmarshalJSON(data, n)
}

func (n EnumTrait) MarshalJSON() ([]byte, error) {
	return json.Marshal(n.Items)
}

type IDRefTrait struct {
	node
	FailWhenMissing BoolNode   `json:"failWhenMissing"`
	Selector        StringNode `json:"selector"`
	ErrorMessage    StringNode `json:"errorMessage"`
}

func (n *IDRefTrait) Decode(dec *json.Decoder) error {
	return decodeToStructPtr(dec, "idRef trait", n)
}

func (n *IDRefTrait) UnmarshalJSON(data []byte) error {
	return unmarshalJSON(data, n)
}

type LengthTrait struct {
	node
	Min *Int64Node `json:"min,omitempty"`
	Max *Int64Node `json:"max,omitempty"`
}

func (n *LengthTrait) Decode(dec *json.Decoder) error {
	return decodeToStructPtr(dec, "length trait", n)
}

func (n *LengthTrait) UnmarshalJSON(data []byte) error {
	return unmarshalJSON(data, n)
}

type RangeTrait struct {
	node
	Min BigFloatNode `json:"min,omitempty"`
	Max BigFloatNode `json:"max,omitempty"`
}

func (n *RangeTrait) Decode(dec *json.Decoder) error {
	return decodeToStructPtr(dec, "range trait", n)
}

func (n *RangeTrait) UnmarshalJSON(data []byte) error {
	return unmarshalJSON(data, n)
}

type DeprecatedTrait struct {
	node
	Message StringNode
	Since   StringNode
}

type ExamplesTrait struct {
	node
	Items []ExamplesTraitItem
}

type ExamplesTraitItem struct {
	Title         StringNode
	Documentation StringNode
	Input         map[string]InterfaceNode
	Output        map[string]InterfaceNode
	Error         *ExamplesTraitError
}

type ExamplesTraitError struct {
	ShapeID AbsShapeIDNode
	Content map[string]InterfaceNode
}

type RecommendedTrait struct {
	node
	Reason StringNode
}

type ProtocolDefinitionTrait struct {
	node
	Traits                  []AbsShapeIDNode `json:"traits"`
	NoInlineDocumentSupport BoolNode         `json:"noInlineDocumentSupport"`
}

type AuthDefinitionTrait struct {
	node
	Traits []AbsShapeIDNode `json:"traits"`
}

type PaginatedTrait struct {
	node
	InputToken  StringNode `json:"inputToken"`
	OutputToken StringNode `json:"outputToken"`
	Items       StringNode `json:"items"`
	PageSize    StringNode `json:"pageSize"`
}

type ReferencesTrait struct {
	node
	Service  AbsShapeIDNode        `json:"service"`
	Resource AbsShapeIDNode        `json:"resource"`
	IDs      map[string]StringNode `json:"ids"`
	Rel      StringNode            `json:"rel"`
}

type HTTPTrait struct {
	node
	Method StringNode `json:"method"`
	URI    StringNode `json:"uri"`
	Code   Int32Node  `json:"code"` // TODO: Nodify and check node types in Smithy, it is integer so int32 https://awslabs.github.io/smithy/1.0/spec/core/model.html#trait-node-values
}

type CORSTrait struct {
	node
	Origin                   StringNode   `json:"origin"`
	MaxAge                   Int32Node    `json:"maxAge"` // TODO: nodify
	AdditionalAllowedHeaders []StringNode `json:"additionalAllowedHeaders"`
	AdditionalExposedHeaders []StringNode `json:"additionalExposedHeaders"`
}

type XMLNamespaceTrait struct {
	node
	URI    StringNode `json:"uri"`
	Prefix StringNode `json:"prefix"`
}

type EndpointTrait struct {
	node
	HostPrefix StringNode `json:"hostPrefix"`
}

var builtinTraits = map[AbsShapeID]reflect.Type{
	TraitTraitID:       nil,
	UnitTraitID:        nil,
	SuppressionTraitID: nil,

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
