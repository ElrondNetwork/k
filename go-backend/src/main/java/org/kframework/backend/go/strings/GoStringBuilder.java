package org.kframework.backend.go.strings;

import java.util.Stack;
import java.util.function.Consumer;

/**
 * Wrapper around a StringBuilder that accumulates Go code. Also handles tabsIndent.
 */
public class GoStringBuilder {

    public static final int FUNCTION_BODY_INDENT = 1;

    private StringBuilder sb;
    private int tabsIndent = 0; // main indent
    private int spacesIndent = 0; // only for special cases

    private class AfterBlockEndCallback {
        final int indent;
        final Consumer<GoStringBuilder> callback;
        public AfterBlockEndCallback(int indent, Consumer<GoStringBuilder> callback) {
            this.indent = indent;
            this.callback = callback;
        }
    }

    private final Stack<AfterBlockEndCallback> blockEndCallbackStack = new Stack<>();

    public GoStringBuilder() {
        sb = new StringBuilder();
    }

    public GoStringBuilder(int tabsIndent, int spacesIndent) {
        sb = new StringBuilder();
        this.tabsIndent = tabsIndent;
        this.spacesIndent = spacesIndent;
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

    public GoStringBuilder appendIndentedLine(String... strings) {
        writeIndent();
        for (String s : strings) {
            sb.append(s);
        }
        newLine();
        return this;
    }

    public GoStringBuilder writeIndent() {
        for (int i = 0; i < tabsIndent; i++) {
            sb.append('\t');
        }
        for (int i = 0; i < spacesIndent; i++) {
            sb.append(' ');
        }
        return this;
    }

    public GoStringBuilder beginBlock() {
        sb.append(" {\n");
        tabsIndent++;
        return this;
    }

    public GoStringBuilder beginBlock(String comment) {
        sb.append(" { // ");
        sb.append(comment);
        sb.append('\n');
        tabsIndent++;
        return this;
    }

    public GoStringBuilder addCallbackWhenReturningFromBlock(int blockIndent, Consumer<GoStringBuilder> callback) {
        blockEndCallbackStack.push(new AfterBlockEndCallback(blockIndent, callback));
        return this;
    }

    public GoStringBuilder endOneBlockNoNewline() {
        tabsIndent--;
        writeIndent();
        sb.append("}");
        if (!blockEndCallbackStack.empty() && blockEndCallbackStack.peek().indent == tabsIndent) {
            blockEndCallbackStack.pop().callback.accept(this);
        }
        return this;
    }

    public GoStringBuilder endOneBlock() {
        endOneBlockNoNewline();
        newLine();
        return this;
    }

    public GoStringBuilder endAllBlocks(int finalIndent) {
        while (tabsIndent > finalIndent) {
            endOneBlock();
        }
        return this;
    }

    public void increaseIndent() {
        increaseIndent(1);
    }

    public void increaseIndent(int amount) {
        tabsIndent += amount;
    }

    public void decreaseIndent() {
        decreaseIndent(1);
    }

    public void decreaseIndent(int amount) {
        tabsIndent -= amount;
    }

    public void forceIndent(int newIndent) {
        this.tabsIndent = newIndent;
    }

    public int getCurrentIndent() {
        return tabsIndent;
    }

    /**
     * Add some spaces to each indent <br/>
     * e.g. in the if-condition, it's nice to align stuff like newlines after '&&' to the 'if'
     *
     * @param reference the spaces indent should have the width of this string, e.g. "if "
     */
    public void enableMiniIndent(String reference) {
        spacesIndent = reference.length();
    }

    public void disableMiniIndent() {
        spacesIndent = 0;
    }

    public String toString() {
        return sb.toString();
    }

    public StringBuilder sb() {
        return sb;
    }

}
