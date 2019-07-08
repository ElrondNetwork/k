%COMMENT%

package %PACKAGE%

import (
	"fmt"
	"strings"
)

func addIndent(sb *strings.Builder, indent int) {
	for i := 0; i < indent; i++ {
		sb.WriteString("    ")
	}
}

func writeEscapedChar(sb *strings.Builder, c byte) {
	if c == '\n' {
		sb.WriteString("\\n")
	} else if c == '\t' {
		sb.WriteString("\\t")
	} else if c == '"' {
		sb.WriteString("\\\"")
	} else if c < ' ' || c >= 0x7f {
		sb.WriteString(fmt.Sprintf("\\x%02x", c))
	} else {
		sb.WriteByte(c)
	}
}

func writeEscapedString(sb *strings.Builder, s string) {
	for _, c := range []byte(s) {
		writeEscapedChar(sb, c)
	}
}
