package namespace

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestNamespace_VariousInputs_ReturnsExpectedNamespace(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		input    string
		expected string
	}{
		"hyphen delimiter":           {input: "build-bin", expected: "build"},
		"dot delimiter":              {input: "tidy.format", expected: "tidy"},
		"mixed delimiters":           {input: "build-bin.linux", expected: "build-bin"},
		"no delimiter returns empty": {input: "deploy", expected: ""},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			g := NewWithT(t)

			// Act
			ns := Namespace(c.input)

			// Assert
			g.Expect(ns).To(Equal(c.expected))
		})
	}
}

func TestParent_VariousInputs_ReturnsExpectedParent(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		input    string
		expected string
	}{
		"hyphen nested":           {input: "build-bin", expected: "build"},
		"dot nested":              {input: "tidy.format", expected: "tidy"},
		"top level returns empty": {input: "cmd", expected: ""},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			g := NewWithT(t)

			// Act
			p := Parent(c.input)

			// Assert
			g.Expect(p).To(Equal(c.expected))
		})
	}
}

func TestDepth_VariousNamespaces_ReturnsCorrectDepth(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		ns       string
		expected int
	}{
		"top-level":     {ns: "build", expected: 0},
		"nested hyphen": {ns: "build-bin", expected: 1},
		"nested dot":    {ns: "tidy.format", expected: 1},
		"deep mixed":    {ns: "build-bin.linux", expected: 2},
		"no delimiter":  {ns: "cmd", expected: 0},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			g := NewWithT(t)

			// Act
			d := Depth(c.ns)

			// Assert
			g.Expect(d).To(Equal(c.expected))
		})
	}
}

func TestCompileMatchPattern_VariousPatterns_ConvertsToExpectedRegex(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		pattern  string
		expected string
	}{
		"glob star":     {pattern: "build*", expected: "^build.*$"},
		"glob question": {pattern: "build?", expected: "^build.$"},
		"special chars": {pattern: "a.b(c)", expected: `^a\.b\(c\)$`},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			g := NewWithT(t)

			// Act
			re, err := CompileMatchPattern(c.pattern)

			// Assert
			g.Expect(err).NotTo(HaveOccurred())
			g.Expect(re.String()).To(Equal(c.expected))
		})
	}
}

func TestCompileMatchPattern_CharacterClass_PassedThrough(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	// Act
	re, err := CompileMatchPattern("build[-.]*")

	// Assert
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(re.String()).To(Equal(`^build[-.].*$`))
	g.Expect(re.MatchString("build-bin")).To(BeTrue())
	g.Expect(re.MatchString("build.image")).To(BeTrue())
	g.Expect(re.MatchString("buildall")).To(BeFalse())
}

func TestCompileMatchPattern_InvalidPattern_ReturnsError(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	// Act — unclosed bracket
	_, err := CompileMatchPattern("[unclosed")

	// Assert
	g.Expect(err).To(HaveOccurred())
}

func TestCompileMatchPattern_ExistingPatterns_BackwardCompatible(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		pattern string
		match   string
		noMatch string
	}{
		"prefix glob": {
			pattern: "build*",
			match:   "build-bin",
			noMatch: "test",
		},
		"contains glob": {
			pattern: "*test*",
			match:   "mytest-unit",
			noMatch: "build",
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			g := NewWithT(t)

			// Act
			re, err := CompileMatchPattern(c.pattern)

			// Assert
			g.Expect(err).NotTo(HaveOccurred())
			g.Expect(re.MatchString(c.match)).To(BeTrue())
			g.Expect(re.MatchString(c.noMatch)).To(BeFalse())
		})
	}
}

func TestMatchPattern_Namespace_ReturnsDelimiterPattern(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	// Act
	p := MatchPattern("build")

	// Assert
	g.Expect(p).To(Equal("build[-.]*"))
}
