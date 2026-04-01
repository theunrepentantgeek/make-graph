package parser

import (
	"testing"

	. "github.com/onsi/gomega"
)

//nolint:funlen // Number of test cases
func TestDetectMakeCall_VariousPatterns_ReturnsExpectedResults(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		input    string
		target   string
		expected bool
	}{
		{
			name:     "dollar-paren-make",
			input:    "$(MAKE) test",
			target:   "test",
			expected: true,
		},
		{
			name:     "dollar-brace-make",
			input:    "${MAKE} build",
			target:   "build",
			expected: true,
		},
		{
			name:     "bare-make",
			input:    "make clean",
			target:   "clean",
			expected: true,
		},
		{
			name:     "make-with-directory",
			input:    "$(MAKE) -C subdir all",
			target:   "all",
			expected: true,
		},
		{
			name:     "make-with-flags",
			input:    "$(MAKE) -s -k test",
			target:   "test",
			expected: true,
		},
		{
			name:     "at-prefix",
			input:    "@$(MAKE) build",
			target:   "build",
			expected: true,
		},
		{
			name:     "plus-prefix",
			input:    "+$(MAKE) build",
			target:   "build",
			expected: true,
		},
		{
			name:     "variable-target",
			input:    "$(MAKE) $(TARGET)",
			target:   "",
			expected: false,
		},
		{
			name:     "not-a-make-call",
			input:    "gcc -o main main.c",
			target:   "",
			expected: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			g := NewGomegaWithT(t)

			// Act
			target, found := DetectMakeCall(tc.input)

			// Assert
			g.Expect(found).To(Equal(tc.expected))
			g.Expect(target).To(Equal(tc.target))
		})
	}
}
