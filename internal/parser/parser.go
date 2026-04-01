package parser

import (
	"bufio"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/rotisserie/eris"
)

// Parse reads a Makefile and returns its parsed representation.
func Parse(filename string) (*Makefile, error) {
	visited := make(map[string]bool)

	return parseFile(filename, visited)
}

func parseFile(filename string, visited map[string]bool) (*Makefile, error) {
	absPath, err := filepath.Abs(filename)
	if err != nil {
		return nil, eris.Wrapf(err, "resolving path %s", filename)
	}

	if visited[absPath] {
		return &Makefile{}, nil
	}

	visited[absPath] = true

	f, err := os.Open(absPath)
	if err != nil {
		return nil, eris.Wrapf(err, "opening %s", filename)
	}

	defer f.Close()

	return parseReader(f, absPath, visited)
}

//nolint:revive,cyclop // Function is already long and splitting it doesn't improve readability.
func parseReader(
	r io.Reader,
	filePath string,
	visited map[string]bool,
) (*Makefile, error) {
	lines, err := readLines(r)
	if err != nil {
		return nil, eris.Wrap(err, "reading lines")
	}

	mf := &Makefile{}

	var currentRuleIndices []int

	for _, line := range lines {
		if line == "" || strings.HasPrefix(strings.TrimSpace(line), "#") {
			continue
		}

		if strings.HasPrefix(line, "\t") {
			recipe := strings.TrimPrefix(line, "\t")
			for _, idx := range currentRuleIndices {
				mf.Rules[idx].Recipes = append(mf.Rules[idx].Recipes, recipe)
			}

			continue
		}

		handled, includeErr := handleInclude(line, filePath, visited, mf)
		if includeErr != nil {
			return nil, includeErr
		}

		if handled {
			currentRuleIndices = nil

			continue
		}

		rules := parseRuleLine(line)
		if len(rules) > 0 {
			currentRuleIndices = nil

			for i := range rules {
				mf.Rules = append(mf.Rules, rules[i])
				currentRuleIndices = append(currentRuleIndices, len(mf.Rules)-1)
			}

			continue
		}

		currentRuleIndices = nil
	}

	mf.Rules = mergeRules(mf.Rules)

	return mf, nil
}

// readLines reads all lines from a reader, joining continuation lines (ending with `\`).
func readLines(r io.Reader) ([]string, error) {
	//nolint:prealloc // We don't know how many lines there will be
	var lines []string

	scanner := bufio.NewScanner(r)

	var pending string

	for scanner.Scan() {
		text := scanner.Text()
		if before, ok := strings.CutSuffix(text, "\\"); ok {
			pending += before

			continue
		}

		if pending != "" {
			text = pending + text
			pending = ""
		}

		lines = append(lines, text)
	}

	if pending != "" {
		lines = append(lines, pending)
	}

	if err := scanner.Err(); err != nil {
		return nil, eris.Wrap(err, "scanning lines")
	}

	return lines, nil
}

// parseRuleLine attempts to parse a line as a rule definition.
// Returns nil if the line is not a rule.
//
//nolint:revive // Function is already long and splitting it doesn't improve readability.
func parseRuleLine(line string) []Rule {
	colonIdx := findRuleColon(line)
	if colonIdx < 0 {
		return nil
	}

	before := line[:colonIdx]
	if containsAssignment(before) {
		return nil
	}

	// Check for := assignment (colon immediately followed by =)
	after := line[colonIdx+1:]
	if strings.HasPrefix(after, "=") {
		return nil
	}

	// Handle double-colon rules
	after = strings.TrimPrefix(after, ":")

	targets := strings.Fields(before)
	if len(targets) == 0 {
		return nil
	}

	// Filter dot-prefix targets
	var validTargets []string

	for _, t := range targets {
		if !strings.HasPrefix(t, ".") {
			validTargets = append(validTargets, t)
		}
	}

	if len(validTargets) == 0 {
		return nil
	}

	prereqs, description := parseAfterColon(after)

	//nolint:prealloc // We don't know how many rules there will be
	var rules []Rule

	for _, target := range validTargets {
		rules = append(rules, Rule{
			Target:        target,
			Prerequisites: prereqs,
			Description:   description,
		})
	}

	return rules
}

// parseAfterColon splits the text after the colon into prerequisites and description.
func parseAfterColon(after string) ([]string, string) {
	var description string

	// Look for ## (description marker) first
	if idx := strings.Index(after, "##"); idx >= 0 {
		description = strings.TrimSpace(after[idx+2:])
		after = after[:idx]
	} else if idx := strings.Index(after, "#"); idx >= 0 {
		// Single # is a comment delimiter (not a description), discard the rest
		after = after[:idx]
	}

	fields := strings.Fields(after)

	return fields, description
}

// findRuleColon finds the index of the colon that separates targets from prerequisites.
// Returns -1 if no rule colon is found.
// Skips colons inside $(...) or ${...} expansions.
//
//nolint:revive,cyclop // Function is already long and splitting it doesn't improve readability.
func findRuleColon(line string) int {
	depth := 0

	for i := 0; i < len(line); i++ {
		switch line[i] {
		case '$':
			if i+1 < len(line) && (line[i+1] == '(' || line[i+1] == '{') {
				depth++
				i++ // skip the opening bracket
			}
		case ')', '}':
			if depth > 0 {
				depth--
			}
		case ':':
			if depth == 0 {
				return i
			}
		default:
			// Nothing
		}
	}

	return -1
}

// containsAssignment checks if text before a colon contains an assignment operator.
func containsAssignment(text string) bool {
	for i := range len(text) {
		switch text[i] {
		case '=':
			return true
		case '?', '+':
			if i+1 < len(text) && text[i+1] == '=' {
				return true
			}
		default:
			// Nothing
		}
	}

	return false
}

// mergeRules combines rules with the same target, merging prerequisites, recipes,
// and keeping the first non-empty description.
func mergeRules(rules []Rule) []Rule {
	seen := make(map[string]int)

	var result []Rule

	for _, r := range rules {
		if idx, ok := seen[r.Target]; ok {
			result[idx].Prerequisites = append(result[idx].Prerequisites, r.Prerequisites...)
			result[idx].Recipes = append(result[idx].Recipes, r.Recipes...)

			if result[idx].Description == "" && r.Description != "" {
				result[idx].Description = r.Description
			}
		} else {
			seen[r.Target] = len(result)
			result = append(result, r)
		}
	}

	return result
}
