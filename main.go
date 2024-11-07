package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"time"

	"github.com/DemmyDemon/wordplop/pile"
	"github.com/DemmyDemon/wordplop/plopper"
)

const INTERVAL = time.Duration(0.05 * float64(time.Second))

func MaybePanic(err error) {
	if err != nil {
		if !errors.Is(err, io.EOF) {
			panic(err)
		}
	}
}

func AllFiles() []string {
	files := []string{}

	entries, err := os.ReadDir("./")
	if err != nil {
		panic(err)
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		files = append(files, entry.Name())
	}

	return files
}

func FileList() []string {
	if len(os.Args) == 1 {
		return AllFiles()
	}
	return os.Args[1:]
}

func main() {

	wordPile := pile.New()

	for _, name := range FileList() {
		wordPile.AddFile(name)
	}

	plop := plopper.New(wordPile)
	fmt.Print("\033[2J")   // Clear
	fmt.Print("\033[?25l") // Hide cursor

	fmt.Printf("\033]0;%s WordPile\007", wordPile.TopWord)

	ticker := time.NewTicker(INTERVAL)

	inputChan := make(chan rune, 1)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	go func() {
		for {
			select {
			case <-ticker.C:
				plop.Update()
				plop.Draw()
			case <-inputChan:
				close(quit)
			case <-quit:
				ticker.Stop()
				fmt.Print("\033[?25h") // Enable cursor
				fmt.Print("\033[0m")   // Reset color
				fmt.Print("\033[2J")   // Clear
				fmt.Print("\033[0;0H") // Set 0,0
				os.Exit(0)
			}
		}
	}()
	for {
		reader := bufio.NewReader(os.Stdin)
		keyStr, _, err := reader.ReadRune()
		MaybePanic(err)
		inputChan <- keyStr
	}

}
