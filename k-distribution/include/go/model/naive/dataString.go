%COMMENT%

package %PACKAGE%

import (
	"strings"
)

// String is a KObject that contains a string
type String struct {
	Value string
}

// Bytes is a KObject that contains a slice of bytes
type Bytes struct {
	Value []byte
}

// StringBuffer is a KObject that contains a string buffer
type StringBuffer struct {
	Value strings.Builder
}

// StringEmpty is a reference to an empty string
var StringEmpty = &String{Value: ""}

// BytesEmpty is a reference to a Bytes item with no bytes (length 0)
var BytesEmpty = &Bytes{Value: nil}

// IsString returns true if reference points to a string
func IsString(ref KReference) bool {
	_, typeOk := ref.(*String)
	return typeOk
}

// IsBytes returns true if reference points to a byte array
func IsBytes(ref KReference) bool {
	_, typeOk := ref.(*Bytes)
	return typeOk
}

// IsStringBuffer returns true if reference points to a string buffer
func IsStringBuffer(ref KReference) bool {
	_, typeOk := ref.(*StringBuffer)
	return typeOk
}

// GetStringObject yields the cast object for a String reference, if possible.
// Deprecated
func (ms *ModelState) GetStringObject(ref KReference) (*String, bool) {
	castObj, typeOk := ref.(*String)
	if !typeOk {
		return nil, false
	}
	return castObj, true
}

// GetString converts reference to a Go string, if possbile
func (ms *ModelState) GetString(ref KReference) (string, bool) {
	castObj, typeOk := ms.GetStringObject(ref)
	if !typeOk {
		return "", false
	}
	return castObj.Value, true
}

// GetBytesObject yields the cast object for a Bytes reference, if possible.
func (ms *ModelState) GetBytesObject(ref KReference) (*Bytes, bool) {
	castObj, typeOk := ref.(*Bytes)
	if !typeOk {
		return nil, false
	}
	return castObj, true
}

// GetStringBufferObject yields the cast object for a StringBuffer reference, if possible.
func (ms *ModelState) GetStringBufferObject(ref KReference) (*StringBuffer, bool) {
	castObj, typeOk := ref.(*StringBuffer)
	if !typeOk {
		return nil, false
	}
	return castObj, true
}

// NewString creates a new K string object from a Go string.
func (ms *ModelState) NewString(str string) KReference {
	return &String{Value: str}
}

// NewStringConstant creates a new K string object from a Go string.
func NewStringConstant(str string) KReference {
	return &String{Value: str}
}

// NewBytes creates a new K string object from a Go string
func (ms *ModelState) NewBytes(value []byte) KReference {
	return &Bytes{Value: value}
}

// NewStringBuffer creates a new object and returns the reference.
func (ms *ModelState) NewStringBuffer() KReference {
	return &StringBuffer{Value: strings.Builder{}}
}

// IsEmpty returns true if Bytes is the empty byte slice
func (k *Bytes) IsEmpty() bool {
	return len(k.Value) == 0
}

// String yields a Go string representation of the K String
func (k *String) String() string {
	return k.Value
}

// IsEmpty returns true if it is the empty string
func (k *String) IsEmpty() bool {
	return len(k.Value) == 0
}
