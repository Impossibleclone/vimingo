package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gdamore/tcell/v2"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Please read the instructions properly.")
	}

	filename := os.Args[1]
	buf, err := LoadFile(filename)
	if err != nil {
		log.Fatalf("failed to open file: %v", err)
	}

	if len(buf.Lines) == 0 {
		buf.Lines = append(buf.Lines, "") //we can't edit empty buffer err: index out of range [0] with length 0
	}

	buffer := &Buffer{
		Filename: filename,
		Lines:    buf.Lines,
		Cursor:   &Cursor{X: 0, Y: 0},
		Mode:     Normal,
		ScrollX:  0,
		ScrollY:  0,
		Register: "",
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

	_, screenH := screen.Size()
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

			//NormalMode
			case Normal:
				switch ev.Rune() {
				//To switch to InsertMode
				case 'v':
					visualStart = *cursor
					mode.SwitchTo(Visual)

				case 'p':
					Paste(buffer, buffer.Register)

				case 'i':
					mode.SwitchTo(Insert)

				case 'I':
					cursor.X = 0
					mode.SwitchTo(Insert)

				//to edit after the cursor.
				case 'a':
					mode.SwitchTo(Insert)
					cursor.MoveRightinInsert(buffer)

				case 'A':
					cursor.X = len(buffer.Lines[cursor.Y])
					mode.SwitchTo(Insert)

				//to edit on a new line.
				case 'O':
					cursor.X = 0
					NewLine(buffer)
					cursor.Y--
					mode.SwitchTo(Insert)

				case 'o':
					mode.SwitchTo(Insert)
					cursor.X = len(buffer.Lines[cursor.Y])
					// cursor.MoveDown(buffer)
					NewLine(buffer)
					// cursor.MoveUp(buffer)

				case ':':
					mode.SwitchTo(Command)
					buffer.Command = nil

				case 'h':
					cursor.MoveLeft()

				case 'j':
					cursor.MoveDown(buffer)
					adjustScroll(buffer, screenH)

				case 'k':
					cursor.MoveUp(buffer)
					adjustScroll(buffer, screenH)

				case 'l':
					cursor.MoveRightinNormal(buffer)

					// case 'q':
					// 	quit()
				}

			//insertMode
			case Insert:
				switch {
				case ev.Key() == tcell.KeyBackspace, ev.Key() == tcell.KeyBackspace2:
					if cursor.X == 0 {
						RemoveLine(buffer)
					}
					if cursor.X > 0 {
						cursor.MoveLeft()
						buffer.Lines[cursor.Y] = RemoveCh(buffer.Lines[cursor.Y], cursor.X) //delete a character and update the line
					}

				case ev.Key() == tcell.KeyEnter, ev.Key() == tcell.KeyCR:
					NewLine(buffer)

				case ev.Key() == tcell.KeyEscape, ev.Key() == tcell.KeyCtrlC:
					mode.SwitchTo(Normal)

				//Logic for Typing:
				case ev.Rune() != 0:
					r := ev.Rune()                                                       //save the typed character
					buffer.Lines[cursor.Y] = TypeCh(buffer.Lines[cursor.Y], cursor.X, r) //update the line
					cursor.MoveRightinInsert(buffer)                                     //increment the position of the cursor in X.

				case ev.Key() == tcell.KeyLeft:
					cursor.MoveLeft()

				case ev.Key() == tcell.KeyDown:
					cursor.MoveDown(buffer)
					adjustScroll(buffer, screenH)

				case ev.Key() == tcell.KeyUp:
					cursor.MoveUp(buffer)
					adjustScroll(buffer, screenH)

				case ev.Key() == tcell.KeyRight:
					cursor.MoveRightinInsert(buffer)
				}
			case Visual:
				switch {
				case ev.Key() == tcell.KeyEscape, ev.Key() == tcell.KeyCtrlC:
					mode.SwitchTo(Normal)

				case ev.Rune() == 'v':
					if mode.Current() == Normal {
						mode.SwitchTo(Visual)
						visualStart = *cursor
					} else {
						mode.SwitchTo(Normal)
					}
				case ev.Rune() == 'y':
					start := min(visualStart.X, cursor.X)
					end := max(visualStart.X, cursor.X)
					startline := min(visualStart.Y, cursor.Y)
					endline := max(visualStart.Y, cursor.Y)
					buffer.Register = ""
					if startline == endline {

						toYankFromLine := buffer.Lines[startline]
						toYankTheCharacters := []rune(toYankFromLine[start:end])
						buffer.Register = string(toYankTheCharacters)
					} else {
						for y := startline; y <= endline; y++ {
							if y == startline {
								toYankFromLine := buffer.Lines[y]
								toYankTheCharacters := []rune(toYankFromLine[start:])
								buffer.Register += string(toYankTheCharacters)
								buffer.Register += "\n"

							} else if y == endline {
								toYankFromLine := buffer.Lines[y]
								toYankTheCharacters := []rune(toYankFromLine[:end])
								buffer.Register += string(toYankTheCharacters)

							} else {
								toYankFromLine := buffer.Lines[y]
								toYankTheCharacters := []rune(toYankFromLine[:])
								buffer.Register += string(toYankTheCharacters)
								buffer.Register += "\n"
							}
						}
					}
					cursor.Y = startline
					cursor.X = start
					mode.SwitchTo(Normal)

				case ev.Key() == ':':
					mode.SwitchTo(Command)

					buffer.Command = nil
				case ev.Rune() == 'p':
					Paste(buffer, buffer.Register)

				case ev.Rune() == 'h':
					cursor.MoveLeft()

				case ev.Rune() == 'j':
					cursor.MoveDown(buffer)
					adjustScroll(buffer, screenH)

				case ev.Rune() == 'k':
					cursor.MoveUp(buffer)
					adjustScroll(buffer, screenH)

				case ev.Rune() == 'l':
					cursor.MoveRightinNormal(buffer)
					// line := buffer.Lines[cursor.Y]
					// screen.SetContent(cursor.X, cursor.Y-buffer.ScrollY, rune(line[cursor.X]), nil, tcell.StyleDefault.Reverse(true))
				}

			//CommandMode

			case Command:
				switch {
				case ev.Key() == tcell.KeyEnter, ev.Key() == tcell.KeyCR:

					cmd := string(buffer.Command)

					switch cmd {
					case "w":
						SaveFile(filename, buffer)
						mode.SwitchTo(Normal)
						buffer.Command = nil

					case "q":
						quit()

					case "wq":
						SaveFile(filename, buffer)
						quit()
					}
					buffer.Command = nil

				case ev.Key() == tcell.KeyEscape:
					mode.SwitchTo(Normal)

				case ev.Rune() != 0:
					r := ev.Rune() //preserve the typed character
					buffer.Command = append(buffer.Command, r)

				}
			}

		}
		screen.Clear()

		for y := 0; y < screenH; y++ {
			lineIndex := y + buffer.ScrollY
			if lineIndex >= len(buffer.Lines) {
				break
			}
			line := buffer.Lines[lineIndex]
			// for x, r := range line {
			// 	screen.SetContent(x, y, r, nil, tcell.StyleDefault)
			// }
			for x, r := range line {
				style := tcell.StyleDefault

				if mode.Current() == Visual && isInSelection(visualStart, *cursor, x, lineIndex) {
					style = style.Reverse(true)
				}

				screen.SetContent(x, y, r, nil, style)
			}

		}

		if mode.Current() == Command {
			for i, r := range buffer.Command {
				screen.SetContent(i, screenH-1, r, nil, tcell.StyleDefault)
			}
		}

		// screen.SetContent(cursor.X, cursor.Y-buffer.ScrollY, '█', nil, tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorBlue))
		screen.ShowCursor(cursor.X, cursor.Y-buffer.ScrollY)
		screen.Show()
	}
}
