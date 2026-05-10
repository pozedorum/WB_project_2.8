package anagrams

import (
	"reflect"
	"testing"
)

func TestAnagramsBasic(t *testing.T) {
	input := []string{
		"пятак", "пятка", "тяпка",
		"листок", "слиток", "столик",
		"стол",
	}

	got := Anagrams(input)

	want := map[string][]string{
		"пятак":  {"пятак", "пятка", "тяпка"},
		"листок": {"листок", "слиток", "столик"},
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %#v, want %#v", got, want)
	}
}

func TestAnagramsIgnoresSingleWords(t *testing.T) {
	input := []string{"кот", "дом", "стол"}

	got := Anagrams(input)

	want := map[string][]string{}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %#v, want %#v", got, want)
	}
}

func TestAnagramsCaseInsensitive(t *testing.T) {
	input := []string{"Пятак", "пятка", "ТЯПКА"}

	got := Anagrams(input)

	want := map[string][]string{
		"пятак": {"пятак", "пятка", "тяпка"},
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %#v, want %#v", got, want)
	}
}

func TestAnagramsEmptyInput(t *testing.T) {
	got := Anagrams(nil)

	want := map[string][]string{}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %#v, want %#v", got, want)
	}
}
