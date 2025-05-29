package main

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"log"
	"os"
)

func main() {
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

	cur := &Cursor{X: 0, Y: 0}
	MaxW, MaxH := screen.Size()
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
			MaxW, MaxH = screen.Size()
		case *tcell.EventKey:
			switch mode.Current() {
            case Normal:
                switch ev.Rune() {
                case 'i':
                    mode.SwitchTo(Insert)
                // switch ev.Rune() {
				case 'h':
					cur.MoveLeft()
				case 'j':
					cur.MoveDown(MaxH)
				case 'k':
					cur.MoveUp()
				case 'l':
					cur.MoveRight(MaxW)
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
		screen.SetContent(cur.X, cur.Y, '█', nil, tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorBlue))
		screen.Show()
	}
}
