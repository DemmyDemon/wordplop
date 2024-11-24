package plopper

import (
	"fmt"
	"math/rand/v2"
	"time"

	"github.com/DemmyDemon/wordplop/pile"
	tsize "github.com/kopoli/go-terminal-size"
)

type Plopper struct {
	active      bool
	withTime    bool
	wordsTarget int
	pile        pile.WordPile
	columns     int
	rows        int
	words       []PlopWord
	colorName   string
}

func remove(slice []PlopWord, i int) []PlopWord {
	return append(slice[:i], slice[i+1:]...)
}

func New(pile pile.WordPile, colorName string) Plopper {
	pl := Plopper{
		active:    true,
		pile:      pile,
		words:     []PlopWord{},
		colorName: colorName,
	}

	pl.maybeResize()

	return pl
}

func (pl *Plopper) Render() string {

	s := ""

	for _, word := range pl.words {
		s += word.Render()
	}
	return s
}

func (pl *Plopper) Draw() {
	fmt.Printf("\033[%d;%dH", 0, 0)
	fmt.Print(pl.Render())
	if pl.withTime {
		pl.drawTime()
	}
}

func (pl *Plopper) drawTime() {
	now := time.Now()
	column := int(pl.rows / 2)
	row := int(pl.columns/2) - 6

	color := GetColorByName(pl.colorName)
	fmt.Printf("\033[38;2;%d;%d;%dm", color[0], color[1], color[2])

	fmt.Printf("\033[%d;%dH┏━━━━━━━━━━┓", column-1, row)
	fmt.Printf("\033[%d;%dH┃ %02d:%02d:%02d ┃", column, row, now.Hour(), now.Minute(), now.Second())
	fmt.Printf("\033[%d;%dH┗━━━━━━━━━━┛", column+1, row)
}

func (pl *Plopper) clearTime() {
	column := int(pl.rows / 2)
	row := int(pl.columns/2) - 6

	fmt.Printf("\033[%d;%dH            ", column-1, row)
	fmt.Printf("\033[%d;%dH            ", column, row)
	fmt.Printf("\033[%d;%dH            ", column+1, row)
}

func (pl *Plopper) ToggleTimeDrawing() {
	pl.withTime = !pl.withTime
	if !pl.withTime {
		pl.clearTime()
	}
}

func (pl *Plopper) Add(word ...PlopWord) {
	pl.words = append(pl.words, word...)
}

func (pl *Plopper) Clear() {
	pl.words = []PlopWord{}
	fmt.Print("\033[2J") // Clear
}

func (pl *Plopper) ToggleActive() {
	pl.active = !pl.active
}

func (pl *Plopper) SetActive(active bool) {
	pl.active = active
}

func (pl *Plopper) IsActive() bool {
	return pl.active
}

func (pl *Plopper) SetColorName(colorName string) {
	pl.colorName = colorName
}

func (pl *Plopper) Resize(size tsize.Size) {

	pl.wordsTarget = (size.Height / 2) * (size.Width / 10)
	pl.columns = size.Width
	pl.rows = size.Height

	pl.Clear()
}

func (pl *Plopper) maybeResize() {
	size, err := tsize.GetSize()
	if err != nil {
		panic(err) // Hypothetical error
	}

	if size.Width != pl.columns || size.Height != pl.rows {
		pl.Resize(size)
	}
}

func (pl *Plopper) Update() {

	pl.maybeResize()

	dead := []int{}
	for i := 0; i < len(pl.words); i++ {
		if pl.words[i].Tick() {
			dead = append(dead, i)
		}
	}

	if len(dead) > 0 {
		for i := len(dead) - 1; i >= 0; i-- {
			pl.words = remove(pl.words, dead[i])
		}
	}

	if len(pl.words) < pl.wordsTarget && pl.active {
		if rand.IntN(1000) <= 900 {
			pl.words = append(pl.words, NewWord(pl.pile, pl.colorName))
		}
	}

}
