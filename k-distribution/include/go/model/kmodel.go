package %PACKAGE_MODEL%

import (
	"fmt"
	"math/big"
	"strings"
)

// K ... Defines a K entity
type K interface {
	Equals(K) bool
	PrettyTreePrint(indent int) string
}

// KSequence ... a sequence of K items
type KSequence struct {
	Ks []K
}

// EmptyKSequence ... the KSequence with no elements
var EmptyKSequence = &KSequence{Ks: nil}

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
	Data  map[KMapKey]K
}

// Set ... a type of KItem, TODO: document
type Set struct {
	Sort  Sort
	Label KLabel
	Data  map[KMapKey]bool
}

// List ... a type of KItem, TODO: document
type List struct {
	Sort  Sort
	Label KLabel
	Data  []K
}

// Array ... array of K Items of fixed size
type Array struct {
	Sort Sort
	Data *DynamicArray
}

// Int ... integer type, implemented via a big int
type Int struct {
	Value *big.Int
}

// IntZero ... K Int with value zero
var IntZero = &Int{Value: big.NewInt(0)}

// IntMinusOne ... K Int with value -1
var IntMinusOne = &Int{Value: big.NewInt(-1)}

// MInt ... machine integer
type MInt struct {
	Value int32
}

// Float ... a type of KItem, TODO: document
type Float struct {
	Value float32
}

// String ... a type of KItem, TODO: document
type String struct {
	Value string
}

// StringBuffer ... a type of KItem, TODO: document
type StringBuffer struct {
}

// Bytes ... a type of KItem, TODO: document
type Bytes struct {
	Value []byte
}

// Bool ... K boolean value
type Bool struct {
	Value bool
}

// BoolTrue ... K boolean value with value true
var BoolTrue = &Bool{Value: true}

// BoolFalse ... K boolean value with value false
var BoolFalse = &Bool{Value: false}

// Bottom ... a K item that contains no data
type Bottom struct {
}

// InternedBottom ... usually used as a dummy object
var InternedBottom = &Bottom{}

// NoResult ... what to return when a function returns an error
var NoResult = &Bottom{}

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

// Equals ... Deep comparison
func (k *KApply) Equals(arg K) bool {
	other, typeOk := arg.(*KApply)
	if !typeOk {
		return false
	}
	if k.Label != other.Label {
		return false
	}
	if len(k.List) != len(other.List) {
		return false
	}
	for i := 0; i < len(k.List); i++ {
		if !k.List[i].Equals(other.List[i]) {
			return false
		}
	}
	return true
}

// PrettyTreePrint ... A tree representation of a KApply object
func (k *KApply) PrettyTreePrint(indent int) string {
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

// Equals ... Deep comparison
func (k *InjectedKLabel) Equals(arg K) bool {
	other, typeOk := arg.(*InjectedKLabel)
	if !typeOk {
		return false
	}
	if k.Label != other.Label {
		return false
	}
	return true
}

// PrettyTreePrint ... A tree representation of a InjectedKLabel object
func (k *InjectedKLabel) PrettyTreePrint(indent int) string {
	return simplePrint(indent, fmt.Sprintf("InjectedKLabel {Label:%s}", k.Label.Name()))
}

// Equals ... Deep comparison
func (k *KToken) Equals(arg K) bool {
	other, typeOk := arg.(*KToken)
	if !typeOk {
		return false
	}
	if k.Sort != other.Sort {
		return false
	}
	return k.Value == other.Value
}

// PrettyTreePrint ... A tree representation of a KToken object
func (k *KToken) PrettyTreePrint(indent int) string {
	return simplePrint(indent, fmt.Sprintf("KToken {Sort:%s, Value:%s}", k.Sort.Name(), k.Value))
}

// Equals ... Deep comparison
func (k *KVariable) Equals(arg K) bool {
	other, typeOk := arg.(*KVariable)
	if !typeOk {
		return false
	}
	if k.Name != other.Name {
		return false
	}
	return true
}

// PrettyTreePrint ... A tree representation of a KVariable object
func (k *KVariable) PrettyTreePrint(indent int) string {
	return simplePrint(indent, fmt.Sprintf("KVariable {Name:%s}", k.Name))
}

// Equals ... Deep comparison
func (k *Map) Equals(arg K) bool {
	other, typeOk := arg.(*Map)
	if !typeOk {
		return false
	}
	if len(k.Data) != len(other.Data) {
		return false
	}
	for key, val := range k.Data {
		otherVal, found := other.Data[key]
		if !found {
			return false
		}
		if !val.Equals(otherVal) {
			return false
		}
	}
	return true
}

// PrettyTreePrint ... A tree representation of a Map object
func (k *Map) PrettyTreePrint(indent int) string {
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
			sb.WriteString(k.String())
			sb.WriteString("  value: ")
			sb.WriteString(v.PrettyTreePrint(0))
		}
		sb.WriteRune('\n')
		addIndent(&sb, indent)
		sb.WriteRune('}')
	}

	return sb.String()
}

// Equals ... Deep comparison
func (k *List) Equals(arg K) bool {
	other, typeOk := arg.(*List)
	if !typeOk {
		return false
	}
	if k.Sort != other.Sort {
		return false
	}
	if k.Label != other.Label {
		return false
	}
	if len(k.Data) != len(other.Data) {
		return false
	}
	for i := 0; i < len(k.Data); i++ {
		if !k.Data[i].Equals(other.Data[i]) {
			return false
		}
	}
	return true
}

// PrettyTreePrint ... A tree representation of a List object
func (k *List) PrettyTreePrint(indent int) string {
	return simplePrint(indent, fmt.Sprintf("List {Sort:%s, Label:%s}", k.Sort.Name(), k.Label.Name()))
}

// Equals ... Deep comparison
func (k *Set) Equals(arg K) bool {
	other, typeOk := arg.(*Set)
	if !typeOk {
		return false
	}
	if len(k.Data) != len(other.Data) {
		return false
	}
	for key := range k.Data {
		_, found := other.Data[key]
		if !found {
			return false
		}
	}
	return true
}

// PrettyTreePrint ... A tree representation of a Set object
func (k *Set) PrettyTreePrint(indent int) string {
	return simplePrint(indent, fmt.Sprintf("Set {Sort:%s, Label:%s}", k.Sort.Name(), k.Label.Name()))
}

// Equals ... Deep comparison
func (k *Array) Equals(arg K) bool {
	other, typeOk := arg.(*Array)
	if !typeOk {
		return false
	}
	if k.Sort != other.Sort {
		return false
	}
	return k.Data.Equals(other.Data)
}

// PrettyTreePrint ... A tree representation of a KApply object
func (k *Array) PrettyTreePrint(indent int) string {
	return simplePrint(indent, fmt.Sprintf("Array {Sort:%s}", k.Sort.Name()))
}

// NewInt ... provides new Int instance
func NewInt(bi *big.Int) *Int {
	return &Int{Value: bi}
}

// NewIntFromInt ... provides new Int instance
func NewIntFromInt(x int) *Int {
	return NewIntFromInt64(int64(x))
}

// NewIntFromInt64 ... provides new Int instance
func NewIntFromInt64(x int64) *Int {
	return &Int{Value: big.NewInt(x)}
}

// NewIntFromUint64 ... provides new Int instance
func NewIntFromUint64(x uint64) *Int {
	var z big.Int
	z.SetUint64(x)
	return &Int{Value: &z}
}

// ParseInt ... creates K int from string representation
func ParseInt(s string) (*Int, error) {
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
func NewIntFromString(s string) *Int {
	i, err := ParseInt(s)
	if err != nil {
		panic(err)
	}
	return i
}

// Equals ... Deep comparison
func (k *Int) Equals(arg K) bool {
	other, typeOk := arg.(*Int)
	if !typeOk {
		return false
	}
	return k.Value.Cmp(other.Value) == 0
}

// PrettyTreePrint ... A tree representation of an Int object
func (k *Int) PrettyTreePrint(indent int) string {
	return simplePrint(indent, fmt.Sprintf("Int (%s)", k.Value.String()))
}

// Equals ... Deep comparison
func (k *MInt) Equals(arg K) bool {
	other, typeOk := arg.(*MInt)
	if !typeOk {
		return false
	}
	return k.Value == other.Value
}

// PrettyTreePrint ... A tree representation of a KApply object
func (k *MInt) PrettyTreePrint(indent int) string {
	return simplePrint(indent, fmt.Sprintf("MInt (%d)", k.Value))
}

// Equals ... Deep comparison
func (k *Float) Equals(arg K) bool {
	other, typeOk := arg.(*Float)
	if !typeOk {
		return false
	}
	return k.Value == other.Value
}

// PrettyTreePrint ... A tree representation of a KApply object
func (k *Float) PrettyTreePrint(indent int) string {
	return simplePrint(indent, fmt.Sprintf("Float (%f)", k.Value))
}

// NewString ... Creates a new K string object from a Go string
func NewString(str string) *String {
	return &String{Value: str}
}

// String ... Yields a Go string representation of the K String
func (k *String) String() string {
	return k.Value
}

// Equals ... Deep comparison
func (k *String) Equals(arg K) bool {
	other, typeOk := arg.(*String)
	if !typeOk {
		return false
	}
	return k.Value == other.Value
}

// PrettyTreePrint ... A tree representation of a KApply object
func (k *String) PrettyTreePrint(indent int) string {
	return simplePrint(indent, fmt.Sprintf("String (%s)", k))
}

// Equals ... Deep comparison
func (k *StringBuffer) Equals(arg K) bool {
	panic("StringBuffer not yet implemented.")
}

// PrettyTreePrint ... A tree representation of a KApply object
func (k *StringBuffer) PrettyTreePrint(indent int) string {
	return simplePrint(indent, "StringBuffer [not yet implemented]")
}

// Equals ... Deep comparison
func (k *Bytes) Equals(arg K) bool {
	panic("Bytes not yet implemented.")
}

// PrettyTreePrint ... A tree representation of a KApply object
func (k *Bytes) PrettyTreePrint(indent int) string {
	return simplePrint(indent, fmt.Sprintf("Bytes (%b)", k))
}

// ToBool ... Convert Go bool to K Bool
func ToBool(b bool) *Bool {
	if b {
		return BoolTrue
	}
	return BoolFalse
}

// Equals ... Deep comparison
func (k *Bool) Equals(arg K) bool {
	other, typeOk := arg.(*Bool)
	if !typeOk {
		return false
	}
	return k.Value == other.Value
}

// PrettyTreePrint ... A tree representation of a KApply object
func (k *Bool) PrettyTreePrint(indent int) string {
	return simplePrint(indent, fmt.Sprintf("Bool (%t)", k.Value))
}

// Equals ... Deep comparison
func (k *Bottom) Equals(arg K) bool {
	_, typeOk := arg.(*Bottom)
	if !typeOk {
		return false
	}
	return true
}

// PrettyTreePrint ... A tree representation of a KApply object
func (k *Bottom) PrettyTreePrint(indent int) string {
	return simplePrint(indent, "Bottom")
}

// Equals ... Deep comparison
func (k *KSequence) Equals(arg K) bool {
	other, typeOk := arg.(*KSequence)
	if !typeOk {
		return false
	}
	if len(k.Ks) != len(other.Ks) {
		return false
	}
	for i := 0; i < len(k.Ks); i++ {
		if !k.Ks[i].Equals(other.Ks[i]) {
			return false
		}
	}
	return true
}

// PrettyTreePrint ... A tree representation of a KApply object
func (k *KSequence) PrettyTreePrint(indent int) string {
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
func (k *KSequence) IsEmpty() bool {
	return len(k.Ks) == 0
}
