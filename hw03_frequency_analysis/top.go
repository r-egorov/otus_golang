package hw03frequencyanalysis

import (
	"fmt"
	"sort"
	"strings"
)

type wordCounter struct {
	word  string
	count int
}

func Top10(inputStr string) []string {
	words := strings.Fields(inputStr)

	if len(words) == 0 {
		return []string{}
	}

	wordCount := map[string]int{}
	for _, word := range words {
		wordCount[word]++
	}

	var wordCounters []wordCounter
	for word, count := range wordCount {
		wordCounters = append(wordCounters, wordCounter{word, count})
	}

	sort.Slice(wordCounters, func(left, right int) bool {
		return wordCounters[left].count > wordCounters[right].count
	})

	wordCounters = wordCounters[:10]

	fmt.Println(wordCounters)
	for i := 0; i < 10; i++ {
		if i < 9 && wordCounters[i+1].count == wordCounters[i].count {
			j := i
			for ; j < 9 && wordCounters[j].count == wordCounters[j+1].count; j++ {
			}
			subSlice := wordCounters[i : j+1]
			sort.Slice(subSlice, func(left, right int) bool {
				return subSlice[left].word < subSlice[right].word
			})
			i = j
		}
	}

	var topWords []string
	for _, wc := range wordCounters {
		topWords = append(topWords, wc.word)
	}

	return topWords
}
