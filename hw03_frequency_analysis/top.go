package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

var regexpPunctuation = regexp.MustCompile(`\.?,?;?:?!?\??\(?\)?\[?]?{?}?`)

const (
	returnAmount  = 10
	invalidSymbol = "-"
)

type UniqueWord struct {
	Name   string
	Amount int
}

func Top10(input string) []string {
	inputSplit := strings.Fields(input)

	uniqueWordsMap := make(map[string]int)

	for _, word := range inputSplit {
		if word == invalidSymbol {
			continue
		}

		word = strings.ToLower(regexpPunctuation.ReplaceAllString(word, ""))

		uniqueWordsMap[word]++
	}

	if len(uniqueWordsMap) == 0 {
		return nil
	}

	uniqueWords := GetUniqueWordsSlice(uniqueWordsMap)
	SortDesc(uniqueWords)

	return GetTopResults(uniqueWords, returnAmount)
}

func GetUniqueWordsSlice(words map[string]int) []UniqueWord {
	uniqueWords := make([]UniqueWord, 0)

	for key, value := range words {
		uniqueWords = append(uniqueWords, UniqueWord{
			Name:   key,
			Amount: value,
		})
	}

	return uniqueWords
}

func SortDesc(words []UniqueWord) {
	sort.Slice(words, func(i, j int) bool {
		if words[i].Amount == words[j].Amount {
			return words[i].Name < words[j].Name
		}

		return words[i].Amount > words[j].Amount
	})
}

func GetTopResults(words []UniqueWord, returnAmount int) []string {
	result := make([]string, returnAmount)

	for i, value := range words {
		if i < returnAmount {
			result[i] = value.Name
		}
	}

	return result
}
