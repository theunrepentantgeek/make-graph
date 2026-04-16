package graphviz

import (
	"fmt"
	"strings"

	"github.com/theunrepentantgeek/make-graph/internal/indentwriter"
)

// record represents a record-shaped node in the Graphviz output.
type record struct {
	parts []string
}

// newRecord creates a new record with no parts.
func newRecord() *record {
	return &record{}
}

// add adds a part to the record.
func (r *record) add(text string) {
	if len(text) > 0 {
		r.parts = append(r.parts, text)
	}
}

// addf adds a formatted part to the record.
func (r *record) addf(format string, args ...any) {
	r.add(fmt.Sprintf(format, args...))
}

// addWrapped adds a part to the record, wrapping the text to the specified width.
func (r *record) addWrapped(width int, text string) {
	lines := indentwriter.WordWrap(text, width)
	r.add(strings.Join(lines, "\\n"))
}

// addWrappedf adds a formatted part to the record, wrapping the text to the specified width.
func (r *record) addWrappedf(width int, format string, args ...any) {
	r.addWrapped(width, fmt.Sprintf(format, args...))
}

// String returns the string representation of the record, which is the parts joined by " | ".
func (r *record) String() string {
	if len(r.parts) == 1 {
		return escapeRecordContent(r.parts[0])
	}

	content := strings.Join(r.parts, " | ")
	content = escapeRecordContent(content)

	return fmt.Sprintf("{%s}", content)
}

// escapeRecordContent escapes characters that have special meaning in Graphviz
// record-shaped node labels. Braces delimit record fields, so literal braces in
// content (e.g. from Makefile variable references like ${VAR}) must be escaped.
// Double-quotes are NOT escaped here because the caller (quoteString) handles that
// when writing the value into a dot attribute.
func escapeRecordContent(s string) string {
	s = strings.ReplaceAll(s, `{`, `\{`)
	s = strings.ReplaceAll(s, `}`, `\}`)

	return s
}
