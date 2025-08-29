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
