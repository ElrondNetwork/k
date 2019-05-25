%COMMENT%

package %PACKAGE_MODEL%

import (
	"fmt"
	"strings"
)

// StructPrint ... returns a representation of a K item that somewhat resembles a Go declaration
func StructPrint(k K) string {
	var sb strings.Builder
	k.structPrint(&sb, 0)
	return sb.String()
}

func simplePrint(sb *strings.Builder, indent int, str string) {
	addIndent(sb, indent)
	sb.WriteString(str)
}

func (k *KApply) structPrint(sb *strings.Builder, indent int) {
	addIndent(sb, indent)
	sb.WriteString("KApply {Label:")
	sb.WriteString(k.Label.Name())
	sb.WriteString(", List:")
	if len(k.List) == 0 {
		sb.WriteString("[] }")
	} else {
		for _, childk := range k.List {
			sb.WriteRune('\n')
			childk.structPrint(sb, indent+1)
		}
		sb.WriteRune('\n')
		addIndent(sb, indent)
		sb.WriteRune('}')
	}
}

func (k *InjectedKLabel) structPrint(sb *strings.Builder, indent int) {
	simplePrint(sb, indent, fmt.Sprintf("InjectedKLabel {Label:%s}", k.Label.Name()))
}

func (k *KToken) structPrint(sb *strings.Builder, indent int) {
	simplePrint(sb, indent, fmt.Sprintf("KToken {Sort:%s, Value:%s}", k.Sort.Name(), k.Value))
}

func (k *KVariable) structPrint(sb *strings.Builder, indent int) {
	simplePrint(sb, indent, fmt.Sprintf("KVariable {Name:%s}", k.Name))
}

func (k *Map) structPrint(sb *strings.Builder, indent int) {
	addIndent(sb, indent)
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
			addIndent(sb, indent+1)
			sb.WriteString("key: ")
			sb.WriteString(k.String())
			sb.WriteString("  value: ")
			v.structPrint(sb, 0)
		}
		sb.WriteRune('\n')
		addIndent(sb, indent)
		sb.WriteRune('}')
	}
}

func (k *List) structPrint(sb *strings.Builder, indent int) {
	// TODO: print data
	simplePrint(sb, indent, fmt.Sprintf("List {Sort:%s, Label:%s}", k.Sort.Name(), k.Label.Name()))
}

func (k *Set) structPrint(sb *strings.Builder, indent int) {
	// TODO: print data
	simplePrint(sb, indent, fmt.Sprintf("Set {Sort:%s, Label:%s}", k.Sort.Name(), k.Label.Name()))
}

func (k *Array) structPrint(sb *strings.Builder, indent int) {
	// TODO: print data
	simplePrint(sb, indent, fmt.Sprintf("Array {Sort:%s}", k.Sort.Name()))
}

func (k *Int) structPrint(sb *strings.Builder, indent int) {
	simplePrint(sb, indent, fmt.Sprintf("Int (%s)", k.Value.String()))
}

func (k *MInt) structPrint(sb *strings.Builder, indent int) {
	simplePrint(sb, indent, fmt.Sprintf("MInt (%d)", k.Value))
}

func (k *Float) structPrint(sb *strings.Builder, indent int) {
	simplePrint(sb, indent, fmt.Sprintf("Float (%f)", k.Value))
}

func (k *String) structPrint(sb *strings.Builder, indent int) {
	simplePrint(sb, indent, fmt.Sprintf("String (%s)", k))
}

func (k *StringBuffer) structPrint(sb *strings.Builder, indent int) {
	simplePrint(sb, indent, fmt.Sprintf("StringBuffer (%s)", k.Value.String()))
}

func (k *Bytes) structPrint(sb *strings.Builder, indent int) {
	simplePrint(sb, indent, fmt.Sprintf("Bytes (%b)", k))
}

func (k *Bool) structPrint(sb *strings.Builder, indent int) {
	simplePrint(sb, indent, fmt.Sprintf("Bool (%t)", k.Value))
}

func (k *Bottom) structPrint(sb *strings.Builder, indent int) {
	simplePrint(sb, indent, "Bottom")
}

func (k KSequence) structPrint(sb *strings.Builder, indent int) {
	ks := k.ToSlice()
	addIndent(sb, indent)
	sb.WriteString("KSequence {")
	if len(ks) == 0 {
		sb.WriteString(" <empty> }")
	} else {
		for i, childk := range ks {
			sb.WriteString("\n")
			childk.structPrint(sb, indent+1)
			if i < len(ks)-1 {
				sb.WriteString(" ~>")
			}
		}
		sb.WriteRune('\n')
		addIndent(sb, indent)
		sb.WriteRune('}')
	}
}
