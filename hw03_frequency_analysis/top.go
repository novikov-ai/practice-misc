package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

const returnAmount = 10
const invalidSymbol = "-"

type UniqueWord struct {
	Name   string
	Amount int
}

func Top10(s string) []string {

	inputSplit := strings.Fields(s)

	uniqueWordsMap := make(map[string]int)

	for _, word := range inputSplit {

		if word == invalidSymbol {
			continue
		}

		word = strings.ToLower(ExcludePunctuation(word))

		elem, ok := uniqueWordsMap[word]
		if ok {
			uniqueWordsMap[word] = elem + 1
		} else {
			uniqueWordsMap[word] = 1
		}
	}

	if len(uniqueWordsMap) == 0 {
		return nil
	}

	uniqueWords := GetUniqueWordsSlice(uniqueWordsMap)
	SortDesc(uniqueWords)

	return GetTopResults(uniqueWords, returnAmount)
}

func ExcludePunctuation(word string) string {
	excludeSymbols := [...]string{".", ",", ":", ";", "!", "?"}

	for i := range excludeSymbols {
		word = strings.Trim(word, excludeSymbols[i])
	}

	return word
}

func GetUniqueWordsSlice(words map[string]int) []UniqueWord {

	uniqueWords := make([]UniqueWord, 0)

	for key, value := range words {
		uniqueWords = append(uniqueWords, UniqueWord{
			Name:   key,
			Amount: value})
	}

	return uniqueWords
}

func SortDesc(words []UniqueWord) {
	sort.Slice(words, func(i, j int) bool {

		isAmountEqual := words[i].Amount == words[j].Amount
		if isAmountEqual {
			return words[i].Name < words[j].Name
		}

		return words[i].Amount > words[j].Amount
	})
}

func GetTopResults(words []UniqueWord, returnAmount int) []string {
	result := make([]string, returnAmount)

	i := 0
	for _, value := range words {
		if i < returnAmount {
			result[i] = value.Name
			i++
		}
	}

	return result
}
