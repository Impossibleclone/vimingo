package main

func wMotion(buffer *Buffer, cursor *Cursor) int {
	line := buffer.Lines[cursor.Y]
	pos := cursor.X + 1
	// Skip current word characters
	for pos < len(line) && line[pos] != ' ' {
		pos++
	}
	// Skip spaces to reach start of next word
	for pos < len(line) && line[pos] == ' ' {
		pos++
	}
	return pos
}

const TabStop = 4 // Define tab width for cursor movement consistency

func eMotion(buffer *Buffer, cursor *Cursor) int {
	line := buffer.Lines[cursor.Y]
	pos := cursor.X
	if pos < len(line)-1 && line[pos+1] == ' ' {
		pos++
	}
	if pos < len(line)-1 && line[pos+1] == '\t'{
		pos+= TabStop
	}
	for x := pos; x < len(line)-1; x++ {
		if x < len(line) && (line[x] == ' ' || line[x] == '\t'){
			continue
		}
		if x < len(line) && (line[x+1] == ' ' || line[x+1] == '\t'){
			return x
		}
	}
	return len(line)-1
}

// func eMotion(buffer *Buffer, cursor *Cursor) int {
//     line := buffer.Lines[cursor.Y]
//     pos := cursor.X
//
//     // Skip whitespace (spaces and tabs) to find start of next word
//     for pos < len(line) && (line[pos] == ' ' || line[pos] == '\t') {
//         if line[pos] == '\t' {
//             pos += TabStop // Account for tab's visual width
//         } else {
//             pos++
//         }
//     }
//
//     // Find end of the word
//     for pos < len(line) && line[pos] != ' ' && line[pos] != '\t' {
//         pos++
//     }
//
//     // Move back to last character of word (or stay if at end)
//     if pos > 0 && (pos == len(line) || line[pos] == ' ' || line[pos] == '\t') {
//         pos--
//     }
//
//     return pos
// }
