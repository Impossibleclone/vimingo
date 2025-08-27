package main

import (
	"github.com/gdamore/tcell/v2"
)

func HandleEvent(ev tcell.Event, buffer *Buffer, cursor *Cursor, visualStart *Cursor, mode *EditorMode, screen tcell.Screen, quit func()) {
	switch ev := ev.(type) {
	case *tcell.EventResize:
		screen.Sync()
		// MaxW, MaxH = screen.Size()

	case *tcell.EventKey:
		_, screenH := screen.Size()
		switch mode.Current() {

		//NormalMode
		case Normal:
			switch ev.Key() {
			case tcell.KeyCtrlD:
				cursor.HalfDown(buffer, screen)
				adjustScroll(buffer, screenH)

			case tcell.KeyCtrlU:
				cursor.HalfUp(buffer, screen)
				adjustScroll(buffer, screenH)

			}
			switch ev.Rune() {
			//To switch to InsertMode
			case 'v':
				*visualStart = *cursor
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
				if cursor.X > len(buffer.Lines[cursor.Y])-1 {
					cursor.MoveLeft()
				}
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
					*visualStart = *cursor
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
					toYankTheCharacters := []rune(toYankFromLine[start : end+1])
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
			case ev.Key() == tcell.KeyBackspace, ev.Key() == tcell.KeyBackspace2:
				if len(buffer.Command) > 0 {
					buffer.Command = buffer.Command[:len(buffer.Command)-1]
				}

			case ev.Key() == tcell.KeyEnter, ev.Key() == tcell.KeyCR:

				cmd := string(buffer.Command)

				switch cmd {
				case "w":
					SaveFile(buffer.Filename, buffer)
					mode.SwitchTo(Normal)
					buffer.Command = nil

				case "q":
					quit()

				case "wq":
					SaveFile(buffer.Filename, buffer)
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
}
