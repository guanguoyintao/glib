package ujson

import (
	"io"
)

// uJsonAPI is a binding of specific config.
// This interface is inspired by github.com/json-iterator/go,
// and has same behaviors under equavilent config.
type uJsonAPI interface {
	// MarshalToString returns the JSON encoding string of v
	MarshalToString(v interface{}) (string, error)
	// Marshal returns the JSON encoding bytes of v.
	Marshal(v interface{}) ([]byte, error)
	// MarshalIndent returns the JSON encoding bytes with indent and prefix.
	MarshalIndent(v interface{}, prefix, indent string) ([]byte, error)
	// UnmarshalFromString parses the JSON-encoded bytes and stores the result in the value pointed to by v.
	UnmarshalFromString(str string, v interface{}) error
	// Unmarshal parses the JSON-encoded string and stores the result in the value pointed to by v.
	Unmarshal(data []byte, v interface{}) error
	// NewEncoder create a Encoder holding writer
	NewEncoder(writer io.Writer) Encoder
	// NewDecoder create a Decoder holding reader
	NewDecoder(reader io.Reader) Decoder
	// Valid validates the JSON-encoded bytes and reportes if it is valid
	Valid(data []byte) bool
}

// Encoder encodes JSON into io.Writer
type Encoder interface {
	// Encode writes the JSON encoding of v to the stream, followed by a newline character.
	Encode(val interface{}) error
	// SetEscapeHTML specifies whether problematic HTML characters
	// should be escaped inside JSON quoted strings.
	// The default behavior NOT ESCAPE
	SetEscapeHTML(on bool)
	// SetIndent instructs the encoder to format each subsequent encoded value
	// as if indented by the package-level function Indent(dst, src, prefix, indent).
	// Calling SetIndent("", "") disables indentation
	SetIndent(prefix, indent string)
}

// Decoder decodes JSON from io.Read
type Decoder interface {
	// Decode reads the next JSON-encoded value from its input and stores it in the value pointed to by v.
	Decode(val interface{}) error
	// Buffered returns a reader of the data remaining in the Decoder's buffer.
	// The reader is valid until the next call to Decode.
	Buffered() io.Reader
	// DisallowUnknownFields causes the Decoder to return an error when the destination is a struct
	// and the input contains object keys which do not match any non-ignored, exported fields in the destination.
	DisallowUnknownFields()
	// More reports whether there is another element in the current array or object being parsed.
	More() bool
	// UseNumber causes the Decoder to unmarshal a number into an interface{} as a Number instead of as a float64.
	UseNumber()
}
