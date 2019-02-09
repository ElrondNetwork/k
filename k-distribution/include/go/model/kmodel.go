package %PACKAGE_MODEL%

import (
	"fmt"
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

// Array ... a type of KItem, TODO: document
type Array struct {
	Sort  Sort
	Label KLabel
	Data  []K
}

// Int ... a type of KItem, TODO: document
type Int int32

// MInt ... a type of KItem, TODO: document
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

// Bool ... a type of KItem, TODO: document
type Bool bool

// Bottom ... a type of KItem, TODO: document
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

// PrettyTreePrint ... A tree representation of a KApply object
func (k InjectedKLabel) PrettyTreePrint(indent int) string {
	return simplePrint(indent, fmt.Sprintf("InjectedKLabel {Label:%s}", k.Label.Name()))
}

// PrettyTreePrint ... A tree representation of a KApply object
func (k KToken) PrettyTreePrint(indent int) string {
	return simplePrint(indent, fmt.Sprintf("KToken {Sort:%s, Value:%s}", k.Sort.Name(), k.Value))
}

// PrettyTreePrint ... A tree representation of a KApply object
func (k KVariable) PrettyTreePrint(indent int) string {
	return simplePrint(indent, fmt.Sprintf("KVariable {Name:%s}", k.Name))
}

// PrettyTreePrint ... A tree representation of a KApply object
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
			addIndent(&sb, indent + 1)
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

// PrettyTreePrint ... A tree representation of a KApply object
func (k List) PrettyTreePrint(indent int) string {
	return simplePrint(indent, fmt.Sprintf("List {Sort:%s, Label:%s}", k.Sort.Name(), k.Label.Name()))
}

// PrettyTreePrint ... A tree representation of a KApply object
func (k Set) PrettyTreePrint(indent int) string {
	return simplePrint(indent, fmt.Sprintf("Set {Sort:%s, Label:%s}", k.Sort.Name(), k.Label.Name()))
}

// PrettyTreePrint ... A tree representation of a KApply object
func (k Array) PrettyTreePrint(indent int) string {
	return simplePrint(indent, fmt.Sprintf("Array {Sort:%s, Label:%s}", k.Sort.Name(), k.Label.Name()))
}

// PrettyTreePrint ... A tree representation of a KApply object
func (k Int) PrettyTreePrint(indent int) string {
	return simplePrint(indent, fmt.Sprintf("Int (%d)", k))
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
