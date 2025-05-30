package main

// define mode
type Mode int

const (
    Normal Mode = iota
    Insert
	Command
)

// currnt mode
type EditorMode struct {
    current Mode
}

// to start with normal
func NewEditorMode() *EditorMode {
    return &EditorMode{current: Normal}
}
// give current to show
func (e *EditorMode) Current() Mode {
    return e.current
}
// change the mode
func (e *EditorMode) SwitchTo(m Mode) {
    e.current = m
}


