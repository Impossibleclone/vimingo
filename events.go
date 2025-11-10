package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	// Handle key press events
	case tea.KeyMsg:
		quit := func() {
			m.quit = true
		}

		screenH := m.height

		switch m.mode.Current() {

		//NormalMode
		case Normal:
			switch msg.Type {
			case tea.KeyCtrlD:
				m.cursor.HalfDown(m.buffer, m.height)
				adjustScroll(m.buffer, screenH)

			case tea.KeyCtrlU:
				m.cursor.HalfUp(m.buffer, m.height)
				adjustScroll(m.buffer, screenH)

			case tea.KeyDelete:
				m.buffer.Register = string(m.buffer.Lines[m.cursor.Y][m.cursor.X])
				m.buffer.Lines[m.cursor.Y] = RemoveCh(m.buffer.Lines[m.cursor.Y], m.cursor.X) //delete a character and update the line
				m.cursor.MoveLeft()

			case tea.KeyRunes:
				// Handle rune-based keys
				switch r := msg.Runes[0]; r {
				case 'v':
					m.visualStart = *m.cursor
					m.mode.SwitchTo(Visual)

				case 'p':
					Paste(m.buffer, m.buffer.Register)

				case 'i':
					m.mode.SwitchTo(Insert)

				case 'I':
					m.cursor.X = 0
					m.mode.SwitchTo(Insert)

				case 'a':
					m.mode.SwitchTo(Insert)
					m.cursor.MoveRightinInsert(m.buffer)

				case 'A':
					m.cursor.X = len(m.buffer.Lines[m.cursor.Y])
					m.mode.SwitchTo(Insert)

				case 'O':
					m.cursor.X = 0
					NewLine(m.buffer)
					m.cursor.Y--
					m.mode.SwitchTo(Insert)

				case 'o':
					m.mode.SwitchTo(Insert)
					m.cursor.X = len(m.buffer.Lines[m.cursor.Y])
					NewLine(m.buffer)

				case ':':
					m.mode.SwitchTo(Command)
					m.buffer.Command = nil

				case 'h':
					m.cursor.MoveLeft()
				case '^':
					m.cursor.X = 0

				case 'j':
					m.cursor.MoveDown(m.buffer)
					adjustScroll(m.buffer, screenH)

				case 'k':
					m.cursor.MoveUp(m.buffer)
					adjustScroll(m.buffer, screenH)

				case 'l':
					m.cursor.MoveRightinNormal(m.buffer)

				case '$':
					m.cursor.X = len(m.buffer.Lines[m.cursor.Y]) - 1

				case 'd':
					r := msg.Runes[0]
					if m.buffer.KeyReg == nil {
						m.buffer.KeyReg = append(m.buffer.KeyReg, r)
					} else {
						m.buffer.Lines[m.cursor.Y] = RemoveChs(m.buffer.Lines[m.cursor.Y], 0, len(m.buffer.Lines[m.cursor.Y]))
						m.buffer.StatusMsg = "line deleted"
						m.buffer.KeyReg = nil
					}
				case 'y':
					r := msg.Runes[0]
					if m.buffer.KeyReg == nil {
						m.buffer.KeyReg = append(m.buffer.KeyReg, r)
					} else {
						m.buffer.KeyReg = nil
					}
				case 'x':
					m.buffer.Register = string(m.buffer.Lines[m.cursor.Y][m.cursor.X])
					m.buffer.Lines[m.cursor.Y] = RemoveCh(m.buffer.Lines[m.cursor.Y], m.cursor.X) //delete a character and update the line
					m.cursor.MoveLeft()

				case 'w':
					if len(m.buffer.KeyReg) > 0 && m.buffer.KeyReg[0] == 'd' {
						m.buffer.StatusMsg = "word deleted"
					} else if len(m.buffer.KeyReg) > 0 && m.buffer.KeyReg[0] == 'y' {
						YankRange(m.buffer, m.cursor, wMotion(m.buffer, m.cursor))
						m.buffer.StatusMsg = fmt.Sprintf("yanked till %d", wMotion(m.buffer, m.cursor))
						m.buffer.KeyReg = nil
					} else {
						movedcur := wMotion(m.buffer, m.cursor)
						m.cursor.X = movedcur
						if m.cursor.X == len(m.buffer.Lines[m.cursor.Y]) {
							m.cursor.X--
						}
					}
				case 'e':
					if len(m.buffer.KeyReg) > 0 && m.buffer.KeyReg[0] == 'd' {
						m.buffer.StatusMsg = "deleted"
					} else if len(m.buffer.KeyReg) > 0 && m.buffer.KeyReg[0] == 'y' {
						YankRange(m.buffer, m.cursor, eMotion(m.buffer, m.cursor)+1)
						m.buffer.StatusMsg = fmt.Sprintf("yanked till %d", eMotion(m.buffer, m.cursor))
						m.buffer.KeyReg = nil
					} else {
						movedcur := eMotion(m.buffer, m.cursor)
						m.cursor.X = movedcur
						if m.cursor.X == len(m.buffer.Lines[m.cursor.Y]) {
							m.cursor.X--
						}
					}
				case 'E':
					if len(m.buffer.KeyReg) > 0 && m.buffer.KeyReg[0] == 'd' {
						m.buffer.StatusMsg = "deleted"
					} else if len(m.buffer.KeyReg) > 0 && m.buffer.KeyReg[0] == 'y' {
						YankRange(m.buffer, m.cursor, EMotion(m.buffer, m.cursor)+1)
						m.buffer.StatusMsg = fmt.Sprintf("yanked till %d", EMotion(m.buffer, m.cursor))
						m.buffer.KeyReg = nil
					} else {
						movedcur := EMotion(m.buffer, m.cursor)
						m.cursor.X = movedcur
						if m.cursor.X == len(m.buffer.Lines[m.cursor.Y]) {
							m.cursor.X--
						}
					}
				}
			}

		//insertMode
		case Insert:
			switch msg.Type {
			case tea.KeyBackspace:
				if m.cursor.X == 0 {
					RemoveLine(m.buffer)
				}
				if m.cursor.X > 0 {
					m.cursor.MoveLeft()
					m.buffer.Lines[m.cursor.Y] = RemoveCh(m.buffer.Lines[m.cursor.Y], m.cursor.X) //delete a character and update the line
				}

			case tea.KeyEnter:
				NewLine(m.buffer)

			case tea.KeyEscape, tea.KeyCtrlC:
				if m.cursor.X > len(m.buffer.Lines[m.cursor.Y])-1 {
					m.cursor.MoveLeft()
				}
				m.mode.SwitchTo(Normal)

			case tea.KeyLeft:
				m.cursor.MoveLeft()

			case tea.KeyDown:
				m.cursor.MoveDown(m.buffer)
				adjustScroll(m.buffer, screenH)

			case tea.KeyUp:
				m.cursor.MoveUp(m.buffer)
				adjustScroll(m.buffer, screenH)

			case tea.KeyRight:
				m.cursor.MoveRightinInsert(m.buffer)

			case tea.KeyRunes:
				//Logic for Typing:
				r := msg.Runes[0]                                                          //save the typed character
				m.buffer.Lines[m.cursor.Y] = TypeCh(m.buffer.Lines[m.cursor.Y], m.cursor.X, r) //update the line
				m.cursor.MoveRightinInsert(m.buffer)                                       //increment the position of the cursor in X.
			}
		case Visual:
			switch msg.Type {
			case tea.KeyEscape, tea.KeyCtrlC:
				m.mode.SwitchTo(Normal)

			case tea.KeyRunes:
				switch r := msg.Runes[0]; r {
				case 'v':
					if m.mode.Current() == Normal {
						m.mode.SwitchTo(Visual)
						m.visualStart = *m.cursor
					} else {
						m.mode.SwitchTo(Normal)
					}
				case 'y':
					YankSelection(m.buffer, m.cursor, &m.visualStart)
					m.mode.SwitchTo(Normal)

				case 'x':
					CutSelection(m.buffer, m.cursor, &m.visualStart)
					m.mode.SwitchTo(Normal)

				case 'c':
					CutSelection(m.buffer, m.cursor, &m.visualStart)
					m.mode.SwitchTo(Insert)

				case ':':
					m.mode.SwitchTo(Command)
					m.buffer.Command = nil
				case 'p':
					Paste(m.buffer, m.buffer.Register)

				case 'h':
					m.cursor.MoveLeft()

				case 'j':
					m.cursor.MoveDown(m.buffer)
					adjustScroll(m.buffer, screenH)

				case 'k':
					m.cursor.MoveUp(m.buffer)
					adjustScroll(m.buffer, screenH)

				case 'l':
					m.cursor.MoveRightinNormal(m.buffer)
				}
			}

		//CommandMode
		case Command:
			switch msg.Type {
			case tea.KeyBackspace:
				if len(m.buffer.Command) > 0 {
					m.buffer.Command = m.buffer.Command[:len(m.buffer.Command)-1]
				}

			case tea.KeyEnter:
				cmds := strings.Fields(string(m.buffer.Command))

				if len(cmds) == 0 {
					m.mode.SwitchTo(Normal)
					m.buffer.Command = nil
					return m, nil
				}

				switch cmds[0] {
				case "w":
					if len(cmds) > 1 {
						if err := SaveFile(cmds[1], m.buffer); err != nil {
							m.buffer.StatusMsg = err.Error()
						} else {
							m.buffer.Filename = cmds[1]
							m.buffer.StatusMsg = "Written " + m.buffer.Filename
						}
					} else {
						if err := SaveFile(m.buffer.Filename, m.buffer); err != nil {
							m.buffer.StatusMsg = err.Error()
						} else {
							m.buffer.StatusMsg = "Written " + m.buffer.Filename
						}
					}
					m.mode.SwitchTo(Normal)
					m.buffer.Command = nil

				case "q":
					quit()

				case "wq":
					if err := SaveFile(m.buffer.Filename, m.buffer); err != nil {
						m.buffer.StatusMsg = err.Error()
						m.mode.SwitchTo(Normal)
						m.buffer.Command = nil
					} else {
						m.buffer.StatusMsg = "Written " + m.buffer.Filename
						quit()
					}
				}
				m.buffer.Command = nil

			case tea.KeyEscape:
				m.mode.SwitchTo(Normal)

			case tea.KeyRunes:
				r := msg.Runes[0] //preserve the typed character
				m.buffer.Command = append(m.buffer.Command, r)
			}
		}
	}

	if m.quit {
		return m, tea.Quit
	}

	return m, nil
}
