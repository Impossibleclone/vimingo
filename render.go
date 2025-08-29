package main

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
)

func RenderScreen(screen tcell.Screen, buffer *Buffer, visualStart Cursor, mode *EditorMode) {
	screen.Clear()

	screenW, screenH := screen.Size()

	for y := 0; y < screenH-1; y++ {
		lineIndex := y + buffer.ScrollY
		if lineIndex >= len(buffer.Lines) {
			break
		}
		runes := []rune(buffer.Lines[lineIndex])
		x := 0 // screen column

		for rn := 0; rn < len(runes); rn++ {
			r := runes[rn]
			style := tcell.StyleDefault
			tab := ' '
			if mode.Current() == Visual && isInSelection(visualStart, *buffer.Cursor, rn, lineIndex) {
				style = style.Reverse(true)
			}

			if r == '\t' {
				spaces := 4 - (x % 4)
				for t := 0; t < spaces; t++ {
					screen.SetContent(x, y, tab, nil, style)
					x++
				}
			} else {
				screen.SetContent(x, y, r, nil, style)
				x++
			}
		}
	}

	status := ""
	if mode.Current() == Normal {
		status = "-- NORMAL --" + " \\ " + buffer.Filename
	} else if mode.Current() == Insert {
		status = "-- INSERT --" + " \\ " + buffer.Filename
	} else if mode.Current() == Command {
		status = ":" + string(buffer.Command)
	}
	if buffer.StatusMsg != "" && mode.Current() != Command {
		status += " | " + buffer.StatusMsg
	}
	// clear bottom line first
	for x := 0; x < screenW; x++ {
		screen.SetContent(x, screenH-1, ' ', nil, tcell.StyleDefault)
	}

	for x, r := range status {
		screen.SetContent(x, screenH-1, r, nil, tcell.StyleDefault)
	}

		
	//define coordinates
	coords := fmt.Sprintf("/ %d:%d ", buffer.Cursor.Y+1, buffer.Cursor.X+1)
	startcoords := screenW - len(coords)
	
	operators := string(buffer.KeyReg)
	startoperators := startcoords - len(coords)
	// statusend := operators + coords

	//render coordinates at bottom-right corner
	for i, r := range operators {
		screen.SetContent(startoperators+i, screenH-1,r,nil, tcell.StyleDefault)
	}
	for i, r := range coords {
		screen.SetContent(startcoords+i, screenH-1, r, nil, tcell.StyleDefault)
	}
	//render cursor
	atcol := visualColumn([]rune(buffer.Lines[buffer.Cursor.Y]), buffer.Cursor.X, 4)
	screen.ShowCursor(atcol, buffer.Cursor.Y-buffer.ScrollY)

	screen.Show()
}
