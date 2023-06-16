package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/jhr22/daily-snip/internal/snips"
)

const (
	SnipVimPackLoc = "%s/pack/minpac/start/vim-go/gosnippets/UltiSnips/go.snippets"
)

func main() {
	snipFile := fmt.Sprintf(SnipVimPackLoc, os.Getenv("VIMCONFIG"))

	readFile, err := os.Open(snipFile)
	defer readFile.Close()

	if err != nil {
		log.Fatal(err)
	}

	fileScanner := bufio.NewScanner(readFile)

	fileScanner.Split(bufio.ScanLines)

	snippets, err := snips.ParseLines(fileScanner)

	if err != nil {
		log.Fatal(err)
	}

	// snip := snippets[time.Now().YearDay()%len(snippets)]
	// fmt.Println(snip)

	runTUI(snippets)
}
