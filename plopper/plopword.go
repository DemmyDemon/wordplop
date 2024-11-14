package plopper

import (
	"fmt"
	"math/rand/v2"
	"strings"

	"github.com/DemmyDemon/wordplop/pile"
	tsize "github.com/kopoli/go-terminal-size"
)

type PlopWord struct {
	Column int
	Row    int
	Word   string
	Colors [3]int
	Life   int
	Max    int
	Intro  int
}

func NewWord(pl pile.WordPile, colorName string) PlopWord {
	word := pl.GetWord()
	size, err := tsize.GetSize()
	if err != nil {
		panic(err)
	}

	column := rand.IntN(size.Width - len(word))
	row := rand.IntN(size.Height)

	color := GetColorByName(colorName)

	life := 10
	count := pl.Count(word)
	if count > life {
		life = count
	}

	return PlopWord{
		Column: column,
		Row:    row,
		Colors: color,
		Life:   life,
		Max:    life,
		Word:   word,
		Intro:  0,
	}
}

func (word PlopWord) GetColor() int {
	ratio := float64(word.Life) / float64(word.Max)
	index := int(ratio * float64(len(word.Colors)-1))
	return word.Colors[index]
}

func (word PlopWord) InRGB(what string) string {
	r := word.Colors[0]
	g := word.Colors[1]
	b := word.Colors[2]

	ratio := float64(word.Life) / float64(word.Max)

	r = int(ratio * float64(r))
	g = int(ratio * float64(g))
	b = int(ratio * float64(b))

	return fmt.Sprintf("\033[38;2;%d;%d;%dm%s", r, g, b, what)
}

func (word PlopWord) Render() string {
	s := fmt.Sprintf("\033[%d;%dH", word.Row, word.Column)
	if word.Intro < len(word.Word) {
		s += word.InRGB(word.Word[:word.Intro])
	} else if word.Life > 0 {
		s += word.InRGB(word.Word)
	} else {
		s += fmt.Sprintf("\033[0m%s", strings.Repeat(" ", len(word.Word)))
	}
	return s
}

func (word *PlopWord) Tick() bool {
	if word.Intro < len(word.Word) {
		if rand.IntN(1000) >= 750 {
			word.Intro++
		}
		return false
	}

	word.Life--
	return word.Life < 0
}
