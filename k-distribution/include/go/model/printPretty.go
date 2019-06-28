%COMMENT%

package %PACKAGE_MODEL%

import (
	"fmt"
	"strings"
)

// PrettyPrint ... returns a representation of a K item that tries to be as readable as possible
// designed for debugging purposes only
func (ms *ModelState) PrettyPrint(ref KReference) string {
	var sb strings.Builder
	ms.prettyPrintToStringBuilder(&sb, ref, 0)
	return sb.String()
}

func (ms *ModelState) prettyPrintToStringBuilder(sb *strings.Builder, ref KReference, indent int) {
	switch ref.refType {
	case boolRef:
		sb.WriteString(fmt.Sprintf("Bool (%t)", IsTrue(ref)))
	case bottomRef:
		sb.WriteString("Bottom")
	case emptyKseqRef:
		sb.WriteString(" .K")
	case nonEmptyKseqRef:
		ks := ms.KSequenceToSlice(ref)
		if len(ks) == 0 {
			panic("K sequences of length 0 should have type emptyKseqRef, not nonEmptyKseqRef")
		} else if len(ks) == 1 {
			ms.prettyPrintToStringBuilder(sb, ks[0], indent)
		} else {
			for i, child := range ks {
				if i > 0 {
					addIndent(sb, indent)
				}
				ms.prettyPrintToStringBuilder(sb, child, indent)
				if i < len(ks)-1 {
					sb.WriteString(" ~>\n")
				} else {
					sb.WriteString(" ~> . ")
				}
			}
		}
	default:
		// object types
		obj := ms.getObject(ref)
		obj.prettyPrint(ms, sb, indent)
	}
}


func (k *KApply) prettyPrint(ms *ModelState, sb *strings.Builder, indent int) {
	lblName := k.Label.Name()
	isKCell := strings.HasPrefix(lblName, "<") && strings.HasSuffix(lblName, ">")

	// begin
	sb.WriteString(lblName)
	if !isKCell {
		sb.WriteString("(")
	}

	// contents
	done := false
	if len(k.List) == 0 {
		done = true
	}
	if !done && len(k.List) == 1 {
		var tempSb strings.Builder
		ms.prettyPrintToStringBuilder(&tempSb, k.List[0], 0)
		childStr := tempSb.String()
		if !strings.Contains(childStr, "\n") {
			// if only one child and its representation not too big, just put everything in one row
			if isKCell {
				sb.WriteString(" ")
			}
			sb.WriteString(childStr)
			if isKCell {
				sb.WriteString(" ")
			}
			done = true
		}
	}
	if !done {
		for i, child := range k.List {
			sb.WriteRune('\n')
			addIndent(sb, indent+1)
			ms.prettyPrintToStringBuilder(sb, child, indent+1)
			if !isKCell && i < len(k.List)-1 {
				sb.WriteString(",")
			}
		}
		sb.WriteRune('\n')
		addIndent(sb, indent)
	}

	// end
	if isKCell {
		sb.WriteString("</")
		sb.WriteString(strings.TrimPrefix(lblName, "<"))
	} else {
		sb.WriteString(")")
	}
}

func (k *InjectedKLabel) prettyPrint(ms *ModelState, sb *strings.Builder, indent int) {
	sb.WriteString(fmt.Sprintf("InjectedKLabel(%s)", k.Label.Name()))
}

func (k *KToken) prettyPrint(ms *ModelState, sb *strings.Builder, indent int) {
	sb.WriteString(fmt.Sprintf("%s: %s", k.Sort.Name(), k.Value))
}

func (k *KVariable) prettyPrint(ms *ModelState, sb *strings.Builder, indent int) {
	sb.WriteString(fmt.Sprintf("var %s", k.Name))
}

func (k *Map) prettyPrint(ms *ModelState, sb *strings.Builder, indent int) {
	sb.WriteString("Map Sort:")
	sb.WriteString(k.Sort.Name())
	sb.WriteString(", Label:")
	sb.WriteString(k.Label.Name())
	if len(k.Data) == 0 {
		sb.WriteString(" <empty>")
	} else {
		sb.WriteString(", Data: (")
		orderedKVPairs := ms.MapOrderedKeyValuePairs(k)
		for _, pair := range orderedKVPairs {
			sb.WriteString("\n")
			addIndent(sb, indent+1)
			ms.prettyPrintToStringBuilder(sb, pair.Key, indent+1)
			sb.WriteString(" => ")
			ms.prettyPrintToStringBuilder(sb, pair.Value, indent+1)
		}
		sb.WriteRune('\n')
		addIndent(sb, indent)
		sb.WriteRune(')')
	}
}

func (k *List) prettyPrint(ms *ModelState, sb *strings.Builder, indent int) {
	sb.WriteString("List Sort:")
	sb.WriteString(k.Sort.Name())
	sb.WriteString(", Label:")
	sb.WriteString(k.Label.Name())
	if len(k.Data) == 0 {
		sb.WriteString(" <empty>")
	} else {
		sb.WriteString(", Data: [")
		for _, item := range k.Data {
			sb.WriteString("\n")
			addIndent(sb, indent+1)
			if item == NullReference {
				sb.WriteString("nil")
			} else {
				ms.prettyPrintToStringBuilder(sb, item, indent+1)
			}
		}
		sb.WriteRune('\n')
		addIndent(sb, indent)
		sb.WriteRune(']')
	}
}

func (k *Set) prettyPrint(ms *ModelState, sb *strings.Builder, indent int) {
	sb.WriteString("Set Sort:")
	sb.WriteString(k.Sort.Name())
	sb.WriteString(", Label:")
	sb.WriteString(k.Label.Name())
	if len(k.Data) == 0 {
		sb.WriteString(" <empty>")
	} else {
		sb.WriteString(", Data: {")
		orderedElems := ms.SetOrderedElements(k)
		for _, elem := range orderedElems {
			sb.WriteString("\n")
			addIndent(sb, indent+1)
			ms.prettyPrintToStringBuilder(sb, elem, indent+1)
		}

		sb.WriteRune('\n')
		addIndent(sb, indent)
		sb.WriteRune('}')
	}
}

func (k *Array) prettyPrint(ms *ModelState, sb *strings.Builder, indent int) {
	sb.WriteString("Array Sort:")
	sb.WriteString(k.Sort.Name())
	slice := k.Data.ToSlice()
	if len(slice) == 0 {
		sb.WriteString(" <empty>")
	} else {
		sb.WriteString(", Data: [")
		for i, item := range slice {
			sb.WriteString("\n")
			addIndent(sb, indent+1)
			sb.WriteString(fmt.Sprintf("[%d] => ", i))
			ms.prettyPrintToStringBuilder(sb, item, indent+1)
		}
		sb.WriteRune('\n')
		addIndent(sb, indent)
		sb.WriteRune(']')
	}
}

func (k *BigInt) prettyPrint(ms *ModelState, sb *strings.Builder, indent int) {
	sb.WriteString(fmt.Sprintf("Int (0x%s | %s)", k.Value.Text(16), k.Value.String()))
}

func (k *MInt) prettyPrint(ms *ModelState, sb *strings.Builder, indent int) {
	sb.WriteString(fmt.Sprintf("MInt (%d)", k.Value))
}

func (k *Float) prettyPrint(ms *ModelState, sb *strings.Builder, indent int) {
	sb.WriteString(fmt.Sprintf("Float (%f)", k.Value))
}

func (k *String) prettyPrint(ms *ModelState, sb *strings.Builder, indent int) {
	sb.WriteString("String(\"")
	writeEscapedString(sb, k.Value)
	sb.WriteString("\")")
}

func (k *StringBuffer) prettyPrint(ms *ModelState, sb *strings.Builder, indent int) {
	sb.WriteString("StringBuffer(\"")
	writeEscapedString(sb, k.Value.String())
	sb.WriteString("\")")
}

func (k *Bytes) prettyPrint(ms *ModelState, sb *strings.Builder, indent int) {
	sb.WriteString("Bytes(")
	if len(k.Value) == 0 {
		sb.WriteString("empty")
	} else {
		for i, b := range k.Value {
			sb.WriteString(fmt.Sprintf("%02x", b))
			if i < len(k.Value)-1 {
				sb.WriteByte(' ')
			}
		}
	}
	sb.WriteString(")")
}
