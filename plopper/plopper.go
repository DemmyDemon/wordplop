package plopper

import (
	"fmt"
	"math/rand/v2"
	"strings"

	"github.com/DemmyDemon/wordplop/pile"
	tsize "github.com/kopoli/go-terminal-size"
)

type Plopper struct {
	wordsTarget int
	pile        pile.WordPile
	columns     int
	rows        int
	words       []PlopWord
}

func remove(slice []PlopWord, i int) []PlopWord {
	return append(slice[:i], slice[i+1:]...)
}

func New(pile pile.WordPile) Plopper {
	size, err := tsize.GetSize()
	if err != nil {
		panic(err) // Hypothetical error
	}
	return Plopper{
		wordsTarget: (size.Height / 2) * (size.Width / 10),
		pile:        pile,
		words:       []PlopWord{},
	}
}

func (pl *Plopper) Render() string {

	s := ""

	for _, word := range pl.words {
		s += fmt.Sprintf("\033[%d;%dH", word.Row, word.Column)
		if word.Life > 0 {

			s += fmt.Sprintf("\033[38;5;%dm%s", word.GetColor(), word.Word)
		} else {
			s += fmt.Sprintf("\033[0m%s", strings.Repeat(" ", len(word.Word)))
		}
		// s += "\n"
	}
	return s
}

func (pl *Plopper) Draw() {
	fmt.Printf("\033[%d;%dH", 0, 0)
	fmt.Print(pl.Render())
}

func (pl *Plopper) Resize(size tsize.Size) {

	pl.wordsTarget = (size.Height / 2) * (size.Width / 10)
	pl.words = []PlopWord{}

	pl.columns = size.Width
	pl.rows = size.Height

	fmt.Print("\033[2J") // Clear
}

func (pl *Plopper) Update() {

	size, err := tsize.GetSize()
	if err != nil {
		panic(err) // Hypothetical error
	}

	if size.Width != pl.columns || size.Height != pl.rows {
		pl.Resize(size)
	}

	dead := []int{}
	for i, word := range pl.words {
		word.Life--
		if word.Life < 0 {
			dead = append(dead, i)
		}
		pl.words[i] = word
	}
	// s += fmt.Sprintf("Dead:  %#v\n", dead)

	if len(dead) > 0 {
		for i := len(dead) - 1; i >= 0; i-- {
			pl.words = remove(pl.words, dead[i])
		}
	}

	if len(pl.words) < pl.wordsTarget {
		if rand.IntN(1000) < 900 {
			pl.words = append(pl.words, NewWord(pl.pile))
		}
		// s += "Added a word!\n"
	}

}
