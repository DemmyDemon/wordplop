package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"

	"atomicgo.dev/keyboard"
	"atomicgo.dev/keyboard/keys"
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

	loading := []plopper.PlopWord{}

	for _, name := range FileList() {
		message := wordPile.AddFile(name)
		fmt.Println(message)
		loading = append(loading, plopper.PlopWord{
			Column: 1,
			Row:    len(loading) + 1,
			Word:   message,
			Colors: plopper.GetColorByName("white"),
			Life:   250,
			Max:    250,
			Intro:  len(message),
		})
	}

	plop := plopper.New(wordPile, "red")
	plop.Add(loading...)
	fmt.Print("\033[2J")   // Clear
	fmt.Print("\033[?25l") // Hide cursor

	fmt.Printf("\033]0;%s WordPlop\007", wordPile.TopWord)

	ticker := time.NewTicker(INTERVAL)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		keyboard.Listen(func(key keys.Key) (stop bool, err error) {
			switch key.Code {
			case keys.CtrlC, keys.Escape:
				close(quit)
				return true, nil
			case keys.Space:
				plop.ToggleActive()
				plop.Clear()
			case keys.RuneKey:
				switch key.String() {
				case "q", "Q":
					close(quit)
					return true, nil
				case "1":
					plop.SetColorName("red")
				case "2":
					plop.SetColorName("green")
				case "3":
					plop.SetColorName("yellow")
				case "4":
					plop.SetColorName("orange")
				case "5":
					plop.SetColorName("blue")
				case "6":
					plop.SetColorName("white")
				case "7":
					plop.SetColorName("dragonberry")
				}
			}
			return false, nil
		})
	}()

	for {
		select {
		case <-ticker.C:
			plop.Update()
			plop.Draw()
		case <-quit:
			ticker.Stop()
			fmt.Print("\033[?25h") // Enable cursor
			fmt.Print("\033[0m")   // Reset color
			fmt.Print("\033[2J")   // Clear
			fmt.Print("\033[0;0H") // Move cursor top left
			fmt.Println("OK BYE")
			os.Exit(0)
		}
	}
}
