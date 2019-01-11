package koreparser

import (
	"fmt"
	"strings"
)

// K ... Defines a K entity, this is either a KItem, or a KSequence of KItems
type K interface {
	PrettyTreePrint(indent int) string
}

// KSequence ... a sequence of K items
type KSequence []K

// KItem ...
type KItem interface {
}

// KApply ... a type of KItem, TODO: document
type KApply struct {
	Label string
	List  []K
}

// InjectedKLabel ... a type of KItem, TODO: document
type InjectedKLabel struct {
	Label string
}

// KToken ... a type of KItem, TODO: document
type KToken struct {
	Value string
	Sort  string
}

// KVariable ... a type of KItem, TODO: document
type KVariable struct {
	Name string
}

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

func (k KApply) PrettyTreePrint(indent int) string {
	var sb strings.Builder
	addIndent(&sb, indent)
	sb.WriteString("KApply {label:")
	sb.WriteString(k.Label)
	sb.WriteString(", list:")
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

func (k InjectedKLabel) PrettyTreePrint(indent int) string {
	return simplePrint(indent, fmt.Sprintf("InjectedKLabel {label:%s}", k.Label))
}

func (k KToken) PrettyTreePrint(indent int) string {
	return simplePrint(indent, fmt.Sprintf("KToken {value:%s, sort:%s}", k.Value, k.Sort))
}

func (k KVariable) PrettyTreePrint(indent int) string {
	return simplePrint(indent, fmt.Sprintf("KVariable {name:%s}", k.Name))
}

func (k KSequence) PrettyTreePrint(indent int) string {
	var sb strings.Builder
	addIndent(&sb, indent)
	sb.WriteString("KSequence {")
	if len(k) == 0 {
		sb.WriteString(" <empty> }")
	} else {
		for i, childk := range k {
			sb.WriteString("\n")
			sb.WriteString(childk.PrettyTreePrint(indent + 1))
			if i < len(k)-1 {
				sb.WriteString(" ~>")
			}
		}
		sb.WriteRune('\n')
		addIndent(&sb, indent)
		sb.WriteRune('}')
	}

	return sb.String()
}
