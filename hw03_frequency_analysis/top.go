package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

type wordFrequency struct {
	word  string
	count int
}

var sep = regexp.MustCompile(`[.,!:?;"'\s]+-*\s*`) // Spaces and punctuation or dash

func Top10(inputStr string) []string {
	res := []string{}

	if len(inputStr) > 0 {
		words := splitIntoWords(inputStr)
		wordFrequencies := countWords(words)
		wordFrequencies = getTop10WordFrequencies(wordFrequencies)
		res = getWords(wordFrequencies)
	}

	return res
}

func splitIntoWords(text string) []string {
	return sep.Split(text, -1)
}

func countWords(words []string) []wordFrequency {
	wordCountMap := map[string]int{}
	for _, word := range words {
		wordCountMap[strings.ToLower(word)]++
	}

	var wordFrequencies []wordFrequency
	for word, count := range wordCountMap {
		wordFrequencies = append(wordFrequencies, wordFrequency{word, count})
	}

	return wordFrequencies
}

func getTop10WordFrequencies(wordFrequencies []wordFrequency) []wordFrequency {
	sort.Slice(wordFrequencies, func(left, right int) bool {
		// First we sort by count, but in case it's equal we sort by word
		if wordFrequencies[left].count > wordFrequencies[right].count {
			return true
		}
		if wordFrequencies[left].count < wordFrequencies[right].count {
			return false
		}
		return wordFrequencies[left].word < wordFrequencies[right].word
	})

	var rightBorder int
	if len(wordFrequencies) < 10 {
		rightBorder = len(wordFrequencies)
	} else {
		rightBorder = 10
	}
	return wordFrequencies[:rightBorder]
}

func getWords(wordFrequencies []wordFrequency) (words []string) {
	for _, wc := range wordFrequencies {
		words = append(words, wc.word)
	}
	return
}
