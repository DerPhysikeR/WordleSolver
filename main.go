package main

import (
	"fmt"
	"os"
	"sort"
	"strings"
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

func rmWordsWithNonAsciiChars(words *[]string) *[]string {
	asciiWords := []string{}
	for _, word := range *words {
		skip := false
		for _, letter := range word {
			asciiValue := int(letter)
			if asciiValue < 65 || asciiValue > 90 {
				skip = true
				break
			}
		}
		if !skip {
			asciiWords = append(asciiWords, word)
		}
	}
	return &asciiWords
}

func filterWordsByLength(words *[]string, length int) *[]string {
	uppercaseWords := []string{}
	for _, word := range *words {
		if len(word) != length {
			continue
		}
		uppercaseWords = append(uppercaseWords, strings.ToUpper(word))
	}
	return &uppercaseWords
}

func cleanupWords(words *[]string, length int) *[]string {
	words = applyToWordSlice(strings.ToUpper, words)
	words = rmWordsWithNonAsciiChars(words)
	words = filterWordsByLength(words, length)
	return words
}

func createWordGameFromDictionary(path string, length int) *WordGame {
	words := readDictionary(path)
	words = cleanupWords(words, length)
	remainingWords := make([]string, len(*words))
	copy(remainingWords, *words)
	return &WordGame{words, &remainingWords}
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
		for _, word := range inverseMap[key] {
			result = append(result, word)
		}
	}
	return &result
}

func (wg *WordGame) getBestGuesses() *[]string {
	fmt.Println("start getting best guess")
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
	fmt.Println("before sorting stuff")
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

func main() {
	wg := createWordGameFromDictionary("/usr/share/dict/cracklib-small", 5)
	var guess string
	var score string
	for len(*wg.remainingWords) > 1 {
		fmt.Printf("Best guesses: %v\n", (*wg.getBestGuesses())[:10])

		fmt.Printf("Your guess: ")
		fmt.Scanln(&guess)

		fmt.Printf("score: ")
		fmt.Scanln(&score)

		if score == "HHHHH" {
			os.Exit(0)
		}
		wg.guess(guess, score)
	}
	fmt.Printf("The solution is: %v\n", (*wg.remainingWords)[0])
}
