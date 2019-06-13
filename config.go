package main

import (
	"github.com/leebrinton/tileslib"
	//"code.google.com/p/goncurses"
	//"github.com/rthornton128/goncurses"
	"github.com/gbin/goncurses"
)

// LineCharType - Create a type for a line character enumeration
type LineCharType int

const (
	// ACS - Graphical line characters
	ACS LineCharType = iota
	// ASCII - Normal symbol characters
	ASCII
)

// LineChars - a structure that holds an example of each type of character
// used to draw the game
type LineChars struct {
	UlCorner goncurses.Char
	UrCorner goncurses.Char
	LlCorner goncurses.Char
	LrCorner goncurses.Char
	HLine    goncurses.Char
	VLine    goncurses.Char
	LTee     goncurses.Char
	RTee     goncurses.Char
	TopTee   goncurses.Char
	BotTee   goncurses.Char
	Plus     goncurses.Char
	Space    goncurses.Char
	NullChar goncurses.Char
}

// Load - Set example character depending on lct
func (lc *LineChars) Load(lct LineCharType) {
	if lct == ACS {
		//lc.UlCorner = goncurses.ACS_ULCORNER
		lc.UlCorner = goncurses.ACS_LLCORNER
		lc.UrCorner = goncurses.ACS_URCORNER
		//lc.LlCorner = goncurses.ACS_LLCORNER
		lc.LlCorner = goncurses.ACS_ULCORNER
		lc.LrCorner = goncurses.ACS_LRCORNER
		lc.HLine = goncurses.ACS_HLINE
		lc.VLine = goncurses.ACS_VLINE
		lc.LTee = goncurses.ACS_LTEE
		lc.RTee = goncurses.ACS_RTEE
		lc.TopTee = goncurses.ACS_TTEE
		lc.BotTee = goncurses.ACS_BTEE
		lc.Plus = goncurses.ACS_PLUS
	} else {
		lc.UlCorner = '/'
		lc.UrCorner = '\\'
		lc.LlCorner = '\\'
		lc.LrCorner = '/'
		lc.HLine = '-'
		lc.VLine = '|'
		lc.LTee = '|'
		lc.RTee = '-'
		lc.TopTee = '-'
		lc.BotTee = '-'
		lc.Plus = '+'
	}

	lc.Space = ' '
	lc.NullChar = 0
}

// CursesConfig - The application configuration elements
type CursesConfig struct {
	base    tileslib.Config
	Chars   LineChars
	FgColor int16
	BgColor int16
}

// NewCursesConfig - Create a new instance of an
// initialized application configuration
func NewCursesConfig() *CursesConfig {
	c := new(CursesConfig)

	c.base.ScrambleIterations = tileslib.DEFAULT_SCRAMBLE_ITERATIONS
	c.base.CommandMode = tileslib.EmptyCellCentric
	c.FgColor = goncurses.C_WHITE
	c.BgColor = goncurses.C_BLACK

	return c
}

// SetIterations - Configure the number of scramble iterations
func (c *CursesConfig) SetIterations(iterations int) {
	c.base.ScrambleIterations = iterations
}

// Iterations - Get the configured number of scramble iterations
func (c *CursesConfig) Iterations() int {
	return c.base.ScrambleIterations
}

// SetCommandMode - Configure the command mode
func (c *CursesConfig) SetCommandMode(mode tileslib.CommandModeType) {
	c.base.CommandMode = mode
}

// CommandMode - Get the configured command mode
func (c *CursesConfig) CommandMode() tileslib.CommandModeType {
	return c.base.CommandMode
}
