package main

import (
	//"code.google.com/p/goncurses"
	//"github.com/rthornton128/goncurses"
	"github.com/gbin/goncurses"
)

func blankLines(dialog *goncurses.Window, beginRow int, beginCol int, lines int, width int) {
	blanks := ""

	for i := 0; i < width; i++ {
		blanks += " "
	}

	for i := beginRow; i < lines; i++ {
		dialog.MovePrint(i, beginCol, blanks)
	}
}

func drawBorder(dialog *goncurses.Window, linechars *LineChars) {
	dialog.Border(
		linechars.VLine, linechars.VLine,
		linechars.HLine, linechars.HLine,
		linechars.UlCorner, linechars.UrCorner,
		linechars.LlCorner, linechars.LrCorner)
}
