package main

import (
	"regexp"
)

type Word struct {
	Text  string
	Start int
	End   int
}

func splitIntoWords(s string) []Word {
	re := regexp.MustCompile(`\b\w+\b|[\p{P}-[\.]]+|\.\.+`)
	matchIndexes := re.FindAllStringIndex(s, -1)
	matchStrings := re.FindAllString(s, -1)

	words := make([]Word, len(matchStrings))
	for i := range matchStrings {
		words[i] = Word{
			Text:  matchStrings[i],
			Start: matchIndexes[i][0],
			End:   matchIndexes[i][1],
		}
	}
	return words
}

func findNextWord(words []Word, pos int) *Word {
	for _, word := range words {
		if word.Start > pos {
			return &word
		}
	}
	return nil
}

func findPrevWord(words []Word, pos int) *Word {
	for i := len(words) - 1; i >= 0; i-- {
		if words[i].End <= pos {
			return &words[i]
		}
	}
	return nil
}

func (r *RaptorCfg) JumpToNextLine() {
	if r.CY+r.CurrentRowOffset < r.CurrentFileNumRows-1 {
		r.CY += 1
		r.CX = 0
		currentRowStr := r.Rows[r.CY+r.CurrentRowOffset].Chars
		wlist := splitIntoWords(currentRowStr)
		if len(wlist) > 0 {
			firstWord := wlist[0]
			r.CX = firstWord.Start
		}
	}
}

func (r *RaptorCfg) JumpToNextWordBeginning() {
	currentRowStr := r.Rows[r.CY+r.CurrentRowOffset].Chars
	wlist := splitIntoWords(currentRowStr)
	nextWord := findNextWord(wlist, r.CX)

	if nextWord != nil {
		r.CX = nextWord.Start
	} else {
		r.JumpToNextLine()
	}
}

func (r *RaptorCfg) JumpToPrevWordBeginning() {
	currentRowStr := r.Rows[r.CY+r.CurrentRowOffset].Chars
	wlist := splitIntoWords(currentRowStr)
	prevWord := findPrevWord(wlist, r.CX)

	if prevWord != nil {
		r.CX = prevWord.Start
	} else {
		if r.CY+r.CurrentRowOffset > 0 {
			r.CY -= 1
			currentRowStr = r.Rows[r.CY+r.CurrentRowOffset].Chars
			wlist = splitIntoWords(currentRowStr)
			if len(wlist) > 0 {
				lastWord := wlist[len(wlist)-1]
				r.CX = lastWord.Start
			} else {
				r.CX = 0
			}
		}
	}
}
