package app

import (
	"fmt"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/impossibleclone/vimingo/internal/core" // Imports the new core package
)

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
		case core.Normal: // Prefixed
			switch msg.Type {
			case tea.KeyCtrlD:
				m.cursor.HalfDown(m.buffer, m.height)
				core.AdjustScroll(m.buffer, screenH) // Prefixed
			case tea.KeyCtrlU:
				m.cursor.HalfUp(m.buffer, m.height)
				core.AdjustScroll(m.buffer, screenH) // Prefixed
			case tea.KeyDelete:
				m.buffer.Register = string(m.buffer.Lines[m.cursor.Y][m.cursor.X])
				m.buffer.Lines[m.cursor.Y] = core.RemoveCh(m.buffer.Lines[m.cursor.Y], m.cursor.X) // Prefixed
				m.cursor.MoveLeft()

			case tea.KeyRunes:
				// Handle rune-based keys
				switch r := msg.Runes[0]; r {
				case 'v':
					m.visualStart = *m.cursor
					m.mode.SwitchTo(core.Visual) // Prefixed
				case 'p':
					core.Paste(m.buffer, m.buffer.Register) // Prefixed
				case 'i':
					m.mode.SwitchTo(core.Insert) // Prefixed
				case 'I':
					m.cursor.X = 0
					m.mode.SwitchTo(core.Insert) // Prefixed
				case 'a':
					m.mode.SwitchTo(core.Insert) // Prefixed
					m.cursor.MoveRightinInsert(m.buffer)
				case 'A':
					m.cursor.X = len(m.buffer.Lines[m.cursor.Y])
					m.mode.SwitchTo(core.Insert) // Prefixed
				case 'O':
					m.cursor.X = 0
					core.NewLine(m.buffer) // Prefixed
					m.cursor.Y--
					m.mode.SwitchTo(core.Insert) // Prefixed
				case 'o':
					m.mode.SwitchTo(core.Insert) // Prefixed
					m.cursor.X = len(m.buffer.Lines[m.cursor.Y])
					core.NewLine(m.buffer) // Prefixed
				case ':':
					m.mode.SwitchTo(core.Command) // Prefixed
					m.buffer.Command = nil
				case 'h':
					m.cursor.MoveLeft()
				case '^':
					m.cursor.X = 0
				case 'j':
					m.cursor.MoveDown(m.buffer)
					core.AdjustScroll(m.buffer, screenH) // Prefixed
				case 'k':
					m.cursor.MoveUp(m.buffer)
					core.AdjustScroll(m.buffer, screenH) // Prefixed
				case 'l':
					m.cursor.MoveRightinNormal(m.buffer)
				case '$':
					m.cursor.X = len(m.buffer.Lines[m.cursor.Y]) - 1
				case 'd':
					r := msg.Runes[0]
					if m.buffer.KeyReg == nil {
						m.buffer.KeyReg = append(m.buffer.KeyReg, r)
					} else {
						m.buffer.Lines[m.cursor.Y] = core.RemoveChs(m.buffer.Lines[m.cursor.Y], 0, len(m.buffer.Lines[m.cursor.Y])) // Prefixed
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
					m.buffer.Lines[m.cursor.Y] = core.RemoveCh(m.buffer.Lines[m.cursor.Y], m.cursor.X) // Prefixed
					m.cursor.MoveLeft()
				case 'w':
					if len(m.buffer.KeyReg) > 0 && m.buffer.KeyReg[0] == 'd' {
						m.buffer.StatusMsg = "word deleted"
					} else if len(m.buffer.KeyReg) > 0 && m.buffer.KeyReg[0] == 'y' {
						core.YankRange(m.buffer, m.cursor, core.WMotion(m.buffer, m.cursor)) // Prefixed
						m.buffer.StatusMsg = fmt.Sprintf("yanked till %d", core.WMotion(m.buffer, m.cursor))
						m.buffer.KeyReg = nil
					} else {
						movedcur := core.WMotion(m.buffer, m.cursor) // Prefixed
						m.cursor.X = movedcur
						if m.cursor.X == len(m.buffer.Lines[m.cursor.Y]) {
							m.cursor.X--
						}
					}
				case 'e':
					if len(m.buffer.KeyReg) > 0 && m.buffer.KeyReg[0] == 'd' {
						m.buffer.StatusMsg = "deleted"
					} else if len(m.buffer.KeyReg) > 0 && m.buffer.KeyReg[0] == 'y' {
						core.YankRange(m.buffer, m.cursor, core.EMotion(m.buffer, m.cursor)+1) // Prefixed
						m.buffer.StatusMsg = fmt.Sprintf("yanked till %d", core.EMotion(m.buffer, m.cursor))
						m.buffer.KeyReg = nil
					} else {
						movedcur := core.EMotion(m.buffer, m.cursor) // Prefixed
						m.cursor.X = movedcur
						if m.cursor.X == len(m.buffer.Lines[m.cursor.Y]) {
							m.cursor.X--
						}
					}
				case 'E':
					if len(m.buffer.KeyReg) > 0 && m.buffer.KeyReg[0] == 'd' {
						m.buffer.StatusMsg = "deleted"
					} else if len(m.buffer.KeyReg) > 0 && m.buffer.KeyReg[0] == 'y' {
						core.YankRange(m.buffer, m.cursor, core.CapitalEMotion(m.buffer, m.cursor)+1) // Prefixed
						m.buffer.StatusMsg = fmt.Sprintf("yanked till %d", core.CapitalEMotion(m.buffer, m.cursor))
						m.buffer.KeyReg = nil
					} else {
						movedcur := core.CapitalEMotion(m.buffer, m.cursor) // Prefixed
						m.cursor.X = movedcur
						if m.cursor.X == len(m.buffer.Lines[m.cursor.Y]) {
							m.cursor.X--
						}
					}
				}
			}

		//insertMode
		case core.Insert: // Prefixed
			switch msg.Type {
			case tea.KeySpace:
				m.buffer.Lines[m.cursor.Y] = core.TypeCh(m.buffer.Lines[m.cursor.Y], m.cursor.X, ' ') // Prefixed
				m.cursor.MoveRightinInsert(m.buffer)
			case tea.KeyTab:
				m.buffer.Lines[m.cursor.Y] = core.TypeCh(m.buffer.Lines[m.cursor.Y], m.cursor.X, '\t') // Prefixed
				m.cursor.MoveRightinInsert(m.buffer)
			case tea.KeyBackspace:
				if m.cursor.X == 0 {
					core.RemoveLine(m.buffer) // Prefixed
				}
				if m.cursor.X > 0 {
					m.cursor.MoveLeft()
					m.buffer.Lines[m.cursor.Y] = core.RemoveCh(m.buffer.Lines[m.cursor.Y], m.cursor.X) // Prefixed
				}
			case tea.KeyEnter:
				core.NewLine(m.buffer) // Prefixed
			case tea.KeyEscape, tea.KeyCtrlC:
				if m.cursor.X > len(m.buffer.Lines[m.cursor.Y])-1 {
					m.cursor.MoveLeft()
				}
				m.mode.SwitchTo(core.Normal) // Prefixed
			case tea.KeyLeft:
				m.cursor.MoveLeft()
			case tea.KeyDown:
				m.cursor.MoveDown(m.buffer)
				core.AdjustScroll(m.buffer, screenH) // Prefixed
			case tea.KeyUp:
				m.cursor.MoveUp(m.buffer)
				core.AdjustScroll(m.buffer, screenH) // Prefixed
			case tea.KeyRight:
				m.cursor.MoveRightinInsert(m.buffer)
			case tea.KeyRunes:
				//Logic for Typing:
				r := msg.Runes[0]
				m.buffer.Lines[m.cursor.Y] = core.TypeCh(m.buffer.Lines[m.cursor.Y], m.cursor.X, r) // Prefixed
				m.cursor.MoveRightinInsert(m.buffer)
			}
		case core.Visual: // Prefixed
			switch msg.Type {
			case tea.KeyEscape, tea.KeyCtrlC:
				m.mode.SwitchTo(core.Normal) // Prefixed
			case tea.KeyRunes:
				switch r := msg.Runes[0]; r {
				case 'v':
					if m.mode.Current() == core.Normal { // Prefixed
						m.mode.SwitchTo(core.Visual) // Prefixed
						m.visualStart = *m.cursor
					} else {
						m.mode.SwitchTo(core.Normal) // Prefixed
					}
				case 'y':
					core.YankSelection(m.buffer, m.cursor, &m.visualStart) // Prefixed
					m.mode.SwitchTo(core.Normal)                           // Prefixed
				case 'x':
					core.CutSelection(m.buffer, m.cursor, &m.visualStart) // Prefixed
					m.mode.SwitchTo(core.Normal)                         // Prefixed
				case 'c':
					core.CutSelection(m.buffer, m.cursor, &m.visualStart) // Prefixed
					m.mode.SwitchTo(core.Insert)                         // Prefixed
				case ':':
					m.mode.SwitchTo(core.Command) // Prefixed
					m.buffer.Command = nil
				case 'p':
					core.Paste(m.buffer, m.buffer.Register) // Prefixed
				case 'h':
					m.cursor.MoveLeft()
				case 'j':
					m.cursor.MoveDown(m.buffer)
					core.AdjustScroll(m.buffer, screenH) // Prefixed
				case 'k':
					m.cursor.MoveUp(m.buffer)
					core.AdjustScroll(m.buffer, screenH) // Prefixed
				case 'l':
					m.cursor.MoveRightinNormal(m.buffer)
				}
			}

		//CommandMode
		case core.Command: // Prefixed
			switch msg.Type {
			case tea.KeyBackspace:
				if len(m.buffer.Command) > 0 {
					m.buffer.Command = m.buffer.Command[:len(m.buffer.Command)-1]
				}

			case tea.KeyEnter:
				cmds := strings.Fields(string(m.buffer.Command))

				if len(cmds) == 0 {
					m.mode.SwitchTo(core.Normal) // Prefixed
					m.buffer.Command = nil
					return m, nil
				}

				if num, err := strconv.Atoi(cmds[0]); err == nil {
					if num > len(m.buffer.Lines) {
						m.cursor.Y = len(m.buffer.Lines) - 1
						m.mode.SwitchTo(core.Normal) // Prefixed
						m.buffer.Command = nil
						core.AdjustScroll(m.buffer, screenH) // Prefixed
						return m, nil
					}
					m.cursor.Y = num - 1
					m.mode.SwitchTo(core.Normal) // Prefixed
					m.buffer.Command = nil
					core.AdjustScroll(m.buffer, screenH) // Prefixed
					return m, nil
				}

				switch cmds[0] {
				case "w":
					if len(cmds) > 1 {
						if err := core.SaveFile(cmds[1], m.buffer); err != nil { // Prefixed
							m.buffer.StatusMsg = err.Error()
						} else {
							m.buffer.Filename = cmds[1]
							m.buffer.StatusMsg = "Written " + m.buffer.Filename
						}
					} else {
						if err := core.SaveFile(m.buffer.Filename, m.buffer); err != nil { // Prefixed
							m.buffer.StatusMsg = err.Error()
						} else {
							m.buffer.StatusMsg = "Written " + m.buffer.Filename
						}
					}
					m.mode.SwitchTo(core.Normal) // Prefixed
					m.buffer.Command = nil

				case "q":
					quit()

				case "wq":
					if err := core.SaveFile(m.buffer.Filename, m.buffer); err != nil { // Prefixed
						m.buffer.StatusMsg = err.Error()
						m.mode.SwitchTo(core.Normal) // Prefixed
						m.buffer.Command = nil
					} else {
						m.buffer.StatusMsg = "Written " + m.buffer.Filename
						quit()
					}
				}
				m.buffer.Command = nil

			case tea.KeyEscape:
				m.mode.SwitchTo(core.Normal) // Prefixed
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
