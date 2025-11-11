package app

import (
	"log"
	"os"

	"github.com/impossibleclone/vimingo/internal/core" // Imports the new core package
)

// Model holds the entire state of the text editor.
// It is capitalized to be exported.
type Model struct {
	buffer      *core.Buffer
	cursor      *core.Cursor
	visualStart core.Cursor
	mode        *core.EditorMode
	width       int
	height      int
	quit        bool
}

// InitialModel sets up the editor's starting state.
// It is capitalized to be exported.
func InitialModel() Model {
	var buffer *core.Buffer
	var filename string

	if len(os.Args) == 1 {
		filename = "[No Name]"
		buffer = &core.Buffer{
			Filename: filename,
			Lines:    []string{""}, // at least one line
			Cursor:   &core.Cursor{X: 0, Y: 0},
			Mode:     core.Normal,
			ScrollX:  0,
			ScrollY:  0,
			Register: "",
		}
	} else {
		filename := os.Args[1]
		// LoadFile is in the core package
		buf, err := core.LoadFile(filename)
		if err != nil {
			log.Fatalf("failed to open the file: %v", err)
		}

		if len(buf.Lines) == 0 {
			buf.Lines = append(buf.Lines, "") //we can't edit empty buffer err: index out of range [0] with length 0
		}

		buffer = &core.Buffer{
			Filename: filename,
			Lines:    buf.Lines,
			Cursor:   &core.Cursor{X: 0, Y: 0},
			Mode:     core.Normal,
			ScrollX:  0,
			ScrollY:  0,
			Register: "",
		}
	}

	visualStart := core.Cursor{X: 0, Y: 0} // default value
	cursor := buffer.Cursor
	mode := core.NewEditorMode() // NewEditorMode is in the core package

	return Model{
		buffer:      buffer,
		cursor:      cursor,
		visualStart: visualStart,
		mode:        mode,
	}
}
