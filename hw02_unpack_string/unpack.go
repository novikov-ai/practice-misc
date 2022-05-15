package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(input string) (string, error) {
	inputRunes := []rune(input)
	inputRunesCount := len(inputRunes)

	const backSlash = `\`
	isEscapedNext := false

	var unpackedInput strings.Builder

	for i := 0; i < inputRunesCount; i++ {
		currentRune := inputRunes[i]

		isCurrentRuneSlash := string(currentRune) == backSlash

		if isCurrentRuneSlash && !isEscapedNext {
			isEscapedNext = true
			continue
		}

		isCurrentRuneDigit := unicode.IsDigit(currentRune)
		isEscapedInvalid := isEscapedNext && !isCurrentRuneDigit && !isCurrentRuneSlash

		if (isCurrentRuneDigit && !isEscapedNext) || isEscapedInvalid {
			return "", ErrInvalidString
		}

		isEscapedNext = false

		repeatNumber := 1

		if i+1 < inputRunesCount {
			nextRune := inputRunes[i+1]

			skipNextRune := unicode.IsDigit(nextRune)

			if skipNextRune {
				repeatNumber, _ = strconv.Atoi(string(nextRune))
				i++
			}
		}

		repeatedValue := strings.Repeat(string(currentRune), repeatNumber)

		unpackedInput.WriteString(repeatedValue)
	}

	return unpackedInput.String(), nil
}
