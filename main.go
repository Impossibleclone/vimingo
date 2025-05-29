package main

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Are you Dumb!!!")
	}

	filename := os.Args[1]
	buf, err := LoadFile(filename)
	if err != nil {
		log.Fatalf("failed to open file: %v",err)
	}

	buffer := &Buffer{
		Lines: buf.Lines,
		Cursor: &Cursor{X:0,Y:0},
		Mode: Normal,
		}

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

	// MaxW, MaxH := screen.Size()
	screen.SetContent(0, 0, 'g', nil, tcell.StyleDefault)
	fmt.Println("UYS was here")
	quit := func() {
		screen.Fini()
		os.Exit(0)
	}
	mode := NewEditorMode()
	for {

		ev := screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventResize:
			screen.Sync()
			// MaxW, MaxH = screen.Size()
		case *tcell.EventKey:
			switch mode.Current() {
            case Normal:
                switch ev.Rune() {
                case 'i':
                    mode.SwitchTo(Insert)
                // switch ev.Rune() {
				case 'h':
					cursor.MoveLeft()
				case 'j':
					cursor.MoveDown(buffer)
				case 'k':
					cursor.MoveUp(buffer)
				case 'l':
					cursor.MoveRight(buffer)
				case 'q':
					quit()
				}
				// ... other Normal mode keys
            case Insert:
                if ev.Key() == tcell.KeyEscape {
                    mode.SwitchTo(Normal)
                }
                // ... insert mode keys
            
			}
			
		}
		screen.Clear()
		for y, line := range buf.Lines{
			for x, r := range line {
				screen.SetContent(x,y,r,nil,tcell.StyleDefault)
			}
		}

		screen.SetContent(cursor.X, cursor.Y, '█', nil, tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorBlue))
		screen.Show()
	}
}
