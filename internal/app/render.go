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
        ansiReverse = "\x1b[7m"
        ansiReset   = "\x1b[0m"
    )

    // 1. Setup Syntax Highlighting (Lexer & Style)
    // We try to match the filename, or fallback to plain text
    lexer := lexers.Match(m.buffer.Filename)
    if lexer == nil {
        lexer = lexers.Fallback
    }
    // You can swap "dracula" for "monokai", "nord", etc.
    style := styles.Get("nord")
    if style == nil {
        style = styles.Fallback
    }

    // Render buffer content
    for y := 0; y < screenH-1; y++ {
        lineIndex := y + m.buffer.ScrollY
        if lineIndex >= len(m.buffer.Lines) {
            b.WriteString("~\n")
            continue
        }

        rawLine := m.buffer.Lines[lineIndex]
        
        // 2. Tokenize the line
        // This breaks "func main()" into parts: [Keyword "func"], [Text " "], [Name "main"], ...
        iterator, err := lexer.Tokenise(nil, rawLine)
        if err != nil {
            // Fallback if tokenization fails
            iterator, _ = lexers.Fallback.Tokenise(nil, rawLine)
        }

        x := 0 // Screen column (visual x)
        rnIndex := 0 // Actual rune index in the line string
        var lineBuilder strings.Builder
        prevReversed := false

        // 3. Iterate Tokens instead of just Runes
        for _, token := range iterator.Tokens() {
            
            // Get the color for this token type (e.g., Keywords are Pink)
            entry := style.Get(token.Type)
            tokenColor := getAnsiColor(entry.Colour)
            
            // Apply the Syntax Color
            lineBuilder.WriteString(tokenColor)

            // Iterate the runes INSIDE the token
            for _, r := range token.Value {
                
                // --- Your Existing Cursor/Selection Logic ---
                inSel := m.mode.Current() == core.Visual && core.IsInSelection(m.visualStart, *m.buffer.Cursor, rnIndex, lineIndex)
                isCursorPos := rnIndex == m.buffer.Cursor.X && lineIndex == m.buffer.Cursor.Y
                shouldReverse := inSel || isCursorPos

                // Apply Cursor (Reverse Video) ON TOP of syntax color
                if shouldReverse && !prevReversed {
                    lineBuilder.WriteString(ansiReverse)
                }
                if !shouldReverse && prevReversed {
                    lineBuilder.WriteString(ansiReset)     // Reset everything
                    lineBuilder.WriteString(tokenColor)    // Re-apply syntax color!
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
                
                rnIndex++ // Increment the global rune counter
            }
        }

        // Handle the "Block Cursor" at the very end of the line (e.g. creating new text)
        // We check rnIndex here because that equals len(runes) after the loop
        if lineIndex == m.buffer.Cursor.Y && m.buffer.Cursor.X == rnIndex {
             if !prevReversed {
                lineBuilder.WriteString(ansiReverse)
            }
            lineBuilder.WriteRune(' ') 
            lineBuilder.WriteString(ansiReset)
            prevReversed = false
        }

        // Clean up any remaining styles
        if prevReversed {
            lineBuilder.WriteString(ansiReset)
        } else {
            // Ensure we reset syntax colors at the end of the line
            lineBuilder.WriteString(ansiReset)
        }

        b.WriteString(lineBuilder.String())
        b.WriteString("\n")
    }

    // --- Status Line (unchanged) ---
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
    // Safety check to prevent panic on very small screens
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

// --- Helper to convert Chroma Color to ANSI TrueColor string ---
func getAnsiColor(c chroma.Colour) string {
    r := int(c.Red())
    g := int(c.Green())
    b := int(c.Blue())
    // returns "\x1b[38;2;R;G;Bm"
    return fmt.Sprintf("\x1b[38;2;%d;%d;%dm", r, g, b)
}
