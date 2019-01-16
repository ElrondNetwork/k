package org.kframework.backend.go;

/**
 * Wrapper around a StringBuilder that accumulates Go code. Also handles indent.
 */
public class GoStringBuilder {

    public static final int FUNCTION_BODY_INDENT = 1;

    private StringBuilder sb;
    private int indent = 0;

    public GoStringBuilder() {
        sb = new StringBuilder();
    }

    public GoStringBuilder(StringBuilder sb, int indent) {
        this.sb = sb;
        this.indent = indent;
    }

    public GoStringBuilder append(String s) {
        sb.append(s);
        return this;
    }

    public GoStringBuilder append(int s) {
        sb.append(s);
        return this;
    }

    public GoStringBuilder newLine() {
        sb.append('\n');
        return this;
    }

    public GoStringBuilder writeIndent() {
        for (int i = 0; i < indent; i++) {
            sb.append('\t');
        }
        return this;
    }

    public GoStringBuilder beginBlock() {
        sb.append(" {\n");
        indent++;
        return this;
    }

    public GoStringBuilder beginBlock(String comment) {
        sb.append(" { // ");
        sb.append(comment);
        sb.append('\n');
        indent++;
        return this;
    }

    public GoStringBuilder endAllBlocks(int finalIndent) {
        while (indent > finalIndent) {
            indent--;
            writeIndent();
            sb.append("}\n");
        }
        return this;
    }

    public void increaseIndent() {
        indent++;
    }

    public int currentIndent() {
        return indent;
    }

    public void decreaseIndent() {
        indent--;
    }

    public String toString() {
        return sb.toString();
    }

    public StringBuilder sb() {
        return sb;
    }

}
