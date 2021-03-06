package main

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"unicode"
)

type wordFunc func(string) string

type WordGame struct {
	allWords       *[]string
	remainingWords *[]string
}

func readDictionary(path string) *[]string {
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	words := strings.Split(strings.TrimSpace(string(data)), "\n")
	return &words
}

func applyToWordSlice(f wordFunc, words *[]string) *[]string {
	newWords := []string{}
	for _, word := range *words {
		w := f(word)
		if w == "" {
			continue
		}
		newWords = append(newWords, w)
	}
	return &newWords
}

func hasNoSpecialCharacters(word string) string {
	for _, letter := range word {
		if !unicode.IsLetter(letter) {
			return ""
		}
	}
	return word
}

func hasLength(word string, length int) string {
	if len(word) == length {
		return word
	}
	return ""
}

func cleanupWords(words *[]string, length int) *[]string {
	words = applyToWordSlice(hasNoSpecialCharacters, words)
	words = applyToWordSlice(func(word string) string { return hasLength(word, length) }, words)
	words = applyToWordSlice(strings.ToUpper, words)
	return words
}

func createWordGameFromDictionary(path string, length int) *WordGame {
	words := readDictionary(path)
	return createWordGame(words, length)
}

func createWordGame(words *[]string, length int) *WordGame {
	words = cleanupWords(words, length)
	remainingWords := make([]string, len(*words))
	copy(remainingWords, *words)
	return &WordGame{words, &remainingWords}
}

func createWordGameFromWordLists(allWords *[]string, remainingWords *[]string) *WordGame {
	allWords = applyToWordSlice(strings.ToUpper, allWords)
	remainingWords = applyToWordSlice(strings.ToUpper, remainingWords)
	return &WordGame{allWords, remainingWords}
}

func scoreAgainst(guess, solution string) string {
	if len(guess) != len(solution) {
		panic(fmt.Errorf("Can't score guess '%v' against '%v' with different length", guess, solution))
	}
	score := []rune{}
	for idx, letter := range guess {
		if letter == rune(solution[idx]) {
			score = append(score, 'H')
		} else if strings.ContainsRune(solution, letter) {
			score = append(score, 'h')
		} else {
			score = append(score, '.')
		}
	}
	return string(score)
}

func getKeysSortedByValue(toSort *map[string]int) *[]string {
	inverseMap := map[int][]string{}
	for key, value := range *toSort {
		inverseMap[value] = append(inverseMap[value], key)
	}
	keys := []int{}
	for key := range inverseMap {
		keys = append(keys, key)
	}

	result := []string{}
	sort.Ints(keys)
	for _, key := range keys {
		sort.Strings(inverseMap[key])
		for _, word := range inverseMap[key] {
			result = append(result, word)
		}
	}
	return &result
}

func (wg *WordGame) getBestGuesses() *[]string {
	wordScores := map[string]int{}
	for _, word := range *wg.allWords {
		scores := map[string]int{}
		for _, solution := range *wg.remainingWords {
			score := scoreAgainst(word, solution)
			scores[score] += 1
		}
		maxCount := 0
		for _, count := range scores {
			if count > maxCount {
				maxCount = count
			}
		}
		wordScores[word] = maxCount
	}
	return getKeysSortedByValue(&wordScores)
}

func (wg *WordGame) guess(guess, score string) {
	newRemainingWords := []string{}
	for _, word := range *wg.remainingWords {
		if scoreAgainst(guess, word) == score {
			newRemainingWords = append(newRemainingWords, word)
		}
	}
	wg.remainingWords = &newRemainingWords
}

func toUniqueScore(score string) string {
	uniqueScore := []rune{}
	for _, letter := range score {
		if unicode.IsLower(letter) {
			uniqueScore = append(uniqueScore, 'h')
		} else if unicode.IsUpper(letter) {
			uniqueScore = append(uniqueScore, 'H')
		} else {
			uniqueScore = append(uniqueScore, '.')
		}
	}
	return string(uniqueScore)
}

func main() {
	length := 5
	wg := createWordGameFromWordLists(&allWords, &possibleSolutions)
	var guess string
	var score string
	fmt.Printf("Calculate best guesses ...\n")
	for len(*wg.remainingWords) > 1 {
		fmt.Printf("Best guesses: %v\n", (*wg.getBestGuesses())[:12])

		guess = ""
		for len(guess) != length {
			fmt.Printf("Your guess: ")
			fmt.Scanln(&guess)
			if len(guess) != length {
				fmt.Printf("Invalid length guess '%v'\n", guess)
			}
			guess = strings.ToUpper(guess)
		}

		score = ""
		for len(score) != length {
			fmt.Printf("Score of the guess: ")
			fmt.Scanln(&score)
			if len(score) != length {
				fmt.Printf("Invalid length score '%v'\n", score)
			}
			score = toUniqueScore(score)
		}

		if score == strings.Repeat("H", length) {
			os.Exit(0)
		}
		wg.guess(guess, score)
	}
	if len(*wg.remainingWords) == 1 {
		fmt.Printf("The solution is: %v\n", (*wg.remainingWords)[0])
	} else {
		fmt.Println("No solution found :-(")
	}
    fmt.Scanln()
}
