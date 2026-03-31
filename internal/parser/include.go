package parser

import (
	"path/filepath"
	"strings"

	"github.com/rotisserie/eris"
)

// handleInclude processes include, -include, and sinclude directives.
// Returns true if the line was an include directive, false otherwise.
func handleInclude(line, currentFile string, visited map[string]bool, mf *Makefile) (bool, error) {
	trimmed := strings.TrimSpace(line)

	silent := false

	var rest string

	switch {
	case strings.HasPrefix(trimmed, "-include "):
		silent = true
		rest = strings.TrimPrefix(trimmed, "-include ")
	case strings.HasPrefix(trimmed, "sinclude "):
		silent = true
		rest = strings.TrimPrefix(trimmed, "sinclude ")
	case strings.HasPrefix(trimmed, "include "):
		rest = strings.TrimPrefix(trimmed, "include ")
	default:
		return false, nil
	}

	dir := filepath.Dir(currentFile)

	for _, path := range strings.Fields(rest) {
		includePath := path
		if !filepath.IsAbs(includePath) {
			includePath = filepath.Join(dir, includePath)
		}

		included, err := parseFile(includePath, visited)
		if err != nil {
			if silent {
				continue
			}

			return true, eris.Wrapf(err, "including %s", path)
		}

		mf.Rules = append(mf.Rules, included.Rules...)
	}

	return true, nil
}
