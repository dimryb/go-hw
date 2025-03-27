package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatProgressBar(t *testing.T) {
	tests := []struct {
		name     string
		percent  int
		expected string
	}{
		{
			name:     "0%",
			percent:  0,
			expected: "\rProgress: [>_________________________________________________] 0%",
		},
		{
			name:     "25%",
			percent:  25,
			expected: "\rProgress: [============>_____________________________________] 25%",
		},
		{
			name:     "50%",
			percent:  50,
			expected: "\rProgress: [=========================>________________________] 50%",
		},
		{
			name:     "75%",
			percent:  75,
			expected: "\rProgress: [=====================================>____________] 75%",
		},
		{
			name:     "100%",
			percent:  100,
			expected: "\rProgress: [==================================================] 100%",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatProgressBar(tt.percent)
			assert.Equal(t, tt.expected, result)
		})
	}
}
