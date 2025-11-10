package main

type Cursor struct {
	X    int
	Y    int
	remX int
}

func (c *Cursor) MoveLeft() {
	if c.X > 0 {
		c.X--
	}
	c.remX = c.X
}

func (c *Cursor) MoveRightinNormal(buffer *Buffer) {
	if c.Y >= len(buffer.Lines) {
		return
	}

	if c.X < len(buffer.Lines[c.Y])-1 {
		c.X++
		c.remX = c.X
	} else if c.Y+1 < len(buffer.Lines) {
		// c.Y++
		return
		// c.X = 0
	}
}
func (c *Cursor) MoveRightinInsert(buffer *Buffer) {
	if c.Y >= len(buffer.Lines) {
		return
	}

	if c.X < len(buffer.Lines[c.Y]) {
		c.X++
		c.remX = c.X
	} else if c.Y+1 < len(buffer.Lines) {
		// c.Y++
		return
		// c.X = 0
	}
}

func (c *Cursor) MoveUp(buffer *Buffer) {
	if c.Y > 0 {
		c.Y--
		linelen := len(buffer.Lines[c.Y])
		if linelen == 0 {
			c.X = linelen
		} else if c.X > linelen || c.remX > linelen || c.remX == linelen {
			c.X = linelen - 1
		} else {
			c.X = c.remX
		}
	}
}

func (c *Cursor) MoveDown(buffer *Buffer) {
	if c.Y < len(buffer.Lines)-1 {
		c.Y++
		linelen := len(buffer.Lines[c.Y])
		if linelen == 0 {
			c.X = linelen
		} else if c.X > linelen || c.remX > linelen || c.remX == linelen {
			c.X = linelen - 1
		} else {
			c.X = c.remX
		}
	}
}

func (c *Cursor) HalfDown(buffer *Buffer, screenH int) {
	if c.Y < len(buffer.Lines)-1 {
		// _, screenH := screen.Size() // No longer need this line
		half := screenH / 2
		c.Y += half

		if c.Y >= len(buffer.Lines) {
			c.Y = len(buffer.Lines) - 1
		}

		linelen := len(buffer.Lines[c.Y])
		if linelen == 0 {
			c.X = linelen
		} else if c.X > linelen || c.remX > linelen || c.remX == linelen {
			c.X = linelen - 1
		} else {
			c.X = c.remX
		}
	}
}

func (c *Cursor) HalfUp(buffer *Buffer, screenH int) {
	// _, screenH := screen.Size() // No longer need this line
	half := screenH / 2
	if c.Y > 0 {
		if c.Y < half {
			c.Y = 0
		} else {
			c.Y -= half
		}

		if c.Y >= len(buffer.Lines) {
			c.Y = len(buffer.Lines) - 1
		}

		linelen := len(buffer.Lines[c.Y])
		if linelen == 0 {
			c.X = linelen
		} else if c.X > linelen || c.remX > linelen || c.remX == linelen {
			c.X = linelen - 1
		} else {
			c.X = c.remX
		}
	}
}

func visualColumn(line []rune, runeIndex int, tabSize int) int {
	col := 0
	for i := 0; i < runeIndex && i < len(line); i++ {
		if line[i] == '\t' {
			spaces := tabSize - (col % tabSize)
			col += spaces
		} else {
			col++
		}
	}
	return col
}
