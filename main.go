package main

import (
	"fmt"
	"os"
)

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		printUsage()
		return
	}

	switch args[0] {
	case "--completions":
		printCompletions()
		return
	case "--refresh":
		if err := refreshCache(); err != nil {
			fatal(err)
		}
		fmt.Println("Topic cache refreshed.")
		return
	case "-fzf":
		if err := requireCmd("fzf"); err != nil {
			fatal(err)
		}
		if err := requireCmd("bat"); err != nil {
			fatal(err)
		}
		topic, err := fzfSelect()
		if err != nil {
			fatal(err)
		}
		if topic == "" {
			return
		}
		if err := fetchAndDisplay(topic); err != nil {
			fatal(err)
		}
		return
	}

	// Direct mode: first arg is the topic, rest ignored.
	if err := requireCmd("bat"); err != nil {
		fatal(err)
	}
	if err := fetchAndDisplay(args[0]); err != nil {
		fatal(err)
	}
}

func printUsage() {
	fmt.Fprintln(os.Stderr, `usage: cht <topic>         fetch cheat sheet
       cht -fzf            fuzzy search topics
       cht --refresh        refresh topic cache
       cht --completions    print zsh completions`)
}

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "cht: %s\n", err)
	os.Exit(1)
}
