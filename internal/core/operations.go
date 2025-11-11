package core

import "strings"

func TypeCh(line string, pos int, ch rune) string {
	return line[:pos] + string(ch) + line[pos:]
}

func insertText(line string, pos int, text string) string {
	return line[:pos] + text + line[pos:]
}

func YankRange(buffer *Buffer, cursor *Cursor, end int) {
	start := cursor.X
	buffer.Register = ""
	toYankFromLine := buffer.Lines[cursor.Y]
	if end >= len(toYankFromLine) {
		end = len(toYankFromLine) - 1
	}
	toYankTheCharacters := []rune(toYankFromLine[start:end])
	buffer.Register = string(toYankTheCharacters)
}

func Paste(buffer *Buffer, yankedText string) {
	lines := strings.Split(yankedText, "\n")
	if len(lines) == 1 {
		buffer.Lines[buffer.Cursor.Y] = insertText(buffer.Lines[buffer.Cursor.Y], buffer.Cursor.X, lines[0])
	} else if len(lines) > 1 {
		for i := 0; i < len(lines); i++ {
			if i == 0 { //for the first line yanked
				buffer.Lines[buffer.Cursor.Y] = insertText(buffer.Lines[buffer.Cursor.Y], buffer.Cursor.X, lines[0])
				buffer.Cursor.X += len(lines[i])
				NewLine(buffer)
			} else if i > 0 && i < len(lines)-1 { //for the middle lines
				buffer.Lines[buffer.Cursor.Y] = insertText(buffer.Lines[buffer.Cursor.Y], buffer.Cursor.X, lines[i])
				buffer.Cursor.X = len(lines[i])
				NewLine(buffer)
			} else if i > 0 && i < len(lines) { //for the last line
				buffer.Lines[buffer.Cursor.Y] = insertText(buffer.Lines[buffer.Cursor.Y], buffer.Cursor.X, lines[i])
				buffer.Cursor.X = len(lines[i])
			}
		}

	}
	// for i := 1; i < len(lines); i++ {
	// 	NewLine(buffer)
	// 	buffer.Lines[buffer.Cursor.Y] = lines[i]
	// }
}

func RemoveCh(line string, pos int) string {
	return line[:pos] + line[pos+1:]
}

func RemoveChs(line string, startpos int, endpos int) string {
	return line[:startpos] + line[endpos:]
}

// func yank(line string, start *Cursor) string {
//
// }

func SplitLine(line string, pos int) (before string, after string) {
	return line[:pos], line[pos:]
}

func NewLine(buffer *Buffer) {
	before, after := SplitLine(buffer.Lines[buffer.Cursor.Y], buffer.Cursor.X)
	buffer.Lines[buffer.Cursor.Y] = before
	buffer.Lines = append(
		buffer.Lines[:buffer.Cursor.Y+1],
		append([]string{after}, buffer.Lines[buffer.Cursor.Y+1:]...)...,
	)
	buffer.Cursor.Y++
	buffer.Cursor.X = 0
}

func RemoveLine(buffer *Buffer) {
	if buffer.Cursor.X == 0 && buffer.Cursor.Y > 0 {
		prevLine := buffer.Lines[buffer.Cursor.Y-1]
		buffer.Lines[buffer.Cursor.Y-1] = buffer.Lines[buffer.Cursor.Y-1] + buffer.Lines[buffer.Cursor.Y]
		buffer.Lines = append(buffer.Lines[:buffer.Cursor.Y], buffer.Lines[buffer.Cursor.Y+1:]...)
		buffer.Cursor.Y--
		buffer.Cursor.X = len(prevLine)
	}
}

func AdjustScroll(buffer *Buffer, screenH int) {
	textHeight := screenH - 1 // reserve bottom line

	if buffer.Cursor.Y < buffer.ScrollY {
		buffer.ScrollY = buffer.Cursor.Y
	} else if buffer.Cursor.Y >= buffer.ScrollY+textHeight {
		buffer.ScrollY = buffer.Cursor.Y - textHeight + 1
	}
}
