package pile

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
	"unicode"
)

func MaybeCrap(what error, where string) {
	if what != nil {
		fmt.Fprintf(os.Stderr, "%s: %s", where, what)
		os.Exit(1)
	}
}

type WordPile struct {
	TopWord  string
	topCount int
	counts   map[string]int
	words    []string
}

func clean(line string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsNumber(r) || unicode.IsDigit(r) || unicode.IsSpace(r) {
			return r
		}
		if r == '\'' || r == '-' {
			return r
		}
		return -1
	}, line)
}

func toTitle(word string) string {
	return strings.ToUpper(word[:1]) + strings.ToLower(word[1:])
}

func New() WordPile {
	return WordPile{
		counts: map[string]int{},
		words:  []string{},
	}
}

func (pile WordPile) Count(word string) int {
	if count, ok := pile.counts[word]; ok {
		return count
	}
	return 0
}

func (pile WordPile) GetWord() string {
	return pile.words[rand.Intn(len(pile.words))]
}

func (pile *WordPile) AddFile(fileName string) {
	begin := time.Now()
	fmt.Printf("Reading %s\n", fileName)
	file, err := os.Open(fileName)
	MaybeCrap(err, "Opening file")
	defer func() {
		err := file.Close()
		MaybeCrap(err, "Closing munched file")
	}()

	wordsFound := []string{}

	scanner := bufio.NewScanner(file)
	numLines := 0
	for scanner.Scan() {

		line := scanner.Text()
		line = clean(line)
		if len(line) <= 3 {
			continue
		}

		numLines++

		words := strings.Split(line, " ")
		for _, word := range words {
			if word == "" {
				continue
			}
			word = toTitle(word)
			pile.counts[word]++
			if pile.counts[word] == 1 {
				wordsFound = append(wordsFound, word)
			}
			if pile.counts[word] > pile.topCount {
				pile.topCount = pile.counts[word]
				pile.TopWord = word
			}
		}
	}
	duration := time.Since(begin)
	pile.words = append(pile.words, wordsFound...)
	fmt.Printf("Read %d lines, for %d new words, in %s (total now %d words)\n", numLines, len(wordsFound), duration, len(pile.words))
}
