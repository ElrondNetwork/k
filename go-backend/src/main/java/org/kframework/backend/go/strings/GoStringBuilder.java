// Copyright (c) 2015-2019 K Team. All Rights Reserved.
package org.kframework.backend.go.strings;

import java.util.Stack;
import java.util.function.Consumer;

/**
 * Wrapper around a StringBuilder that accumulates Go code. Also handles tabsIndent.
 */
public class GoStringBuilder {

    public static final int FUNCTION_BODY_INDENT = 1;
    public static final int FUNCTION_MAIN_LOOP_INDENT = 2;

    private StringBuilder sb;
    private int tabsIndent = 0; // main indent
    private int spacesIndent = 0; // only for special cases

    private class BlockEndCallback {
        final int indent;
        final Consumer<GoStringBuilder> callback;
        public BlockEndCallback(int indent, Consumer<GoStringBuilder> callback) {
            this.indent = indent;
            this.callback = callback;
        }
    }

    private final Stack<BlockEndCallback> callbacksBeforeBlockEnd = new Stack<>();
    private final Stack<BlockEndCallback> callbacksAfterBlockEnd = new Stack<>();

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

    /**
     * Opens a new block after an instruction. Increases indent and inserts a newline.
     * @param comments optional text to print as comment after '{'
     * @return this, for call chaining
     */
    public GoStringBuilder beginBlock(String... comments) {
        sb.append(" {");
        if (comments.length > 0) {
            sb.append(" // ");
            for (String comment : comments) {
                sb.append(comment);
            }
        }
        sb.append('\n');
        tabsIndent++;
        return this;
    }

    /**
     * Opens a new block after an instruction. Increases indent and inserts a newline.
     * @param comments optional text to print as comment after '{'
     * @return this, for call chaining
     */
    public GoStringBuilder scopingBlock(String... comments) {
        writeIndent();
        sb.append("{");
        if (comments.length > 0) {
            sb.append(" // ");
            for (String comment : comments) {
                sb.append(comment);
            }
        }
        sb.append('\n');
        tabsIndent++;
        return this;
    }

    public GoStringBuilder addCallbackBeforeReturningFromBlock(int blockIndent, Consumer<GoStringBuilder> callback) {
        callbacksBeforeBlockEnd.push(new BlockEndCallback(blockIndent, callback));
        return this;
    }

    public GoStringBuilder addCallbackAfterReturningFromBlock(int blockIndent, Consumer<GoStringBuilder> callback) {
        callbacksAfterBlockEnd.push(new BlockEndCallback(blockIndent, callback));
        return this;
    }

    public GoStringBuilder endOneBlockNoNewline() {
        if (!callbacksBeforeBlockEnd.empty() && callbacksBeforeBlockEnd.peek().indent == tabsIndent - 1) {
            callbacksBeforeBlockEnd.pop().callback.accept(this);
        }
        tabsIndent--;
        writeIndent();
        sb.append("}");
        if (!callbacksAfterBlockEnd.empty() && callbacksAfterBlockEnd.peek().indent == tabsIndent) {
            callbacksAfterBlockEnd.pop().callback.accept(this);
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
