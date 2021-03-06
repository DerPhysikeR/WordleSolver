package main

import (
	"strings"
	"testing"
)

func compareWordSlices(t testing.TB, slice1, slice2 *[]string) {
	t.Helper()
	if len(*slice1) != len(*slice2) {
		t.Errorf("Can't compare slices of different length.")
	}
	for idx, word := range *slice1 {
		if word != (*slice2)[idx] {
			t.Errorf("Given slices are not equal '%v' != '%v'", word, (*slice2)[idx])
		}
	}
}

func TestApplyToWordSlice(t *testing.T) {
	t.Run("with toUpper on simple string", func(t *testing.T) {
		words := []string{"abc", "aBc"}
		capitalizedWords := applyToWordSlice(strings.ToUpper, &words)
		reference := []string{"ABC", "ABC"}
		compareWordSlices(t, capitalizedWords, &reference)
	})

	t.Run("with filtering function", func(t *testing.T) {
		words := []string{"abc", "aBc"}
		capitalizedWords := applyToWordSlice(func(word string) string { return "" }, &words)
		reference := []string{}
		compareWordSlices(t, capitalizedWords, &reference)
	})
}

func expectGotString(t testing.TB, expect, got string) {
	t.Helper()
	if expect != got {
		t.Errorf("Expected '%v' got '%v'", expect, got)
	}
}

func TestHasLength(t *testing.T) {
	t.Run("with correct length", func(t *testing.T) {
		expect := string("abc")
		got := hasLength(string("abc"), 3)
		expectGotString(t, expect, got)
	})

	t.Run("with wrong length", func(t *testing.T) {
		expect := string("")
		got := hasLength(string("abc"), 1)
		expectGotString(t, expect, got)
	})
}

func TestHasNoSpecialCharacters(t *testing.T) {
	t.Run("with only lowercase letters", func(t *testing.T) {
		expect := string("abc")
		got := hasNoSpecialCharacters(string("abc"))
		expectGotString(t, expect, got)
	})

	t.Run("with only uppercase letters", func(t *testing.T) {
		expect := string("ABC")
		got := hasNoSpecialCharacters(string("ABC"))
		expectGotString(t, expect, got)
	})

	t.Run("with some whitespace", func(t *testing.T) {
		expect := string("")
		got := hasNoSpecialCharacters(string("AB C"))
		expectGotString(t, expect, got)
	})

	t.Run("with some special characters", func(t *testing.T) {
		expect := string("")
		got := hasNoSpecialCharacters(string("AB*C"))
		expectGotString(t, expect, got)
	})
}

func TestCleanupWords(t *testing.T) {
	t.Run("all kinds of words", func(t *testing.T) {
		words := []string{"a,c", "a", "abc"}
		cleanedWords := cleanupWords(&words, 3)
		reference := []string{"ABC"}
		compareWordSlices(t, cleanedWords, &reference)
	})
}

func TestCreateWordGameFromWords(t *testing.T) {
	t.Run("all kinds of words", func(t *testing.T) {
		words := []string{"a,c", "a", "abc"}
		wordGame := createWordGame(&words, 3)
		reference := []string{"ABC"}
		compareWordSlices(t, wordGame.allWords, &reference)
		compareWordSlices(t, wordGame.remainingWords, &reference)
	})
}

func assertPanic(t *testing.T, f func()) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	f()
}

func TestScoreAgainst(t *testing.T) {
	t.Run("with invalid words", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("The code did not panic")
			}
		}()
		scoreAgainst("abc", "abcd")
	})

	t.Run("with identical words", func(t *testing.T) {
		score := scoreAgainst("abc", "abc")
		expectGotString(t, "HHH", score)
	})

	t.Run("with completely different words", func(t *testing.T) {
		score := scoreAgainst("abc", "xyz")
		expectGotString(t, "...", score)
	})

	t.Run("with intersecting words", func(t *testing.T) {
		score := scoreAgainst("abc", "cde")
		expectGotString(t, "..h", score)
	})

	t.Run("with similar", func(t *testing.T) {
		score := scoreAgainst("abc", "cbe")
		expectGotString(t, ".Hh", score)
	})

	t.Run("with repeating letter", func(t *testing.T) {
		score := scoreAgainst("aae", "bca")
		expectGotString(t, "hh.", score)
	})
}

func TestGetKeysSortedByValue(t *testing.T) {
	t.Run("with unique values", func(t *testing.T) {
		m := map[string]int{"abc": 1, "def": 2, "ghi": 3}
		got := getKeysSortedByValue(&m)
		reference := []string{"abc", "def", "ghi"}
		compareWordSlices(t, got, &reference)
	})

	t.Run("with non-unique values", func(t *testing.T) {
		m := map[string]int{"abc": 1, "def": 1, "ghi": 3}
		got := getKeysSortedByValue(&m)
		ref := []string{"abc", "def", "ghi"}
		compareWordSlices(t, got, &ref)
	})
}

func TestWordGameGuess(t *testing.T) {
	t.Run("are words really filtered", func(t *testing.T) {
		words := []string{"ABC", "ACB", "EAD"}
		wordGame := createWordGame(&words, 3)
		wordGame.guess("AEF", "H..")
		reference := []string{"ABC", "ACB"}
		compareWordSlices(t, wordGame.remainingWords, &reference)
	})
}

func TestWordGameGetBestGuesses(t *testing.T) {
	t.Run("with simple word list", func(t *testing.T) {
		words := []string{"AXY", "BXY", "CXY", "ABC", "XXX"}
		wordGame := createWordGame(&words, 3)
		bestGuesses := wordGame.getBestGuesses()
		best := (*bestGuesses)[0]
		expectGotString(t, "ABC", best)
		worst := (*bestGuesses)[4]
		expectGotString(t, "XXX", worst)
	})
}

func TestToUniqueScore(t *testing.T) {
	t.Run("with only dots", func(t *testing.T) {
		got := toUniqueScore("...")
		expect := "..."
		expectGotString(t, expect, got)
	})

	t.Run("with only lowercase letters", func(t *testing.T) {
		got := toUniqueScore("abc")
		expect := "hhh"
		expectGotString(t, expect, got)
	})

	t.Run("with only uppercase letters", func(t *testing.T) {
		got := toUniqueScore("ABC")
		expect := "HHH"
		expectGotString(t, expect, got)
	})

	t.Run("with mixture of all", func(t *testing.T) {
		got := toUniqueScore(".aB")
		expect := ".hH"
		expectGotString(t, expect, got)
	})

	t.Run("with different special character", func(t *testing.T) {
		got := toUniqueScore("*aB")
		expect := ".hH"
		expectGotString(t, expect, got)
	})
}
