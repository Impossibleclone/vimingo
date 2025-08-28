package main

import (
	"github.com/gdamore/tcell/v2"
	"log"
	"os"
)

func main() {
	// if len(os.Args) < 2 {
	// 	log.Fatal("Please read the instructions properly.")
	// }
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

	screen, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("Error creating screen: %v", err)
	}

	// Initialize the screen
	if err := screen.Init(); err != nil {
		log.Fatalf("Error initializing screen: %v", err)
	}

	defer screen.Fini()
	screen.Clear()
	screen.Show()

	// cur := &Cursor{X: 0, Y: 0}
	cursor := buffer.Cursor

	quit := func() {
		screen.Fini()
		os.Exit(0)
	}
	mode := NewEditorMode()
	for {

		ev := screen.PollEvent()
		HandleEvent(ev, buffer, cursor, &visualStart, mode, screen, quit)

		RenderScreen(screen, buffer, visualStart, mode)
	}
}
