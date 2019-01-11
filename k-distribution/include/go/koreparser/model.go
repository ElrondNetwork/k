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
type KSequence struct {
	ks []K
}

/*
kitem =
            | Map of sort * klabel * m
            | List of sort * klabel * t list
            | Set of sort * klabel * s
            | Array of sort * t * t Dynarray.t
            | Int of Z.t
            | Float of Gmp.FR.t * int * int
            | String of string
            | Bytes of bytes
            | StringBuffer of Buffer.t
            | Bool of bool
            | MInt of int * Z.t
            | ThreadLocal
            | Thread of t * t * t * t
            | Bottom
            | KApply0 of klabel
            | KApply1 of klabel * t
            | KApply2 of klabel * t * t
            | KApply3 of klabel * t * t * t
            | KApply4 of klabel * t * t * t * t
*/

// KItem ...
type KItem interface {
}

// KApply ... a type of KItem, TODO: document
type KApply struct {
	label string
	list  []K
}

// InjectedKLabel ... a type of KItem, TODO: document
type InjectedKLabel struct {
	label string
}

// KToken ... a type of KItem, TODO: document
type KToken struct {
	value string
	sort  string
}

// KVariable ... a type of KItem, TODO: document
type KVariable struct {
	name string
}

// Bottom ... a type of KItem, TODO: document
type Bottom struct {
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
	sb.WriteString(k.label)
	sb.WriteString(", list:")
	if len(k.list) == 0 {
		sb.WriteString("[] }")
	} else {
		for _, childk := range k.list {
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
	return simplePrint(indent, fmt.Sprintf("InjectedKLabel {label:%s}", k.label))
}

func (k KToken) PrettyTreePrint(indent int) string {
	return simplePrint(indent, fmt.Sprintf("KToken {value:%s, sort:%s}", k.value, k.sort))
}

func (k KVariable) PrettyTreePrint(indent int) string {
	return simplePrint(indent, fmt.Sprintf("KVariable {name:%s}", k.name))
}

func (k Bottom) PrettyTreePrint(indent int) string {
	return simplePrint(indent, "Bottom")
}

func (k KSequence) PrettyTreePrint(indent int) string {
	var sb strings.Builder
	addIndent(&sb, indent)
	sb.WriteString("KSequence {")
	if len(k.ks) == 0 {
		sb.WriteString(" <empty> }")
	} else {
		for i, childk := range k.ks {
			sb.WriteString("\n")
			sb.WriteString(childk.PrettyTreePrint(indent + 1))
			if i < len(k.ks)-1 {
				sb.WriteString(" ~>")
			}
		}
		sb.WriteRune('\n')
		addIndent(&sb, indent)
		sb.WriteRune('}')
	}

	return sb.String()
}
