package main

import (
	"fmt"
	"strings"

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
			case '_':
				cursor.X = 0

			case 'j':
				cursor.MoveDown(buffer)
				adjustScroll(buffer, screenH)

			case 'k':
				cursor.MoveUp(buffer)
				adjustScroll(buffer, screenH)

			case 'l':
				cursor.MoveRightinNormal(buffer)

			case '$':
				cursor.X = len(buffer.Lines[cursor.Y]) - 1

			case 'd':
				r := ev.Rune()
				if buffer.KeyReg == nil {
					buffer.KeyReg = append(buffer.KeyReg, r)
				} else {
					buffer.Lines[cursor.Y] = RemoveChs(buffer.Lines[cursor.Y], 0, len(buffer.Lines[cursor.Y]))
					buffer.StatusMsg = "line deleted"
					buffer.KeyReg = nil
				}
			case 'y':
				r := ev.Rune()
				if buffer.KeyReg == nil {
					buffer.KeyReg = append(buffer.KeyReg, r)
				} else {
					buffer.KeyReg = nil
				}
			case 'x':
				buffer.Register = string(buffer.Lines[cursor.Y][cursor.X])
				buffer.Lines[cursor.Y] = RemoveCh(buffer.Lines[cursor.Y], cursor.X) //delete a character and update the line
				cursor.MoveLeft()

			case 'w':
				// r := ev.Rune()
				// buffer.KeyReg = append(buffer.KeyReg, r)
				if len(buffer.KeyReg) > 0 && buffer.KeyReg[0] == 'd' {
					buffer.StatusMsg = "word deleted"
				} else if len(buffer.KeyReg) > 0 && buffer.KeyReg[0] == 'y' {
					// start := cursor.X
					YankRange(buffer, cursor, wMotion(buffer, cursor))
					buffer.StatusMsg = fmt.Sprintf("yanked till %d", wMotion(buffer, cursor))
					buffer.KeyReg = nil
				} else {
					movedcur := wMotion(buffer, cursor)
					cursor.X = movedcur
					if cursor.X == len(buffer.Lines[cursor.Y]) {
						cursor.X--
					}
				}
			case 'e':
				// r := ev.Rune()
				// buffer.KeyReg = append(buffer.KeyReg, r)
				if len(buffer.KeyReg) > 0 && buffer.KeyReg[0] == 'd' {
					buffer.StatusMsg = "deleted"
				} else if len(buffer.KeyReg) > 0 && buffer.KeyReg[0] == 'y' {
					// start := cursor.X
					YankRange(buffer, cursor, eMotion(buffer, cursor)+1)
					buffer.StatusMsg = fmt.Sprintf("yanked till %d", eMotion(buffer, cursor))
					buffer.KeyReg = nil
				} else {
					movedcur := eMotion(buffer, cursor)
					cursor.X = movedcur
					if cursor.X == len(buffer.Lines[cursor.Y]) {
						cursor.X--
					}
				}
			case 'E':
				// r := ev.Rune()
				// buffer.KeyReg = append(buffer.KeyReg, r)
				if len(buffer.KeyReg) > 0 && buffer.KeyReg[0] == 'd' {
					buffer.StatusMsg = "deleted"
				} else if len(buffer.KeyReg) > 0 && buffer.KeyReg[0] == 'y' {
					// start := cursor.X
					YankRange(buffer, cursor, EMotion(buffer, cursor)+1)
					buffer.StatusMsg = fmt.Sprintf("yanked till %d", EMotion(buffer, cursor))
					buffer.KeyReg = nil
				} else {
					movedcur := EMotion(buffer, cursor)
					cursor.X = movedcur
					if cursor.X == len(buffer.Lines[cursor.Y]) {
						cursor.X--
					}
				}
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
				YankSelection(buffer, cursor, visualStart)
				mode.SwitchTo(Normal)

			case ev.Rune() == 'x':
				CutSelection(buffer, cursor, visualStart)
				mode.SwitchTo(Normal)

			case ev.Rune() == 'c':
				CutSelection(buffer, cursor, visualStart)
				mode.SwitchTo(Insert)

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

				cmds := strings.Fields(string(buffer.Command))

				//check if nothing is typed
				if len(cmds) == 0 {
					mode.SwitchTo(Normal)
					buffer.Command = nil
					return
				}

				switch cmds[0] {
				case "w":
					if len(cmds) > 1 {
						if err := SaveFile(cmds[1], buffer); err != nil {
							buffer.StatusMsg = err.Error()
						} else {
							buffer.Filename = cmds[1]
							buffer.StatusMsg = "Written " + buffer.Filename
						}
					} else {
						if err := SaveFile(buffer.Filename, buffer); err != nil {
							buffer.StatusMsg = err.Error()
						} else {
							buffer.StatusMsg = "Written " + buffer.Filename
						}
					}

					mode.SwitchTo(Normal)
					buffer.Command = nil

				case "q":
					quit()

				case "wq":
					if err := SaveFile(buffer.Filename, buffer); err != nil {
						buffer.StatusMsg = err.Error()
						mode.SwitchTo(Normal)
						buffer.Command = nil
					} else {
						buffer.StatusMsg = "Written " + buffer.Filename
						quit()
					}
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
