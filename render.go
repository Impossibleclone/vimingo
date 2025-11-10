package main

import (
	"fmt"
	"strings"
)

func (m *model) View() string {
	var b strings.Builder

	screenW, screenH := m.width, m.height
	if screenW == 0 || screenH == 0 {
		return "Initializing..."
	}

	const (
		ansiReverse = "\x1b[7m"
		ansiReset   = "\x1b[0m"
	)

	// Render buffer content
	for y := 0; y < screenH-1; y++ {
		lineIndex := y + m.buffer.ScrollY
		if lineIndex >= len(m.buffer.Lines) {
			b.WriteString("~\n")
			continue
		}

		runes := []rune(m.buffer.Lines[lineIndex])
		x := 0 // screen column
		var lineBuilder strings.Builder
		
		prevReversed := false 

		for rn := 0; rn < len(runes); rn++ {
			r := runes[rn]
			
			inSel := m.mode.Current() == Visual && isInSelection(m.visualStart, *m.buffer.Cursor, rn, lineIndex)
			
			isCursorPos := rn == m.buffer.Cursor.X && lineIndex == m.buffer.Cursor.Y

			shouldReverse := inSel || isCursorPos

			if shouldReverse && !prevReversed {
				lineBuilder.WriteString(ansiReverse)
			}
			if !shouldReverse && prevReversed {
				lineBuilder.WriteString(ansiReset)
			}
			prevReversed = shouldReverse // Update state for next iteration

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
		}
		
		if lineIndex == m.buffer.Cursor.Y && m.buffer.Cursor.X == len(runes) {
			if !prevReversed {
				lineBuilder.WriteString(ansiReverse) // Turn on reverse if not already on
			}
			lineBuilder.WriteRune(' ') // Draw the "block" cursor
			lineBuilder.WriteString(ansiReset)   // Immediately reset
			prevReversed = false // We are now reset
		}

		if prevReversed {
			lineBuilder.WriteString(ansiReset)
		}

		b.WriteString(lineBuilder.String())
		b.WriteString("\n")
	}

	status := ""
	if m.mode.Current() == Normal {
		status = "-- NORMAL --" + " \\ " + m.buffer.Filename
	} else if m.mode.Current() == Insert {
		status = "-- INSERT --" + " \\ " + m.buffer.Filename
	} else if m.mode.Current() == Command {
		status = ":" + string(m.buffer.Command)
	}
	if m.buffer.StatusMsg != "" && m.mode.Current() != Command {
		status += " | " + m.buffer.StatusMsg
	}

	coords := fmt.Sprintf("/ %d:%d ", m.buffer.Cursor.Y+1, m.buffer.Cursor.X+1)
	operators := string(m.buffer.KeyReg)

	statusLine := make([]rune, screenW)
	for i := range statusLine {
		statusLine[i] = ' '
	}

	copy(statusLine, []rune(status))

	coordStart := screenW - len(coords)
	copy(statusLine[coordStart:], []rune(coords))

	opStart := coordStart - len(operators)
	if opStart < 0 {
		opStart = 0
	}
	copy(statusLine[opStart:], []rune(operators))

	b.WriteString(string(statusLine))

	return b.String()
}
