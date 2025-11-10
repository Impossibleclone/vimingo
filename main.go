package main

import (
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

// model holds the entire state of the text editor.
type model struct {
	buffer      *Buffer
	cursor      *Cursor
	visualStart Cursor
	mode        *EditorMode
	width       int
	height      int
	quit        bool
}

// initialModel sets up the editor's starting state.
func initialModel() model {
	var buffer *Buffer
	var filename string

	if len(os.Args) == 1 {
		filename = "[No Name]"
		buffer = &Buffer{
			Filename: filename,
			Lines:    []string{""}, // at least one line
			Cursor:   &Cursor{X: 0, Y: 0},
			Mode:     Normal,
			ScrollX:  0,
			ScrollY:  0,
			Register: "",
		}
	} else {
		filename := os.Args[1]
		buf, err := LoadFile(filename)
		if err != nil {
			log.Fatalf("failed to open the file: %v", err)
		}

		if len(buf.Lines) == 0 {
			buf.Lines = append(buf.Lines, "") //we can't edit empty buffer err: index out of range [0] with length 0
		}

		buffer = &Buffer{
			Filename: filename,
			Lines:    buf.Lines,
			Cursor:   &Cursor{X: 0, Y: 0},
			Mode:     Normal,
			ScrollX:  0,
			ScrollY:  0,
			Register: "",
		}
	}

	visualStart := Cursor{X: 0, Y: 0} // default value
	cursor := buffer.Cursor
	mode := NewEditorMode()

	return model{
		buffer:      buffer,
		cursor:      cursor,
		visualStart: visualStart,
		mode:        mode,
	}
}

func main() {
	m := initialModel()
	// WithAltScreen to have full screen
	p := tea.NewProgram(&m, tea.WithAltScreen(), tea.WithMouseAllMotion())

	if _, err := p.Run(); err != nil {
		log.Fatalf("Error running program: %v", err)
	}
}
