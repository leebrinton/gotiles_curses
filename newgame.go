package main

import (
	"fmt"
	//"log"
	"os"
	"time"

	"github.com/leebrinton/tileslib"
	//`"code.google.com/p/goncurses"
	//"github.com/rthornton128/goncurses"
	"github.com/gbin/goncurses"
)

const (
	// nCols - Number of columns in the new game dialog
	nCols = 45
	// anotherGameString -  String asking the question
	anotherGameString = "Want to play another game? (Y/N)"
)

func getDurationString(startTime time.Time) string {
	duration := time.Since(startTime)
	durationStr := ""

	if duration > time.Hour {
		hours := (duration/time.Hour)
		durationStr := fmt.Sprintf("%d hour", hours)
		if hours > 1 {
			durationStr += "s"
		}
		durationStr += " "
		duration -= (hours * time.Hour)
	}

	if duration > time.Minute {
		minutes := (duration/time.Minute)
		tmpStr := fmt.Sprintf("%d minute", minutes)
		if minutes > 1 {
			tmpStr += "s"
		}
		tmpStr += " "
		durationStr += tmpStr
		duration -= (minutes * time.Minute)
	}

	if duration > time.Second {
		seconds := (duration/time.Second)
		tmpStr := fmt.Sprintf("%d second", seconds)
		if seconds > 1 {
			tmpStr += "s"
		}
		tmpStr += " "
		durationStr += tmpStr
	}
	//return duration.String()
	return durationStr
}

func drawCongratsMessage(dialog *goncurses.Window,
                         row *int, col *int, startTime time.Time) {
	durationStr := getDurationString(startTime)

	dialog.MovePrint(*row, *col, "Congratulations, you solved the puzzle in")

	*row++
	*col = int((nCols - len(durationStr)) / 2)
	dialog.MovePrint(*row, *col, durationStr)
}

func drawQuestionMessage(dialog *goncurses.Window, row *int, col *int) {
	*row += 2
	*col = int((nCols - len(anotherGameString)) / 2)
	dialog.MovePrint(*row, *col, anotherGameString)
}

func drawAskDialog(config *CursesConfig,
                   model *tileslib.Model) *goncurses.Window {
	nlines := 5
	solved := model.Solved()
	termlines := 50
	termcols := 80

	// Make room for the congrats message
	if solved {
		nlines += 3
	}

	// Center the new window on the screen
	beginx := ((termlines - nlines) / 2)
	beginy := ((termcols - nCols) / 2)

	// Create the dialog window
	//log.Printf("nlines %d NCOLS %d beginx %d beginy %d\n", nlines, nCols, beginx, beginy)
	dialog, err := goncurses.NewWindow(nlines, nCols, beginx, beginy)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}

	// Set the colors for the dialog
	if goncurses.HasColors() {
		dialog.ColorOn(1)
	}

	// Write spaces through out the dialog
	blankLines(dialog, 1, 1, (nlines - 1), nCols)

	// Write the congrats message
	currow := 2
	curcol := 2

	if solved {
		drawCongratsMessage(dialog, &currow, &curcol, model.StartTime)
	}

	// Write question message
	drawQuestionMessage(dialog, &currow, &curcol)

	// Draw a border around the dialog
	drawBorder(dialog, &config.Chars)

	// Write the dialog to the screen
	dialog.Refresh()

	return dialog
}

func getAnswer(dialog *goncurses.Window) bool {
	var key int
	answer := false
	done := false

	for !done {
		goncurses.Echo(false)
		key = int(dialog.GetChar())

		if 'y' == int(key) || 'Y' == int(key) {
			answer = true
			done = true
		} else if 'n' == int(key) || 'N' == int(key) {
			answer = false
			done = true
		} else if 'q' == int(key) || 'Q' == int(key) {
			answer = false
			done = true
		} else {
			goncurses.Flash()
		}
	}
	return answer
}

// AskPlayNewGame - Display a dialog asking to play a new game and
// return the reply
func AskPlayNewGame(config *CursesConfig, model *tileslib.Model) bool {
	dialog := drawAskDialog(config, model)
	answer := getAnswer(dialog)
	dialog.Erase()
	dialog.Refresh()

	err := dialog.Delete()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	return answer
}
