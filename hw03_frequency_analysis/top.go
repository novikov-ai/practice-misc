package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

type UniqueWord struct {
	Name   string
	Amount int
}

func Top10(s string) []string {

	const topBoundNumber = 10

	inputSplit := strings.Fields(s)

	words := make(map[string]int)

	for _, v := range inputSplit {

		elem, ok := words[v]
		if ok {
			words[v] = elem + 1
		} else {
			words[v] = 1
		}
	}

	if len(words) == 0 {
		return nil
	}

	sortedWords := make([]UniqueWord, 0)
	for key, value := range words {
		sortedWords = append(sortedWords, UniqueWord{
			Name:   key,
			Amount: value})
	}

	sort.Slice(sortedWords, func(i, j int) bool {

		isAmountEqual := sortedWords[i].Amount == sortedWords[j].Amount
		if isAmountEqual {
			return sortedWords[i].Name < sortedWords[j].Name
		}

		return sortedWords[i].Amount > sortedWords[j].Amount
	})

	result := make([]string, topBoundNumber)

	i := 0
	for _, value := range sortedWords {
		if i < topBoundNumber {
			result[i] = value.Name
			i++
		}
	}

	return result
}
