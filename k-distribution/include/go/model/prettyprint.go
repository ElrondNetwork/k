package %PACKAGE_MODEL%

import (
	"fmt"
	"strings"
)

func addIndent(sb *strings.Builder, indent int) {
	for i := 0; i < indent; i++ {
		sb.WriteString("    ")
	}
}

func printMapKey(sb *strings.Builder, kmk KMapKey) {
	sb.WriteString(kmk.String())
}

// PrettyPrint ... returns a representation of a K item that tries to be as readable as possible
func PrettyPrint(k K) string {
	return k.prettyPrint(0)
}

func (k *KApply) prettyPrint(indent int) string {
	var sb strings.Builder
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
		childStr := k.List[0].prettyPrint(0)
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
		for i, childk := range k.List {
			sb.WriteRune('\n')
			addIndent(&sb, indent+1)
			sb.WriteString(childk.prettyPrint(indent + 1))
			if !isKCell && i < len(k.List)-1 {
				sb.WriteString(",")
			}
		}
		sb.WriteRune('\n')
		addIndent(&sb, indent)
	}

	// end
	if isKCell {
		sb.WriteString("</")
		sb.WriteString(strings.TrimPrefix(lblName, "<"))
	} else {
		sb.WriteString(")")
	}

	return sb.String()
}

func (k *InjectedKLabel) prettyPrint(indent int) string {
	return fmt.Sprintf("InjectedKLabel(%s)", k.Label.Name())
}

func (k *KToken) prettyPrint(indent int) string {
	return fmt.Sprintf("%s: %s", k.Sort.Name(), k.Value)
}

func (k *KVariable) prettyPrint(indent int) string {
	return fmt.Sprintf("var %s", k.Name)
}

func (k *Map) prettyPrint(indent int) string {
	var sb strings.Builder
	sb.WriteString("Map Sort:")
	sb.WriteString(k.Sort.Name())
	sb.WriteString(", Label:")
	sb.WriteString(k.Label.Name())
	if len(k.Data) == 0 {
		sb.WriteString(" <empty>")
	} else {
		sb.WriteString(", Data: (")
		for k, v := range k.Data {
			sb.WriteString("\n")
			addIndent(&sb, indent+1)
			printMapKey(&sb, k)
			sb.WriteString(" => ")
			sb.WriteString(v.prettyPrint(indent + 1))
		}
		sb.WriteRune('\n')
		addIndent(&sb, indent)
		sb.WriteRune(')')
	}

	return sb.String()
}

func (k *List) prettyPrint(indent int) string {
	var sb strings.Builder
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
			addIndent(&sb, indent+1)
			sb.WriteString(item.prettyPrint(indent + 1))
		}
		sb.WriteRune('\n')
		addIndent(&sb, indent)
		sb.WriteRune(']')
	}

	return sb.String()
}

func (k *Set) prettyPrint(indent int) string {
	var sb strings.Builder
	sb.WriteString("Set Sort:")
	sb.WriteString(k.Sort.Name())
	sb.WriteString(", Label:")
	sb.WriteString(k.Label.Name())
	if len(k.Data) == 0 {
		sb.WriteString(" <empty>")
	} else {
		sb.WriteString(", Data: {")
		for k := range k.Data {
			sb.WriteString("\n")
			addIndent(&sb, indent+1)
			printMapKey(&sb, k)
		}
		sb.WriteRune('\n')
		addIndent(&sb, indent)
		sb.WriteRune('}')
	}

	return sb.String()
}

func (k *Array) prettyPrint(indent int) string {
	var sb strings.Builder
	sb.WriteString("Array Sort:")
	sb.WriteString(k.Sort.Name())
	slice := k.Data.ToSlice()
	if len(slice) == 0 {
		sb.WriteString(" <empty>")
	} else {
		sb.WriteString(", Data: [")
		for i, item := range slice {
			sb.WriteString("\n")
			addIndent(&sb, indent+1)
			sb.WriteString(fmt.Sprintf("[%d] => ", i))
			sb.WriteString(item.prettyPrint(indent + 1))
		}
		sb.WriteRune('\n')
		addIndent(&sb, indent)
		sb.WriteRune(']')
	}

	return sb.String()
}

func (k *Int) prettyPrint(indent int) string {
	return fmt.Sprintf("Int (%s)", k.Value.String())
}

func (k *MInt) prettyPrint(indent int) string {
	return fmt.Sprintf("MInt (%d)", k.Value)
}

func (k *Float) prettyPrint(indent int) string {
	return fmt.Sprintf("Float (%f)", k.Value)
}

func (k *String) prettyPrint(indent int) string {
	return fmt.Sprintf("String (%s)", k)
}

func (k *StringBuffer) prettyPrint(indent int) string {
	return "StringBuffer [not yet implemented]"
}

func (k *Bytes) prettyPrint(indent int) string {
	return fmt.Sprintf("Bytes (%b)", k)
}

func (k *Bool) prettyPrint(indent int) string {
	return fmt.Sprintf("Bool (%t)", k.Value)
}

func (k *Bottom) prettyPrint(indent int) string {
	return "Bottom"
}

func (k *KSequence) prettyPrint(indent int) string {
	var sb strings.Builder
	if len(k.Ks) == 0 {
		sb.WriteString(".Kseq")
	} else {
		for i, childk := range k.Ks {
			if i > 0 {
				addIndent(&sb, indent)
			}
			sb.WriteString(childk.prettyPrint(indent))
			if i < len(k.Ks)-1 {
				sb.WriteString(" ~>\n")
			}
		}
	}

	return sb.String()
}
