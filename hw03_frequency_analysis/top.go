package hw03frequencyanalysis

import (
	"sort"
	"strings"
	"unicode"
)

func trimPunctuation(word string) string {
	runes := []rune(word)
	isAllPunctuation := func(word string) bool {
		for _, r := range word {
			if !unicode.IsPunct(r) {
				return false
			}
		}
		return true
	}
	if isAllPunctuation(word) && len(runes) > 1 {
		return word
	}
	startIndex, endIndex := 0, len(runes)
	for startIndex < len(runes) && unicode.IsPunct(runes[startIndex]) {
		startIndex++
	}
	for endIndex > startIndex && unicode.IsPunct(runes[endIndex-1]) {
		endIndex--
	}
	return string(runes[startIndex:endIndex])
}

func Top10(str string) []string {
	words := make(map[string]int)
	splitedStr := strings.Fields(str)
	for _, word := range splitedStr {
		word = trimPunctuation(word)
		if word == "" {
			continue
		}
		lowedWord := strings.ToLower(word)
		words[lowedWord]++
	}

	type wordCount struct {
		word  string
		count int
	}
	sortedWords := make([]wordCount, 0, len(words))
	for word, count := range words {
		sortedWords = append(sortedWords, wordCount{word, count})
	}

	sort.Slice(sortedWords, func(i, j int) bool {
		if sortedWords[i].count == sortedWords[j].count {
			return sortedWords[i].word < sortedWords[j].word
		}
		return sortedWords[i].count > sortedWords[j].count
	})

	topSize := 10
	if len(sortedWords) < topSize {
		topSize = len(sortedWords)
	}

	sortedSlice := make([]string, 0, topSize)
	for _, word := range sortedWords[:topSize] {
		sortedSlice = append(sortedSlice, word.word)
	}

	return sortedSlice
}
