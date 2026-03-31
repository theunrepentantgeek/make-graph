package namespace

import (
	"regexp"
	"strings"

	"github.com/rotisserie/eris"
)

// delimiters are the characters treated as namespace separators in Makefile target names.
var delimiters = "-."

// Namespace returns the namespace portion of a node ID.
// Returns everything before the last "-" or ".".
// Returns "" if no delimiter is found.
func Namespace(id string) string {
	idx := strings.LastIndexAny(id, delimiters)
	if idx < 0 {
		return ""
	}

	return id[:idx]
}

// Parent returns the parent of a namespace string.
// Returns everything before the last "-" or ".".
// Returns "" if the namespace has no parent.
func Parent(ns string) string {
	idx := strings.LastIndexAny(ns, delimiters)
	if idx < 0 {
		return ""
	}

	return ns[:idx]
}

// Depth returns the nesting depth of a namespace.
// Counts hyphens and dots within the namespace.
// A top-level namespace (no internal delimiters) has depth 0.
func Depth(ns string) int {
	count := 0

	for _, c := range ns {
		if strings.ContainsRune(delimiters, c) {
			count++
		}
	}

	return count
}

// CompileMatchPattern converts a glob-style pattern (using *, ?, and [...])
// to a compiled regexp. Returns an error if the resulting regex is invalid.
func CompileMatchPattern(pattern string) (*regexp.Regexp, error) {
	var b strings.Builder

	b.WriteString("^")

	inBracket := false

	for i := range len(pattern) {
		c := pattern[i]
		inBracket = convertGlobChar(&b, c, inBracket)
	}

	b.WriteString("$")

	re, err := regexp.Compile(b.String())
	if err != nil {
		return nil, eris.Wrapf(err, "failed to compile pattern %q", pattern)
	}

	return re, nil
}

// convertGlobChar writes one glob character to the regex builder and returns
// the updated bracket state.
func convertGlobChar(b *strings.Builder, c byte, inBracket bool) bool {
	switch {
	case c == '[' && !inBracket:
		b.WriteByte(c)

		return true
	case c == ']' && inBracket:
		b.WriteByte(c)

		return false
	case inBracket:
		b.WriteByte(c)
	case c == '*':
		b.WriteString(".*")
	case c == '?':
		b.WriteByte('.')
	default:
		b.WriteString(regexp.QuoteMeta(string(c)))
	}

	return inBracket
}

// MatchPattern returns a glob-style pattern string matching all nodes
// in the given namespace. The returned pattern is intended for storage
// in NodeStyleRule.Match and will be compiled via CompileMatchPattern.
//
// Returns "ns[-.]*" to match all delimiter styles.
func MatchPattern(ns string) string {
	return ns + "[-.]*"
}
