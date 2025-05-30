package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gdamore/tcell/v2"
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
		Filename: filename,	
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
                case ':':
					mode.SwitchTo(Command)
					buffer.Command = nil
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
			case Insert:
				switch {
				case ev.Key() == tcell.KeyEnter , ev.Key() == tcell.KeyCR:
					NewLine(buffer)

				case ev.Rune() != 0 :
					r := ev.Rune() //save the typed character
					buffer.Lines[cursor.Y] = TypeCh(buffer.Lines[cursor.Y],cursor.X,r) //update the line 
					cursor.MoveRight(buffer) //increment the position of the cursor in X.
				
				case ev.Key() == tcell.KeyBackspace, ev.Key() == tcell.KeyBackspace2:
					if cursor.X == 0 {
						RemoveLine(buffer)
					}
					if cursor.X > 0 {	
					cursor.MoveLeft()
						buffer.Lines[cursor.Y] = RemoveCh(buffer.Lines[cursor.Y],cursor.X) //delete a character and update the line
					}
				
				case ev.Key() == tcell.KeyEscape :
					mode.SwitchTo(Normal)
				}
				
			case Command:
				switch {
				case ev.Key() == tcell.KeyEnter , ev.Key() == tcell.KeyCR:
					
					cmd := string(buffer.Command)

					switch cmd{
					case "w" :
						fmt.Println("Saving...")
						SaveFile(filename, buffer)
						mode.SwitchTo(Normal)
						buffer.Command = nil

					case "q" :
						quit()

					case "wq" :
						SaveFile(filename, buffer)
						quit()
					}
					buffer.Command = nil

				case ev.Key() == tcell.KeyEscape :
					mode.SwitchTo(Normal)

				case ev.Rune() != 0 :
					r:= ev.Rune() //preserve the typed character
					buffer.Command = append(buffer.Command, r)
				
				
				}
			}

		}
		screen.Clear()
		for y, line := range buffer.Lines{
			for x, r := range line {
				screen.SetContent(x,y,r,nil,tcell.StyleDefault)
			}
		}

		for i, r := range buffer.Command {
			screen.SetContent(i, buffer.Cursor.Y+1, r, nil, tcell.StyleDefault)
		}
		screen.SetContent(cursor.X, cursor.Y, '█', nil, tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorBlue))
		screen.Show()
	}
}
