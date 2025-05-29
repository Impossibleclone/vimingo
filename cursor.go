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

func (c *Cursor) MoveRight(MaxW int) {
	if c.X < MaxW -1  {
		c.X++
	}
}

func (c *Cursor) MoveUp() {
	if c.Y > 0 {
		c.Y--
	}
}

func (c *Cursor) MoveDown(MaxH int) {
	if c.Y < MaxH-1 {
		c.Y++
	}
}



