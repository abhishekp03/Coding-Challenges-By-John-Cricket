package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
)

var manual = fmt.Sprintf(`%s

NAME
%s - print newline, word, and byte counts for each file

SYNOPSIS
%s [OPTION]... [FILE]...

DESCRIPTION
Print  newline,  word,  and  byte counts for each FILE, and a total line if more than one FILE is specified.  A word is a non-zero-length sequence of characters delimited by
white space.

With no FILE, or when FILE is -, read standard input.

The options below may be used to select which counts are printed, always in the following order: newline, word, character, byte.

-c, --bytes
	   print the byte counts

-m, --chars
	   print the character counts

-l, --lines
	   print the newline counts

-w, --words
	   print the word counts

-h, --help display this help and exit

AUTHOR
Written by Abhishek Pandey.
`, os.Args[0], os.Args[0], os.Args[0])

type config struct {
	filePath    string
	printManual bool
	printBytes  bool
	printChars  bool
	printWords  bool
	printLines  bool
}

func parseArgs(args []string) (config, error) {
	c := config{}

	if len(args) == 0 {
		return c, errors.New("include at least one argument")
	}

	for _, arg := range args {
		switch arg {
		case "-h", "--help":
			c.printManual = true
		case "-c", "--bytes":
			c.printBytes = true
		case "-m", "--chars":
			c.printChars = true
		case "-w", "--words":
			c.printWords = true
		case "-l", "--lines":
			c.printLines = true
		default:
			if c.filePath != "" {
				return c, errors.New(fmt.Sprintf("multiple filenames provided or unrecognized argument: %s", arg))
			}
			c.filePath = arg
		}
	}

	if c.printManual && (c.printBytes || c.printChars || c.printWords || c.printLines) {
		return c, errors.New("usage of \"-h\" or \"--help\" with other flags is not permitted")
	}

	if !c.printManual && !c.printBytes && !c.printChars && !c.printWords && !c.printLines {
		c.printBytes = true
		c.printChars = true
		c.printWords = true
		c.printLines = true
	}

	return c, nil
}

func runCmd(r io.Reader, c config) error {
	if c.printManual {
		fmt.Print(manual)
		return nil
	}

	var data []byte
	var err error

	if c.filePath != "" {
		data, err = os.ReadFile(c.filePath)
		if err != nil {
			return fmt.Errorf("error reading from file %s: %w", c.filePath, err)
		}
	} else {
		data, err = io.ReadAll(r)
		if err != nil {
			return fmt.Errorf("error reading from stdin: %w", err)
		}
	}

	if c.printBytes {
		fmt.Printf("%d\t", len(data))
	}

	if c.printChars {
		chars := len(bytes.Runes(data))
		fmt.Printf("%d\t", chars)
	}

	if c.printWords {
		words := len(bytes.Fields(data))
		fmt.Printf("%d\t", words)
	}

	if c.printLines {
		lines := 0
		for _, b := range data {
			if b == '\n' {
				lines++
			}
		}
		fmt.Printf("%d\t", lines)
	}

	if c.filePath != "" {
		fmt.Printf(c.filePath)
	}

	fmt.Println()

	return nil
}

func main() {
	c, err := parseArgs(os.Args[1:])

	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if err = runCmd(os.Stdin, c); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
}
