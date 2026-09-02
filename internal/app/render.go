package app

import (
    "fmt"
    "strings"

    "github.com/impossibleclone/vimingo/internal/core"
    "github.com/alecthomas/chroma/v2"
    "github.com/alecthomas/chroma/v2/lexers"
    "github.com/alecthomas/chroma/v2/styles"
)

func (m *Model) View() string {
    var b strings.Builder

    screenW, screenH := m.width, m.height
    if screenW == 0 || screenH == 0 {
        return "Initializing..."
    }

    const (
        ansiReverse   = "\x1b[7m"
        ansiReverseOff = "\x1b[27m" // Only turns off reverse, keeps color!
        ansiReset     = "\x1b[0m"
    )

    // Setup Syntax Highlighting
    lexer := lexers.Match(m.buffer.Filename)
    if lexer == nil {
        lexer = lexers.Fallback
    }
    style := styles.Get("nord")
    if style == nil {
        style = styles.Fallback
    }

    for y := 0; y < screenH-1; y++ {
        lineIndex := y + m.buffer.ScrollY
        if lineIndex >= len(m.buffer.Lines) {
            b.WriteString("~\n")
            continue
        }

        rawLine := m.buffer.Lines[lineIndex]
        iterator, err := lexer.Tokenise(nil, rawLine)
        if err != nil {
            iterator, _ = lexers.Fallback.Tokenise(nil, rawLine)
        }

        x := 0
        rnIndex := 0
        var lineBuilder strings.Builder
        prevReversed := false

        for _, token := range iterator.Tokens() {
            entry := style.Get(token.Type)
            tokenColor := getAnsiColor(entry.Colour)
            
            lineBuilder.WriteString(tokenColor)

            for _, r := range token.Value {
                inSel := m.mode.Current() == core.Visual && core.IsInSelection(m.visualStart, *m.buffer.Cursor, rnIndex, lineIndex)
                isCursorPos := rnIndex == m.buffer.Cursor.X && lineIndex == m.buffer.Cursor.Y
                shouldReverse := inSel || isCursorPos

                // Turn Reverse ON
                if shouldReverse && !prevReversed {
                    lineBuilder.WriteString(ansiReverse)
                }
                // Turn Reverse OFF (Only Reverse!)
                if !shouldReverse && prevReversed {
                    lineBuilder.WriteString(ansiReverseOff)
                }
                prevReversed = shouldReverse

                // Handle Tabs
                if r == '\t' {
                    spaces := 4 - (x % 4)
                    for t := 0; t < spaces; t++ {
                        lineBuilder.WriteRune(' ')
                        x++
                    }
                } else {
                    lineBuilder.WriteRune(r)
                    x++
                }
                rnIndex++
            }
        }

        if lineIndex == m.buffer.Cursor.Y && m.buffer.Cursor.X == rnIndex {
            if !prevReversed {
                lineBuilder.WriteString(ansiReverse)
            }
            lineBuilder.WriteRune(' ')
            lineBuilder.WriteString(ansiReverseOff)
            prevReversed = false
        }

        if prevReversed {
            lineBuilder.WriteString(ansiReverseOff)
        }
        lineBuilder.WriteString(ansiReset)

        b.WriteString(lineBuilder.String())
        b.WriteString("\n")
    }

    status := ""
    if m.mode.Current() == core.Normal {
        status = " -- NORMAL -- " + m.buffer.Filename
    } else if m.mode.Current() == core.Insert {
        status = " -- INSERT -- " + m.buffer.Filename
    } else if m.mode.Current() == core.Visual {
        status = " -- VISUAL -- " + m.buffer.Filename
    } else if m.mode.Current() == core.Command {
        status = ":" + string(m.buffer.Command)
    }
    if m.buffer.StatusMsg != "" && m.mode.Current() != core.Command {
        status += " | " + m.buffer.StatusMsg
    }

    coords := fmt.Sprintf("%d:%d ", m.buffer.Cursor.Y+1, m.buffer.Cursor.X+1)
    operators := string(m.buffer.KeyReg)

    statusLine := make([]rune, screenW)
    for i := range statusLine {
        statusLine[i] = ' '
    }

    copy(statusLine, []rune(status))

    coordStart := screenW - len(coords)
    if coordStart > 0 && coordStart < len(statusLine) {
        copy(statusLine[coordStart:], []rune(coords))
    }

    opStart := coordStart - len(operators)
    if opStart > 0 && opStart < len(statusLine) {
         copy(statusLine[opStart:], []rune(operators))
    }

    b.WriteString(string(statusLine))

    return b.String()
}

func getAnsiColor(c chroma.Colour) string {
    if int(c) == 0 {
        return "\x1b[39m" // 39m = Default Foreground Color
    }
    r := int(c.Red())
    g := int(c.Green())
    b := int(c.Blue())
    return fmt.Sprintf("\x1b[38;2;%d;%d;%dm", r, g, b)
}
