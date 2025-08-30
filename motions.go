package main

func isAlphabetorNumber(r rune) bool {
	return (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9')
}
func wMotion(buffer *Buffer, cursor *Cursor) int {
	line := buffer.Lines[cursor.Y]
	var pos int
	if pos+1 < len(line) {
		pos = cursor.X + 1
	} else {
		pos = len(line) - 1
	}
	// Skip All Alphabets
	for pos < len(line) && isAlphabetorNumber(rune(line[pos])) {
		pos++
	}
	// Skip spaces to reach start of next word
	for pos < len(line) && line[pos] == ' ' {
		pos++
	}
	if pos >= len(line) {
		return len(line) - 1
	} else {
		return pos
	}
}

const TabStop = 4 // Define tab width for cursor movement consistency

func eMotion(buffer *Buffer, cursor *Cursor) int {
	line := buffer.Lines[cursor.Y]
	pos := cursor.X

	//if end of word go to the space if next is space
	if pos < len(line)-1 && line[pos+1] == ' ' {
		pos++
	}

	//if end of word go to tab if next is tab
	if pos < len(line)-1 && line[pos+1] == '\t' {
		pos += TabStop
	}

	//if on alphabet go to last alphabet
	if pos < len(line)-1 && isAlphabetorNumber(rune(line[pos])) {
		for pos < len(line)-1 && isAlphabetorNumber(rune(line[pos+1])) {
			pos++
		}
		if pos != cursor.X {
			return pos
		}
		//if after spaces and tabs there is a symbol go to it.
		if !isAlphabetorNumber(rune(line[pos+1])){
			return pos+1
		}
	}
	if pos < len(line)-1 && !isAlphabetorNumber(rune(line[pos])) {
		//if on space or tab 
		if line[pos] == ' ' || line[pos] == '\t' {

			//for all spaces one after another
			for pos < len(line)-1 && line[pos] == ' ' {
				pos++
			}

			//for all tabs one after another
			for pos < len(line)-1 && line[pos] == '\t' {
				pos += TabStop
			}

			//if after all spaces and tabs a symbol is found next give that position.
			if !isAlphabetorNumber(rune(line[pos])) {
				// buffer.StatusMsg = fmt.Sprintf("this")
				return pos
			}

			//if after all spaces and tabs an alphabet is found go to all alphabet till last one i.e end
			for pos < len(line)-1 && isAlphabetorNumber(rune(line[pos+1])) {
				pos++
			}
			return pos
		} else if pos < len(line)-1 && !isAlphabetorNumber(rune(line[pos+1])) {
			return pos + 1
		} else if pos < len(line)-1 && isAlphabetorNumber(rune(line[pos+1])) {
			for pos < len(line)-1 && isAlphabetorNumber(rune(line[pos+1])) {
				pos++
			}
			if pos != cursor.X {
				return pos
			}
		}
	}

	return pos

}
