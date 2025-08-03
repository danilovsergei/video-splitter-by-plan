package cmd

import (
	"testing"
)

func TestGenerateTitleFromFilename(t *testing.T) {
	testCases := []struct {
		name     string
		filename string
		expected string
	}{
		{
			name:     "with number prefix and underscore",
			filename: "01_backhand.mp4",
			expected: "backhand",
		},
		{
			name:     "with number prefix and dash",
			filename: "02 - forehand.mp4",
			expected: "forehand",
		},
		{
			name:     "without number prefix",
			filename: "serve.mp4",
			expected: "serve",
		},
		{
			name:     "with multiple words",
			filename: "03_a_long_exercise_name.mp4",
			expected: "a_long_exercise_name",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := generateTitleFromFilename(tc.filename)
			if actual != tc.expected {
				t.Errorf("expected title '%s', but got '%s'", tc.expected, actual)
			}
		})
	}
}
