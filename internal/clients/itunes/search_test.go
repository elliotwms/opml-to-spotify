package itunes

import "testing"

func TestSearch(t *testing.T) {
	res, err := Search("James O'Brien's Mystery Hour", "GB")
	if err != nil {
		t.Fatal(err)
	}

	if len(res) == 0 {
		t.Fatal("no results")
	}
}
