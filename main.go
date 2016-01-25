package main

import (
	"fmt"
	"os"

	"github.com/alexcesaro/log"
	"github.com/jessevdk/go-flags"

	ui "github.com/gizak/termui"
)

// version gets replaced during build
var version string = "dev"

// logger gets replaced by golog
var logger log.Logger = log.NullLogger

// Options contains the flag options
type Options struct {
	Verbose []bool `short:"v" long:"verbose" description:"Show verbose logging."`
	Version bool   `long:"version"`
}

var logLevels = []log.Level{
	log.Warning,
	log.Info,
	log.Debug,
}

func fail(code int, format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
	os.Exit(code)
}

func main() {
	options := Options{}
	parser := flags.NewParser(&options, flags.Default)
	args, err := parser.Parse()
	if err != nil {
		os.Exit(1)
		return
	}

	if options.Version {
		fmt.Println(version)
		os.Exit(0)
	}

	// Figure out the log level
	numVerbose := len(options.Verbose)
	if numVerbose > len(logLevels) {
		numVerbose = len(logLevels) - 1
	}

	logLevel := logLevels[numVerbose]
	//logger = golog.New(os.Stderr, logLevel)
	_ = logLevel
	logger = log.NullLogger

	if len(args) < 1 {
		fail(1, "Missing arg: PATH")
	}

	if err := ui.Init(); err != nil {
		panic(err)
	}
	defer ui.Close()

	g := ui.NewGauge()
	g.Width = 50
	g.Label = "Womps"
	ui.Render(g)

	setter := func(v float64) {
		g.Percent = int(v * 100)
		ui.Render(g)
	}
	vis := BasicVisualizer(setter)

	ui.Handle("/sys/kbd", func(ui.Event) {
		// Any key to quit
		ui.StopLoop()
	})

	go ui.Loop()

	err = PlayPath(args[0], vis)
	if err != nil {
		fail(2, "Play error: %s", err)
	}

	logger.Info("Done.")
	os.Exit(0)
}
