package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

type TFlags struct {
	UseMutex    bool
	ShowVersion bool

	Formatter string

	PathToSrc string
	PathToDst string
}

var (
	ErrEmptyPathToSrc = errors.New("flag '-src' can't be empty")
	ErrEmptyPathToDst = errors.New("flag '-dst' can't be empty")
)

func (this *TFlags) Validate() error {
	if this.PathToSrc == "" {
		return ErrEmptyPathToSrc
	}

	if this.PathToDst == "" {
		return ErrEmptyPathToDst
	}

	return nil
}

var (
	Version string = "unset"

	Flags TFlags
)

func init() {
	flag.BoolVar(&Flags.UseMutex, "safe", false, "use mutex to protect hook methods")
	flag.BoolVar(&Flags.ShowVersion, "version", false, "show the version for hoog")
	flag.StringVar(&Flags.Formatter, "fmt", "gofmt", "go pretty-printer: gofmt, goimports or noop (default gofmt)")
	flag.StringVar(&Flags.PathToSrc, "src", "", "path to interface")
	flag.StringVar(&Flags.PathToDst, "dst", "", "path to output generated file")

	flag.Parse()
}

func main() {
	if Flags.ShowVersion {
		fmt.Printf("hoog version: %s\n", Version)
		os.Exit(0)
	}

	if err := Flags.Validate(); err != nil {
		fmt.Printf("error: %s\n\n", err)
		flag.Usage()
		os.Exit(1)
	}
}
