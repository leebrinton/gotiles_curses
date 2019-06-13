package main

import (
	"fmt"
	//"log"
	"os"

	"github.com/leebrinton/tileslib"
	//"code.google.com/p/goncurses"
	//"github.com/rthornton128/goncurses"
	"github.com/gbin/goncurses"
)

const (
	numLines     = 21
	numCols      = 21
	numGameCells = 16
)

///////////////////////////////////////////////////////////////////////////////

// DrawScreenChars - A structure that holds the set of characters used to draw
// a line in a window
type DrawScreenChars struct {
	LeftSide  goncurses.Char
	Middle    goncurses.Char
	Filler    goncurses.Char
	RightSide goncurses.Char
}

func (drawChars *DrawScreenChars) loadBorderChars(row int, lineChars *LineChars) {
	switch row {
	case 0: /* top of banner area */
		drawChars.LeftSide = lineChars.UlCorner
		drawChars.Middle = lineChars.HLine
		//drawChars.Middle = lineChars.TopTee
		//drawChars.Middle = ' '
		drawChars.Filler = lineChars.HLine
		drawChars.RightSide = lineChars.UrCorner

	case 4: /* bottom of banner area and top of cell area */
		drawChars.LeftSide = lineChars.LTee
		drawChars.Middle = lineChars.TopTee
		drawChars.Filler = lineChars.HLine
		drawChars.RightSide = lineChars.RTee

	case 20: /* bottom of cell area */
		drawChars.LeftSide = lineChars.LlCorner
		drawChars.Middle = lineChars.BotTee
		drawChars.Filler = lineChars.HLine
		drawChars.RightSide = lineChars.LrCorner

	default: /* interior cell border */
		drawChars.LeftSide = lineChars.LTee
		drawChars.Middle = lineChars.Plus
		drawChars.Filler = lineChars.HLine
		drawChars.RightSide = lineChars.RTee
	}
}

func (drawChars *DrawScreenChars) loadCellChars(row int, lineChars *LineChars) {
	drawChars.LeftSide = lineChars.VLine
	drawChars.Filler = lineChars.Space
	drawChars.RightSide = lineChars.VLine

	//  If the line is in the banner area ...
	if row < 4 {
		drawChars.Middle = lineChars.Space
	} else { /* line is in the cell area */
		drawChars.Middle = lineChars.VLine
	}
}

// Load - Load the appropriate characters for drawing a border or a cell
func (drawChars *DrawScreenChars) Load(row int, lineChars *LineChars) {
	// the banner and cell border linea are 4 lines apart
	lineType := row % 4

	if lineType == 0 {
		drawChars.loadBorderChars(row, lineChars)
	} else {
		drawChars.loadCellChars(row, lineChars)
	}
}

///////////////////////////////////////////////////////////////////////////////

// InputCode - Define a type for an input code enumeration
type InputCode int

const (
	icNone InputCode = iota
	icUp
	icDown
	icLeft
	icRight
	icQuit
	icHelp
)

///////////////////////////////////////////////////////////////////////////////

// Game - A representation of a game
type Game struct {
	config *CursesConfig
	model  *tileslib.Model
	stdscr *goncurses.Window
	win    *goncurses.Window
}

// NewGame - Create an instance of a game
func NewGame(config *CursesConfig, model *tileslib.Model, stdscr *goncurses.Window) *Game {
	game := new(Game)

	game.config = config
	game.model = model
	game.stdscr = stdscr

	return game
}

func (game *Game) drawLine(row int, drawChars *DrawScreenChars) {
	game.win.HLine(row, 0, drawChars.LeftSide, 1)
	game.win.HLine(row, 1, drawChars.Filler, 4)
	game.win.HLine(row, 5, drawChars.Middle, 1)
	game.win.HLine(row, 6, drawChars.Filler, 4)
	game.win.HLine(row, 10, drawChars.Middle, 1)
	game.win.HLine(row, 11, drawChars.Filler, 4)
	game.win.HLine(row, 15, drawChars.Middle, 1)
	game.win.HLine(row, 16, drawChars.Filler, 4)
	game.win.HLine(row, 20, drawChars.RightSide, 1)
}

func (game *Game) showCell(cellIndex int) {
	value := game.model.CellValueAt(cellIndex)
	display := "  "

	// build the value string
	if value != numGameCells {
		display = fmt.Sprintf("%2d", value)
	}

	// Calculate the the cell row an column
	row := tileslib.RowFromIndex(byte(cellIndex))
	col := tileslib.ColFromIndex(byte(cellIndex))

	// Write a string to the cell
	game.win.MovePrint(rowAddress(row), colAddress(col), display)
}

func (game *Game) drawScreen() {
	// Set the window colors
	if goncurses.HasColors() {
		//game.win.Color(1)
		game.win.ColorOn(1)
	}

	// Make the cursor invisible
	goncurses.Cursor(0)

	// Draw the banner and cell borders
	var drawChars DrawScreenChars

	for i := 0; i < numLines; i++ {
		drawChars.Load(i, &game.config.Chars)
		game.drawLine(i, &drawChars)
	}

	// Write the banner title
	game.win.MovePrint(2, 6, "T I L E S")

	// Write the cell values to the window
	for i := 0; i < numGameCells; i++ {
		game.showCell(i)
	}

	// Write the window to the screen
	game.win.Refresh()
}

// Init - Initialize a game; formulate the game data, create a window and
// draw the game in the window
func (game *Game) Init() {
	termlines := 50
	termcols := 80

	beginx := ((termlines - numLines) / 2)
	beginy := ((termcols - numCols) / 2)

	game.win = game.stdscr.Sub(numLines, numCols, beginx, beginy)

	game.model.StartNewGame(game.config.Iterations())
	game.drawScreen()
}

// GetInputCode - Get an input and convert it to an InputCode if possible
func (game *Game) GetInputCode() InputCode {
	code := icNone

	for code == icNone {
		goncurses.Echo(false)

		key := goncurses.StdScr().GetChar()
		//log.Printf("getInputCode key = %d\n", int(key))

		switch key {
		case goncurses.KEY_UP:
			code = icUp

		case goncurses.KEY_DOWN:
			code = icDown

		case goncurses.KEY_LEFT:
			code = icLeft

		case goncurses.KEY_RIGHT:
			code = icRight

		case goncurses.KEY_END:
			code = icQuit

		case goncurses.KEY_F1:
			code = icHelp
		}

		if 'k' == int(key) {
			code = icUp
		} else if 'j' == int(key) {
			code = icDown
		} else if 'h' == int(key) {
			code = icLeft
		} else if 'l' == int(key) {
			code = icRight
		} else if 'q' == int(key) {
			code = icQuit
		}

		if code == icNone {
			//log.Printf("Unknown key: %c\n", key)
			goncurses.Flash()
		}
	}
	return code
}

func (game *Game) swapCells() {
	/* Blank out the dest and write new value at the source */
	game.showCell(int(game.model.LastDest()))
	game.showCell(int(game.model.LastSource()))

	// Write the window changes to the screen
	game.win.Move(-1, -1)
	game.win.Refresh()
	//goncurses.Cursor(0)
}

func (game *Game) command(direction tileslib.Direction) {
	game.model.MoveCell(direction, game.config.CommandMode())

	transresult := game.model.LastTransResult()
	switch transresult {
	case tileslib.Pending:
		fatalError("Unexpected result 'result_pending' in game_command")

	case tileslib.Exception:
		goncurses.Flash()

	case tileslib.Error:
		fatalError("Error while processing game transaction")

	case tileslib.Ok:
		game.swapCells()

	default:
		fatalError("Unknown result in game command")
	}
}

func (game *Game) handleInput(inputCode InputCode) (bool, bool) {
	quitRequested := false
	gameSolved := false

	switch inputCode {
	case icUp:
		game.command(tileslib.Up)

	case icDown:
		game.command(tileslib.Down)

	case icLeft:
		game.command(tileslib.Left)

	case icRight:
		game.command(tileslib.Right)

	case icHelp:
		//display_help( config, line_chars );
		game.drawScreen()

	case icQuit:
		quitRequested = true
	}

	// See if the game is solved
	if game.model.Solved() {
		gameSolved = true
	}
	return quitRequested, gameSolved
}

// Play - Get and respond to input until game is solved or quit is requested
func (game *Game) Play() bool {
	game.Init()
	play := true
	solved := false
	quit := false

	for play {
		code := game.GetInputCode()
		quit, solved = game.handleInput(code)
		play = !(quit || solved)
		//log.Printf("quit: %t solved: %t play %t\n", quit, solved, play)
	}
	return !quit
}

///////////////////////////////////////////////////////////////////////////////
func rowAddress(index byte) int {
	return int((index * 4) + 6)
}

func colAddress(index byte) int {
	return int((index * 5) + 2)
}

func fatalError(msg string) {
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(1)
}

///////////////////////////////////////////////////////////////////////////////
