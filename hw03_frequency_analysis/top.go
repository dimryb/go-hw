package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

func Top10(str string) []string {
	words := make(map[string]int)
	splitedStr := strings.Fields(str)
	for _, word := range splitedStr {
		words[word]++
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
