// Here is the main myprolog program.
package main

import (
	"flag"
	"fmt"

	"github.com/xavier268/myprolog/config"
	"github.com/xavier268/myprolog/repl"
)

// File to load and run at startup
var FlagStartFile = ""
var FlagHelp bool

func init() {

	flag.StringVar(&FlagStartFile, "l", "", "Load file at start-up, if none, will run interactive mode.")
	flag.BoolVar(&config.FlagDebug, "d", false, "debug mode")
	flag.BoolVar(&config.FlagVerbose, "v", false, "verbose mode")
	flag.BoolVar(&FlagHelp, "h", false, "print help and exit")

}

func main() {
	flag.Parse()

	if FlagHelp {
		config.PrintFullWelcome()
		flag.PrintDefaults()
		return
	}

	if FlagStartFile == "" {
		repl.REPL()
	} else {
		config.PrintFullWelcome()
		pc, err := repl.RunFile(FlagStartFile)
		if err != nil {
			fmt.Println("Error : ", err)
		}
		fmt.Printf("\nResults : %s\n", pc.ResultString())
	}
}
