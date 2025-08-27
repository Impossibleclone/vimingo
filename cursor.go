package main

type Cursor struct {
	X    int
	Y    int
	remX int
}

// Define method on *Cursor receiver:
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
		} else if c.X > linelen || c.remX > linelen || c.remX == linelen{
			c.X = linelen-1
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
		} else if c.X > linelen || c.remX > linelen || c.remX == linelen{
			c.X = linelen-1
		} else {
			c.X = c.remX
		}
	}
}
