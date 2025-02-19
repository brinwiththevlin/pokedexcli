package main

import "testing"

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "  HellO  world  ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "This is a new sentence ",
			expected: []string{"this", "is", "a", "new", "sentence"},
		},
		{
			input:    "",
			expected: []string{},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)
		if len(actual) != len(c.expected) {
			t.Errorf("length of result doesn't match: %v %v", len(actual), len(c.expected))
			t.Fail()
		}

		for i := range actual {
			if actual[i] != c.expected[i] {
				t.Errorf("words don't match: %v != %v", actual[i], c.expected[i])
				t.Fail()
			}
		}
	}
}
