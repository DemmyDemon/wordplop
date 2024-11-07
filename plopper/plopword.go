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
	Color  []int
	Life   int
}

func NewWord(pl pile.WordPile, factor int) PlopWord {
	word := pl.GetWord()
	size, err := tsize.GetSize()
	if err != nil {
		panic(err)
	}

	column := rand.IntN(size.Width - len(word))
	row := rand.IntN(size.Height)

	color := ColorsGrayscale

	return PlopWord{
		Column: column,
		Row:    row,
		Color:  color,
		Life:   (len(color) - 1) * factor,
		Word:   word,
	}
}
