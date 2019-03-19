package %PACKAGE_MODEL%

import (
	"fmt"
	"strings"
)

// StructPrint ... returns a representation of a K item that somewhat resembles a Go declaration
func StructPrint(k K) string {
	return k.structPrint(0)
}

func simplePrint(indent int, str string) string {
	var sb strings.Builder
	addIndent(&sb, indent)
	sb.WriteString(str)
	return sb.String()
}

func (k *KApply) structPrint(indent int) string {
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
			sb.WriteString(childk.structPrint(indent + 1))
		}
		sb.WriteRune('\n')
		addIndent(&sb, indent)
		sb.WriteRune('}')
	}

	return sb.String()
}

func (k *InjectedKLabel) structPrint(indent int) string {
	return simplePrint(indent, fmt.Sprintf("InjectedKLabel {Label:%s}", k.Label.Name()))
}

func (k *KToken) structPrint(indent int) string {
	return simplePrint(indent, fmt.Sprintf("KToken {Sort:%s, Value:%s}", k.Sort.Name(), k.Value))
}

func (k *KVariable) structPrint(indent int) string {
	return simplePrint(indent, fmt.Sprintf("KVariable {Name:%s}", k.Name))
}

func (k *Map) structPrint(indent int) string {
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
			sb.WriteString(v.structPrint(0))
		}
		sb.WriteRune('\n')
		addIndent(&sb, indent)
		sb.WriteRune('}')
	}

	return sb.String()
}

func (k *List) structPrint(indent int) string {
	// TODO: print data
	return simplePrint(indent, fmt.Sprintf("List {Sort:%s, Label:%s}", k.Sort.Name(), k.Label.Name()))
}

func (k *Set) structPrint(indent int) string {
	// TODO: print data
	return simplePrint(indent, fmt.Sprintf("Set {Sort:%s, Label:%s}", k.Sort.Name(), k.Label.Name()))
}

func (k *Array) structPrint(indent int) string {
	// TODO: print data
	return simplePrint(indent, fmt.Sprintf("Array {Sort:%s}", k.Sort.Name()))
}

func (k *Int) structPrint(indent int) string {
	return simplePrint(indent, fmt.Sprintf("Int (%s)", k.Value.String()))
}

func (k *MInt) structPrint(indent int) string {
	return simplePrint(indent, fmt.Sprintf("MInt (%d)", k.Value))
}

func (k *Float) structPrint(indent int) string {
	return simplePrint(indent, fmt.Sprintf("Float (%f)", k.Value))
}

func (k *String) structPrint(indent int) string {
	return simplePrint(indent, fmt.Sprintf("String (%s)", k))
}

func (k *StringBuffer) structPrint(indent int) string {
	return fmt.Sprintf("StringBuffer (%s)", k.Value.String())
}

func (k *Bytes) structPrint(indent int) string {
	return simplePrint(indent, fmt.Sprintf("Bytes (%b)", k))
}

func (k *Bool) structPrint(indent int) string {
	return simplePrint(indent, fmt.Sprintf("Bool (%t)", k.Value))
}

func (k *Bottom) structPrint(indent int) string {
	return simplePrint(indent, "Bottom")
}

func (k *KSequence) structPrint(indent int) string {
	var sb strings.Builder
	addIndent(&sb, indent)
	sb.WriteString("KSequence {")
	if len(k.Ks) == 0 {
		sb.WriteString(" <empty> }")
	} else {
		for i, childk := range k.Ks {
			sb.WriteString("\n")
			sb.WriteString(childk.structPrint(indent + 1))
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
