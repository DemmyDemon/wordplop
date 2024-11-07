package plopper

import (
	"math/rand/v2"

	"github.com/DemmyDemon/wordplop/pile"
	tsize "github.com/kopoli/go-terminal-size"
)

type PlopWord struct {
	Column int
	Row    int
	Word   string
	Colors []int
	Life   int
	Max    int
}

func NewWord(pl pile.WordPile) PlopWord {
	word := pl.GetWord()
	size, err := tsize.GetSize()
	if err != nil {
		panic(err)
	}

	column := rand.IntN(size.Width - len(word))
	row := rand.IntN(size.Height)

	color := ColorsGrayscale

	life := len(color) - 1
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
	}
}

func (word PlopWord) GetColor() int {
	ratio := float64(word.Life) / float64(word.Max)
	index := int(ratio * float64(len(word.Colors)-1))
	return word.Colors[index]
}
