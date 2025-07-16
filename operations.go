package main

func TypeCh(line string, pos int, ch rune) string {
	return line[:pos] + string(ch) + line[pos:]
}

func RemoveCh(line string, pos int) string {
	return line[:pos] + line[pos+1:]
}

// func yank(line string, pos int) string {
// 	return line[:pos] 
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

func adjustScroll(buffer *Buffer, screenH int) {
    if buffer.Cursor.Y < buffer.ScrollY {
        buffer.ScrollY = buffer.Cursor.Y
    } else if buffer.Cursor.Y >= buffer.ScrollY + screenH {
        buffer.ScrollY = buffer.Cursor.Y - screenH + 1
    }
}
