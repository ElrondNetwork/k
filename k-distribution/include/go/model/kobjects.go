%COMMENT%

package %PACKAGE_MODEL%

import (
	"math/big"
	"strings"
)

// KObject defines a K item object that is managed by the model
type KObject interface {
	referenceType() kreferenceType
	equals(ms *ModelState, other KObject) bool
	deepCopy(ms *ModelState) KObject
	prettyPrint(ms *ModelState, sb *strings.Builder, indent int)
	kprint(ms *ModelState, sb *strings.Builder)
	collectionsToK(ms *ModelState) KReference
}

// InjectedKLabel is a KObject representing an InjectedKLabel item in K
type InjectedKLabel struct {
	Label KLabel
}

func (*InjectedKLabel) referenceType() kreferenceType {
	return injectedKLabelRef
}

// KToken is a KObject representing a KToken item in K
type KToken struct {
	Value string
	Sort  Sort
}

func (*KToken) referenceType() kreferenceType {
	return ktokenRef
}

// KVariable is a KObject representing a KVariable item in K
type KVariable struct {
	Name string
}

func (*KVariable) referenceType() kreferenceType {
	return kvariableRef
}

// Map is a KObject representing a map in K
type Map struct {
	Sort  Sort
	Label KLabel
	Data  map[KMapKey]KReference
}

func (*Map) referenceType() kreferenceType {
	return mapRef
}

// Set is a KObject representing a set in K
type Set struct {
	Sort  Sort
	Label KLabel
	Data  map[KMapKey]bool
}

func (*Set) referenceType() kreferenceType {
	return setRef
}

// List is a KObject representing a list in K
type List struct {
	Sort  Sort
	Label KLabel
	Data  []KReference
}

func (*List) referenceType() kreferenceType {
	return listRef
}

// Array is a KObject holding an array that can grow
type Array struct {
	Sort Sort
	Data *DynamicArray
}

func (*Array) referenceType() kreferenceType {
	return arrayRef
}

// BigInt is a KObject representing a big int in K
type BigInt struct {
	Value *big.Int
}

func (*BigInt) referenceType() kreferenceType {
	return bigIntRef
}

// MInt is a KObject representing a machine integer in K
type MInt struct {
	Value int32
}

func (*MInt) referenceType() kreferenceType {
	return mintRef
}

// Float is a KObject representing a float in K
type Float struct {
	Value float32
}

func (*Float) referenceType() kreferenceType {
	return floatRef
}

// String is a KObject that contains a string
type String struct {
	Value string
}

func (*String) referenceType() kreferenceType {
	return stringRef
}

// StringBuffer is a KObject that contains a string buffer
type StringBuffer struct {
	Value strings.Builder
}

func (*StringBuffer) referenceType() kreferenceType {
	return stringBufferRef
}

// Bytes is a KObject that contains a slice of bytes
type Bytes struct {
	Value []byte
}

func (*Bytes) referenceType() kreferenceType {
	return bytesRef
}
