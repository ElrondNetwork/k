package main

import (
	"fmt"
	"strings"
)

// K ... Defines a K entity, this is either a KItem, or a KSequence of KItems
type K interface {
	prettyTreePrint(indent int) string
}

type KItem interface {
}

type KApply struct {
	label string
	list  []K
}

type InjectedKLabel struct {
	label string
}

type KToken struct {
	value string
	sort  string
}

type KVariable struct {
	name string
}

// KSequence ... a sequence of K items
type KSequence struct {
	ks []K
}

func addIndent(sb *strings.Builder, indent int) {
	for i := 0; i < indent; i++ {
		sb.WriteString("    ")
	}
}

func (k KApply) prettyTreePrint(indent int) string {
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
			sb.WriteString(childk.prettyTreePrint(indent + 1))
		}
		sb.WriteRune('\n')
		addIndent(&sb, indent)
		sb.WriteRune('}')
	}

	return sb.String()
}

func (k InjectedKLabel) prettyTreePrint(indent int) string {
	var sb strings.Builder
	addIndent(&sb, indent)
	sb.WriteString(fmt.Sprintf("InjectedKLabel {label:%s}", k.label))
	return sb.String()
}

func (k KToken) prettyTreePrint(indent int) string {
	var sb strings.Builder
	addIndent(&sb, indent)
	sb.WriteString(fmt.Sprintf("KToken {value:%s, sort:%s}", k.value, k.sort))
	return sb.String()
}

func (k KVariable) prettyTreePrint(indent int) string {
	var sb strings.Builder
	addIndent(&sb, indent)
	sb.WriteString(fmt.Sprintf("KVariable {name:%s}", k.name))
	return sb.String()
}

func (k KSequence) prettyTreePrint(indent int) string {
	var sb strings.Builder
	addIndent(&sb, indent)
	sb.WriteString("KSequence {")
	if len(k.ks) == 0 {
		sb.WriteString(" <empty> }")
	} else {
		for i, childk := range k.ks {
			sb.WriteString("\n")
			sb.WriteString(childk.prettyTreePrint(indent + 1))
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
