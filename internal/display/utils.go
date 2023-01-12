package display

import (
	"unicode"
)

func SplitHeader(longString string, maxLen int) []string {
	splits := []string{}
	splitsRune := [][]rune{}
	runeString := []rune(longString)

	var l, r int
	for l, r = 0, maxLen; r < len(runeString); l, r = r, r+maxLen {
		for !unicode.IsSpace(rune(runeString[r])) {
			r--
		}
		splitsRune = append(splitsRune, runeString[l:r])
	}
	splitsRune = append(splitsRune, runeString[l:])
	for _, v := range splitsRune {
		splits = append(splits, string(v))
	}
	return splits
}
