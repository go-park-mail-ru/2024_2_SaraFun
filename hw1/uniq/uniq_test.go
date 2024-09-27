package main

import (
	"bytes"

	"testing"

	"github.com/stretchr/testify/require"
)

func TestProcessFile(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expected       string
		countFlag      bool
		duplicatesFlag bool
		uniqueFlag     bool
		ignoreCase     bool
		fieldCount     int
		charCount      int
	}{
		{
			name:           "Basic unique lines",
			input:          "I love music.\nI love music.\nI love music.\n\nI love music of Kartik.\nI love music of Kartik.\nThanks.\nI love music of Kartik.\nI love music of Kartik.\n",
			expected:       "I love music.\n\nI love music of Kartik.\nThanks.\nI love music of Kartik.\n",
			countFlag:      false,
			duplicatesFlag: false,
			uniqueFlag:     false,
			ignoreCase:     false,
			fieldCount:     0,
			charCount:      0,
		},
		{
			name:           "Count occurrences -c",
			input:          "I love music.\nI love music.\nI love music.\n\nI love music of Kartik.\nI love music of Kartik.\nThanks.\nI love music of Kartik.\nI love music of Kartik.\n",
			expected:       "3 I love music.\n1 \n2 I love music of Kartik.\n1 Thanks.\n2 I love music of Kartik.\n",
			countFlag:      true,
			duplicatesFlag: false,
			uniqueFlag:     false,
			ignoreCase:     false,
			fieldCount:     0,
			charCount:      0,
		},
		{
			name:           "Only print duplicate lines -d",
			input:          "I love music.\nI love music.\n\nI love music of Kartik.\nI love music of Kartik.\nThanks.\nI love music of Kartik.\n",
			expected:       "I love music.\nI love music of Kartik.\n",
			countFlag:      false,
			duplicatesFlag: true,
			uniqueFlag:     false,
			ignoreCase:     false,
			fieldCount:     0,
			charCount:      0,
		},
		{
			name:           "Only print unique lines -u",
			input:          "I love music.\nI love music.\n\nI love music of Kartik.\nThanks.\nI love music of Kartik.\n",
			expected:       "\nI love music of Kartik.\nThanks.\nI love music of Kartik.\n",
			countFlag:      false,
			duplicatesFlag: false,
			uniqueFlag:     true,
			ignoreCase:     false,
			fieldCount:     0,
			charCount:      0,
		},
		{
			name:           "test with -f 1",
			input:          "We love music.\nI love music.\nThey love music.\n\nI love music of Kartik.\nWe love music of Kartik.\nThanks.",
			expected:       "We love music.\n\nI love music of Kartik.\nThanks.\n",
			countFlag:      false,
			duplicatesFlag: false,
			uniqueFlag:     false,
			ignoreCase:     false,
			fieldCount:     1,
			charCount:      0,
		},
		{
			name:           "test with -s 1",
			input:          "I love music.\nA love music.\nC love music.\n\nI love music of Kartik.\nWe love music of Kartik.\nThanks.",
			expected:       "I love music.\n\nI love music of Kartik.\nWe love music of Kartik.\nThanks.\n",
			countFlag:      false,
			duplicatesFlag: false,
			uniqueFlag:     false,
			ignoreCase:     false,
			fieldCount:     0,
			charCount:      1,
		},
		{
			name:           "test with [-f 1 -s 1]",
			input:          "I 1ove music.\nA 2ove music.\nC 3ove music.\n\nI 4ove music of Kartik.\nWe 5ove music of Kartik.\nThanks.",
			expected:       "I 1ove music.\n\nI 4ove music of Kartik.\nThanks.\n",
			countFlag:      false,
			duplicatesFlag: false,
			uniqueFlag:     false,
			ignoreCase:     false,
			fieldCount:     1,
			charCount:      1,
		},
		{
			name:           "test with [-f 1 -s 1 -c]",
			input:          "I 1ove music.\nA 2ove music.\nC 3ove music.\n\nI 4ove music of Kartik.\nWe 5ove music of Kartik.\nThanks.",
			expected:       "3 I 1ove music.\n1 \n2 I 4ove music of Kartik.\n1 Thanks.\n",
			countFlag:      true,
			duplicatesFlag: false,
			uniqueFlag:     false,
			ignoreCase:     false,
			fieldCount:     1,
			charCount:      1,
		},
		{
			name:           "Ignore case -i",
			input:          "I LOVE MUSIC.\nI love music.\nI LoVe MuSiC.\n\nI love MuSIC of Kartik.\nI love music of kartik.\nThanks.\nI love music of kartik.\nI love MuSIC of Kartik.",
			expected:       "I LOVE MUSIC.\n\nI love MuSIC of Kartik.\nThanks.\nI love music of kartik.\n",
			countFlag:      false,
			duplicatesFlag: false,
			uniqueFlag:     false,
			ignoreCase:     true,
			fieldCount:     0,
			charCount:      0,
		},
		{
			name:           "Options -c, -d, and -u are mutually exclusive.",
			input:          "I LOVE MUSIC.\nI love music.\nI LoVe MuSiC.\n\nI love MuSIC of Kartik.\nI love music of kartik.\nThanks.\nI love music of kartik.\nI love MuSIC of Kartik.",
			expected:       "Error: Options -c, -d, and -u are mutually exclusive.\n",
			countFlag:      true,
			duplicatesFlag: true,
			uniqueFlag:     true,
			ignoreCase:     true,
			fieldCount:     0,
			charCount:      0,
		},
		{
			name:           "Options -c, -d, and -u are mutually exclusive.",
			input:          "I LOVE MUSIC.\nI love music.\nI LoVe MuSiC.\n\nI love MuSIC of Kartik.\nI love music of kartik.\nThanks.\nI love music of kartik.\nI love MuSIC of Kartik.",
			expected:       "Error: Options -c, -d, and -u are mutually exclusive.\n",
			countFlag:      true,
			duplicatesFlag: false,
			uniqueFlag:     true,
			ignoreCase:     true,
			fieldCount:     0,
			charCount:      0,
		},
		{
			name:           "Options -c, -d, and -u are mutually exclusive.",
			input:          "I LOVE MUSIC.\nI love music.\nI LoVe MuSiC.\n\nI love MuSIC of Kartik.\nI love music of kartik.\nThanks.\nI love music of kartik.\nI love MuSIC of Kartik.",
			expected:       "Error: Options -c, -d, and -u are mutually exclusive.\n",
			countFlag:      true,
			duplicatesFlag: true,
			uniqueFlag:     false,
			ignoreCase:     true,
			fieldCount:     0,
			charCount:      0,
		},
		{
			name:           "Empty",
			input:          "",
			expected:       "",
			countFlag:      false,
			duplicatesFlag: false,
			uniqueFlag:     false,
			ignoreCase:     false,
			fieldCount:     0,
			charCount:      0,
		},
		{
			name:           "1 word",
			input:          "90",
			expected:       "90\n",
			countFlag:      false,
			duplicatesFlag: false,
			uniqueFlag:     false,
			ignoreCase:     false,
			fieldCount:     0,
			charCount:      0,
		},
		{
			name:           "1 word and f 10 s 10",
			input:          "90",
			expected:       "90\n",
			countFlag:      false,
			duplicatesFlag: false,
			uniqueFlag:     false,
			ignoreCase:     false,
			fieldCount:     10,
			charCount:      10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inputReader := bytes.NewReader([]byte(tt.input))
			var outputBuffer bytes.Buffer

			processFile(inputReader, &outputBuffer, tt.countFlag, tt.duplicatesFlag, tt.uniqueFlag, tt.ignoreCase, tt.fieldCount, tt.charCount)

			require.Equal(t, tt.expected, outputBuffer.String())
		})
	}
}
