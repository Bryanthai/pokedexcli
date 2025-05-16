package main

import (
	"testing"
)


func TestCleanInput (t *testing.T) {
	cases := []struct{
		input string
		expected []string
	}{
		{
			input: "  Hello   World",
			expected: []string{"hello","world"},
		},
		{
			input: "",
			expected: []string{},
		},
		{
			input: "   ",
			expected: []string{},
		},
		{
			input: " real",
			expected: []string{"real"},
		},
		{
			input: "  Vaporeon is my favorite pokemon",
			expected: []string{"vaporeon","is","my","favorite","pokemon"},
		},
	}
	for j, c := range cases {
		t.Logf("Case %v", j)
		actual := cleanInput(c.input)
		if len(actual) != len(c.expected) {
			t.Errorf("Inequal length of slices %v and %v", len(actual), len(c.expected))
			t.Fail()
			continue
		}
		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			if word != expectedWord {
				t.Errorf("Incorrect word at index: %v", i)
				t.Fail()
			}
		}
	}
	t.Logf("Passed!")
}