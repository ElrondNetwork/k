package %PACKAGE_MODEL%

import (
	"fmt"
	"math/big"
	"strings"
)

// K ... Defines a K entity
type K interface {
	PrettyTreePrint(indent int) string
}

// KSequence ... a sequence of K items
type KSequence struct {
	Ks []K
}

// EmptyKSequence ... the KSequence with no elements
var EmptyKSequence = KSequence{Ks: nil}

// KItem ...
type KItem interface {
}

// KApply ... a type of KItem, TODO: document
type KApply struct {
	Label KLabel
	List  []K
}

// InjectedKLabel ... a type of KItem, TODO: document
type InjectedKLabel struct {
	Label KLabel
}

// KToken ... a type of KItem, TODO: document
type KToken struct {
	Value string
	Sort  Sort
}

// KVariable ... a type of KItem, TODO: document
type KVariable struct {
	Name string
}

// Map ... a type of KItem, TODO: document
type Map struct {
	Sort  Sort
	Label KLabel
	Data  map[K]K
}

// Set ... a type of KItem, TODO: document
type Set struct {
	Sort  Sort
	Label KLabel
	Data  map[K]bool
}

// List ... a type of KItem, TODO: document
type List struct {
	Sort  Sort
	Label KLabel
	Data  []K
}

// Array ... array of K Items of fixed size
type Array struct {
	Sort    Sort
	Data    []*K
	Default *K
}

// Int ... integer type, implemented via a big int
type Int struct {
	Value *big.Int
}

// IntZero ... K Int with value zero
var IntZero = Int{Value: big.NewInt(0)}

// MInt ... machine integer
type MInt int32

// Float ... a type of KItem, TODO: document
type Float float32

// String ... a type of KItem, TODO: document
type String string

// StringBuffer ... a type of KItem, TODO: document
type StringBuffer struct {
}

// Bytes ... a type of KItem, TODO: document
type Bytes []byte

// Bool ... K boolean value
type Bool bool

// BoolTrue ... K boolean value with value true
var BoolTrue = Bool(true)

// BoolFalse ... K boolean value with value false
var BoolFalse = Bool(false)

// Bottom ... a K item that contains no data
type Bottom struct {
}

// InternedBottom ... usually used as a dummy object
var InternedBottom K = Bottom{}

// NoResult ... what to return when a function returns an error
var NoResult K = Bottom{}

func addIndent(sb *strings.Builder, indent int) {
	for i := 0; i < indent; i++ {
		sb.WriteString("    ")
	}
}

func simplePrint(indent int, str string) string {
	var sb strings.Builder
	addIndent(&sb, indent)
	sb.WriteString(str)
	return sb.String()
}

// PrettyTreePrint ... A tree representation of a KApply object
func (k KApply) PrettyTreePrint(indent int) string {
	var sb strings.Builder
	addIndent(&sb, indent)
	sb.WriteString("KApply {Label:")
	sb.WriteString(k.Label.Name())
	sb.WriteString(", List:")
	if len(k.List) == 0 {
		sb.WriteString("[] }")
	} else {
		for _, childk := range k.List {
			sb.WriteRune('\n')
			sb.WriteString(childk.PrettyTreePrint(indent + 1))
		}
		sb.WriteRune('\n')
		addIndent(&sb, indent)
		sb.WriteRune('}')
	}

	return sb.String()
}

// PrettyTreePrint ... A tree representation of a InjectedKLabel object
func (k InjectedKLabel) PrettyTreePrint(indent int) string {
	return simplePrint(indent, fmt.Sprintf("InjectedKLabel {Label:%s}", k.Label.Name()))
}

// PrettyTreePrint ... A tree representation of a KToken object
func (k KToken) PrettyTreePrint(indent int) string {
	return simplePrint(indent, fmt.Sprintf("KToken {Sort:%s, Value:%s}", k.Sort.Name(), k.Value))
}

// PrettyTreePrint ... A tree representation of a KVariable object
func (k KVariable) PrettyTreePrint(indent int) string {
	return simplePrint(indent, fmt.Sprintf("KVariable {Name:%s}", k.Name))
}

// PrettyTreePrint ... A tree representation of a Map object
func (k Map) PrettyTreePrint(indent int) string {
	var sb strings.Builder
	addIndent(&sb, indent)
	sb.WriteString("Map {Sort:")
	sb.WriteString(k.Sort.Name())
	sb.WriteString(", Label:")
	sb.WriteString(k.Label.Name())
	sb.WriteString(", Data:")
	if len(k.Data) == 0 {
		sb.WriteString(" <empty> }")
	} else {
		for k, v := range k.Data {
			sb.WriteString("\n")
			addIndent(&sb, indent+1)
			sb.WriteString("key: ")
			sb.WriteString(k.PrettyTreePrint(0))
			sb.WriteString("  value: ")
			sb.WriteString(v.PrettyTreePrint(0))
		}
		sb.WriteRune('\n')
		addIndent(&sb, indent)
		sb.WriteRune('}')
	}

	return sb.String()
}

// PrettyTreePrint ... A tree representation of a List object
func (k List) PrettyTreePrint(indent int) string {
	return simplePrint(indent, fmt.Sprintf("List {Sort:%s, Label:%s}", k.Sort.Name(), k.Label.Name()))
}

// PrettyTreePrint ... A tree representation of a Set object
func (k Set) PrettyTreePrint(indent int) string {
	return simplePrint(indent, fmt.Sprintf("Set {Sort:%s, Label:%s}", k.Sort.Name(), k.Label.Name()))
}

// PrettyTreePrint ... A tree representation of a KApply object
func (k Array) PrettyTreePrint(indent int) string {
	return simplePrint(indent, fmt.Sprintf("Array {Sort:%s}", k.Sort.Name()))
}

// NewInt ... provides new Int instance
func NewInt(bi *big.Int) Int {
	return Int{Value: bi}
}

// NewIntFromInt ... provides new Int instance
func NewIntFromInt(x int) Int {
	return NewIntFromInt64(int64(x))
}

// NewIntFromInt64 ... provides new Int instance
func NewIntFromInt64(x int64) Int {
	return Int{Value: big.NewInt(x)}
}

// ParseInt ... creates K int from string representation
func ParseInt(s string) (Int, error) {
	b := big.NewInt(0)
	if s != "0" {
		b.UnmarshalText([]byte(s))
		if b.Cmp(IntZero.Value) == 0 {
			return IntZero, &parseIntError{parseVal: s}
		}
	}
	return NewInt(b), nil
}

// NewIntFromString ... same as ParseInt but panics instead of error
func NewIntFromString(s string) Int {
    i, err := ParseInt(s)
    if err != nil {
        panic(err)
    }
	return i
}

// PrettyTreePrint ... A tree representation of a KApply object
func (k Int) PrettyTreePrint(indent int) string {
	return simplePrint(indent, fmt.Sprintf("Int (%s)", k.Value.String()))
}

// PrettyTreePrint ... A tree representation of a KApply object
func (k MInt) PrettyTreePrint(indent int) string {
	return simplePrint(indent, fmt.Sprintf("MInt (%d)", k))
}

// PrettyTreePrint ... A tree representation of a KApply object
func (k Float) PrettyTreePrint(indent int) string {
	return simplePrint(indent, fmt.Sprintf("Float (%f)", k))
}

// PrettyTreePrint ... A tree representation of a KApply object
func (k String) PrettyTreePrint(indent int) string {
	return simplePrint(indent, fmt.Sprintf("String (%s)", k))
}

// PrettyTreePrint ... A tree representation of a KApply object
func (k StringBuffer) PrettyTreePrint(indent int) string {
	return simplePrint(indent, "StringBuffer [not yet implemented]")
}

// PrettyTreePrint ... A tree representation of a KApply object
func (k Bytes) PrettyTreePrint(indent int) string {
	return simplePrint(indent, fmt.Sprintf("Bytes (%b)", k))
}

// PrettyTreePrint ... A tree representation of a KApply object
func (k Bool) PrettyTreePrint(indent int) string {
	return simplePrint(indent, fmt.Sprintf("Bool (%t)", k))
}

// PrettyTreePrint ... A tree representation of a KApply object
func (k Bottom) PrettyTreePrint(indent int) string {
	return simplePrint(indent, "Bottom")
}

// PrettyTreePrint ... A tree representation of a KApply object
func (k KSequence) PrettyTreePrint(indent int) string {
	var sb strings.Builder
	addIndent(&sb, indent)
	sb.WriteString("KSequence {")
	if len(k.Ks) == 0 {
		sb.WriteString(" <empty> }")
	} else {
		for i, childk := range k.Ks {
			sb.WriteString("\n")
			sb.WriteString(childk.PrettyTreePrint(indent + 1))
			if i < len(k.Ks)-1 {
				sb.WriteString(" ~>")
			}
		}
		sb.WriteRune('\n')
		addIndent(&sb, indent)
		sb.WriteRune('}')
	}

	return sb.String()
}

// IsEmpty ... returns true if KSequence has no elements
func (k KSequence) IsEmpty() bool {
	return len(k.Ks) == 0
}
