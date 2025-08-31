package main

func isInSelection(start, end Cursor, x, y int) bool {
	if start.Y > end.Y || (start.Y == end.Y && start.X > end.X) {
		start, end = end, start
	}

	if y < start.Y || y > end.Y {
		return false
	}

	if y == start.Y && y == end.Y {
		return x >= start.X && x <= end.X
	} else if y == start.Y {
		return x >= start.X
	} else if y == end.Y {
		return x <= end.X
	}

	return true
}
func YankSelection(buffer *Buffer, cursor *Cursor, visualStart *Cursor) {
	start := min(visualStart.X, cursor.X)
	end := max(visualStart.X, cursor.X)
	startline := min(visualStart.Y, cursor.Y)
	endline := max(visualStart.Y, cursor.Y)

	buffer.Register = ""

	if startline == endline {
		toYankFromLine := buffer.Lines[startline]
		toYankTheCharacters := []rune(toYankFromLine[start : end+1])
		buffer.Register = string(toYankTheCharacters)
	} else {
		for y := startline; y <= endline; y++ {
			if y == startline {
				toYankFromLine := buffer.Lines[y]
				toYankTheCharacters := []rune(toYankFromLine[start:])
				buffer.Register += string(toYankTheCharacters)
				buffer.Register += "\n"

			} else if y == endline {
				toYankFromLine := buffer.Lines[y]
				toYankTheCharacters := []rune(toYankFromLine[:end])
				buffer.Register += string(toYankTheCharacters)

			} else {
				toYankFromLine := buffer.Lines[y]
				toYankTheCharacters := []rune(toYankFromLine[:])
				buffer.Register += string(toYankTheCharacters)
				buffer.Register += "\n"
			}
		}
	}

	cursor.Y = startline
	cursor.X = start
}

func CutSelection(buffer *Buffer, cursor *Cursor, visualStart *Cursor) {
	start := min(visualStart.X, cursor.X)
	end := max(visualStart.X, cursor.X)
	startline := min(visualStart.Y, cursor.Y)
	endline := max(visualStart.Y, cursor.Y)

	buffer.Register = ""

	if startline == endline {
		// same line
		toCutFromLine := buffer.Lines[startline]
		toCutTheCharacters := []rune(toCutFromLine[start : end+1])
		buffer.Register = string(toCutTheCharacters)

		// remove cut characters from line
		newLine := []rune(toCutFromLine[:start])
		newLine = append(newLine, []rune(toCutFromLine[end+1:])...)
		buffer.Lines[startline] = string(newLine)

	} else {
		// multiple lines
		for y := startline; y <= endline; y++ {
			if y == startline {
				toCutFromLine := buffer.Lines[y]
				toCutTheCharacters := []rune(toCutFromLine[start:])
				buffer.Register += string(toCutTheCharacters)
				buffer.Register += "\n"
				buffer.Lines[y] = string([]rune(toCutFromLine[:start]))

			} else if y == endline {
				toCutFromLine := buffer.Lines[y]
				toCutTheCharacters := []rune(toCutFromLine[:end+1])
				buffer.Register += string(toCutTheCharacters)
				// keep only part after end
				buffer.Lines[y] = string([]rune(toCutFromLine[end+1:]))

			} else {
				toCutFromLine := buffer.Lines[y]
				buffer.Register += toCutFromLine + "\n"
				// delete whole line
				buffer.Lines[y] = ""
			}
		}
	}

	DeleteSelection(buffer, cursor, visualStart)
}

func DeleteSelection(buffer *Buffer, cursor *Cursor, visualStart *Cursor) {
	start := min(visualStart.X, cursor.X)
	end := max(visualStart.X, cursor.X)
	startline := min(visualStart.Y, cursor.Y)
	endline := max(visualStart.Y, cursor.Y)

	buffer.Register = ""

	if startline == endline {
		// single line delete
		line := []rune(buffer.Lines[startline])
		toDelete := line[start : end+1]
		buffer.Register = string(toDelete)

		newLine := append(line[:start], line[end+1:]...)
		buffer.Lines[startline] = string(newLine)

	} else {
		// multi-line delete
		startLineRunes := []rune(buffer.Lines[startline])
		endLineRunes := []rune(buffer.Lines[endline])

		// yank text into register
		buffer.Register += string(startLineRunes[start:]) + "\n"
		for y := startline + 1; y < endline; y++ {
			buffer.Register += buffer.Lines[y] + "\n"
		}
		buffer.Register += string(endLineRunes[:end+1])

		// merge first + last line pieces
		newStart := string(startLineRunes[:start]) + string(endLineRunes[end+1:])
		buffer.Lines[startline] = newStart

		// delete the in-between lines
		buffer.Lines = append(buffer.Lines[:startline+1], buffer.Lines[endline+1:]...)
	}

	// place cursor at start of deleted region
	cursor.Y = startline
	cursor.X = start
}
