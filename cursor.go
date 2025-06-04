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

func (c *Cursor) MoveRight(buffer *Buffer) {
	if c.Y >= len(buffer.Lines) {
		return
	}

	if c.X < len(buffer.Lines[c.Y])  {
		c.X++
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


