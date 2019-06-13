package main

import (
	"flag"
	"fmt"
	//"log"
	"os"
	"strings"

	"github.com/leebrinton/tileslib"
	//`"code.google.com/p/goncurses"
	//"github.com/rthornton128/goncurses"
	"github.com/gbin/goncurses"
)

func showCmdLineHelp() {
	fmt.Fprintf(os.Stderr, "Usage: %s [OPTION] ...\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(0)
}

func showVersion() {
	fmt.Printf("%s %s A puzzle game\n", os.Args[0], "1.0")
	os.Exit(0)
}

func convertColor(colorStr string, defaultColor int16) int16 {
	result := defaultColor

	if "black" == strings.ToLower(colorStr) {
		result = goncurses.C_BLACK
	} else if "red" == strings.ToLower(colorStr) {
		result = goncurses.C_RED
	} else if "green" == strings.ToLower(colorStr) {
		result = goncurses.C_GREEN
	} else if "yellow" == strings.ToLower(colorStr) {
		result = goncurses.C_YELLOW
	} else if "blue" == strings.ToLower(colorStr) {
		result = goncurses.C_BLUE
	} else if "magenta" == strings.ToLower(colorStr) {
		result = goncurses.C_MAGENTA
	} else if "cyan" == strings.ToLower(colorStr) {
		result = goncurses.C_CYAN
	} else if "white" == strings.ToLower(colorStr) {
		result = goncurses.C_WHITE
	} else {
		fmt.Fprintf(os.Stderr, "Unknown color: %s", colorStr)
	}
	return result
}

func processCmdLine(config *CursesConfig) {
	var asciiChars = flag.Bool("ascii", false, "use ascii chars for line drawing")
	var bg = flag.String("bg", "black", "the background color")
	var valueCellCentric = flag.Bool("value-cell-centric", false, "directional key commands move a value cell")
	var emptyCellCentric = flag.Bool("empty-cell-centric", false, "directional key commands move the empty cell")
	var fg = flag.String("fg", "white", "the foreground color")
	var help = flag.Bool("help", false, "show a help message and exit")
	var iterations = flag.Int("iterations", tileslib.DEFAULT_SCRAMBLE_ITERATIONS, "the number of scramble iterations")
	var acsChars = flag.Bool("acs", false, "use acs chars for line drawing")
	var version = flag.Bool("version", false, "show the version and exit")

	flag.Parse()

	if *asciiChars && !*acsChars {
		config.Chars.Load(ASCII)
	} else {
		config.Chars.Load(ACS)
	}

	if *valueCellCentric && !*emptyCellCentric {
		config.SetCommandMode(tileslib.ValueCellCentric)
	}

	if *help {
		showCmdLineHelp()
	}

	if *version {
		showVersion()
	}

	config.FgColor = convertColor(*fg, goncurses.C_WHITE)
	config.BgColor = convertColor(*bg, goncurses.C_BLACK)
	config.SetIterations(*iterations)
}

func initColor(config *CursesConfig) {
	if goncurses.HasColors() {
		err := goncurses.StartColor()

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}

		err = goncurses.InitPair(1, config.FgColor, config.BgColor)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Unable to initialize a color pair: "+err.Error())
		}
	}
}

func initialize(config *CursesConfig, model *tileslib.Model) *goncurses.Window {
	stdscr, err := goncurses.Init()

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	initColor(config)
	stdscr.Keypad(true)

	return stdscr
}

func main() {
	config := NewCursesConfig()
	model := tileslib.NewModel()

	// logfile, err := os.OpenFile("tiles.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	// if err != nil {
	// 	log.Fatalln("Failed to open log file [tiles.log]: ", err)
	// } else {
	// 	defer logfile.Close()
	// }

	// log.SetOutput(logfile)
	// log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	processCmdLine(config)
	stdscr := initialize(config, model)

	defer goncurses.End()

	game := NewGame(config, model, stdscr)
	answer := tileslib.Accept

	for answer == tileslib.Accept {
		answer = tileslib.Quit

		gameResult := game.Play()
		//log.Printf("gameResult: %t\n", gameResult)

		if gameResult {
			x := AskPlayNewGame(config, model)
			//log.Println(fmt.Sprintf("AskPlayNewGame: %t", x))
			if x {
				answer = tileslib.Accept
			}
		}
	}
}
