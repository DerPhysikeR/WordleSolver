package main

import (
	"strings"
	"testing"
)

func TestApplyToWordSlice(t *testing.T) {
	compareWordSlices := func(t testing.TB, slice1, slice2 *[]string) {
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

	t.Run("with toUpper", func(t *testing.T) {
		words := []string{"abc", "aBc"}
		capitalizedWords := applyToWordSlice(strings.ToUpper, &words)
		reference := []string{"ABC", "ABC"}
		compareWordSlices(t, capitalizedWords, &reference)
	})

	t.Run("with special characters", func(t *testing.T) {
		words := []string{"ab'c", "aB_c"}
		capitalizedWords := applyToWordSlice(strings.ToUpper, &words)
		reference := []string{"AB'C", "AB_C"}
		compareWordSlices(t, capitalizedWords, &reference)
	})

}

func TestHasLength(t *testing.T) {
	t.Run("with correct length", func(t *testing.T) {
		expect := string("abc")
		got := hasLength(string("abc"), 3)
		if expect != got {
			t.Errorf("Expected '%v' got '%v'", expect, got)
		}
	})

	t.Run("with wrong length", func(t *testing.T) {
		expect := string("")
		got := hasLength(string("abc"), 1)
		if expect != got {
			t.Errorf("Expected '%v' got '%v'", expect, got)
		}
	})
}
