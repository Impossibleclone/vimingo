package main

import "github.com/gdamore/tcell/v2"

func RenderScreen(screen tcell.Screen, buffer *Buffer, visualStart Cursor, mode *EditorMode) {
	screen.Clear()

	screenW, screenH := screen.Size()

	for y := 0; y < screenH-1; y++ {
		lineIndex := y + buffer.ScrollY
		if lineIndex >= len(buffer.Lines) {
			break
		}
		line := buffer.Lines[lineIndex]
		for x, r := range line {
			style := tcell.StyleDefault
			if mode.Current() == Visual && isInSelection(visualStart, *buffer.Cursor, x, lineIndex) {
				style = style.Reverse(true)
			}
			screen.SetContent(x, y, r, nil, style)
		}
	}

	status := ""
	if mode.Current() == Normal {
		status = "-- NORMAL --"
	} else if mode.Current() == Insert {
		status = "-- INSERT --"
	} else if mode.Current() == Command {
		status = ":" + string(buffer.Command)
	}

	// clear bottom line first
	for x := 0; x < screenW; x++ {
		screen.SetContent(x, screenH-1, ' ', nil, tcell.StyleDefault)
	}

	for x, r := range status {
		screen.SetContent(x, screenH-1, r, nil, tcell.StyleDefault)
	}

	screen.ShowCursor(buffer.Cursor.X, buffer.Cursor.Y-buffer.ScrollY)
	screen.Show()
}
