package main

type Cursor struct {
    X int
    Y int
}

// Define method on *Cursor receiver:
func (c *Cursor) MoveLeft() {
    if c.X > 0 {
        c.X--
    }
}

func (c *Cursor) MoveRightinNormal(buffer *Buffer) {
	if c.Y >= len(buffer.Lines) {
		return
	}

	if c.X < len(buffer.Lines[c.Y])-1  {
		c.X++
	}else if c.Y+1 < len(buffer.Lines){
		// c.Y++
		return
		// c.X = 0
	}
}
func (c *Cursor) MoveRightinInsert(buffer *Buffer) {
	if c.Y >= len(buffer.Lines) {
		return
	}

	if c.X < len(buffer.Lines[c.Y])  {
		c.X++
	}else if c.Y+1 < len(buffer.Lines){
		// c.Y++
		return
		// c.X = 0
	}
}

func (c *Cursor) MoveUp(buffer *Buffer) {
	if c.Y > 0 {
		c.Y--
		linelen := len(buffer.Lines[c.Y])
		if c.X > linelen {
			c.X = linelen
		}
	}
}

func (c *Cursor) MoveDown(buffer *Buffer) {
	if c.Y < len(buffer.Lines)-1 {
		c.Y++
		linelen := len(buffer.Lines[c.Y])
		if c.X > linelen{
			c.X = linelen
		}
	}
}


