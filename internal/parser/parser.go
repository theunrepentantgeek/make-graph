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
	return parseFile(filename, nil)
}

func parseFile(filename string, visited map[string]bool) (*Makefile, error) {
	if visited == nil {
		visited = make(map[string]bool)
	}

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

func parseReader(r io.Reader, filePath string, visited map[string]bool) (*Makefile, error) {
	lines, err := readLines(r)
	if err != nil {
		return nil, eris.Wrap(err, "reading lines")
	}

	mf := &Makefile{}
	var currentRule *Rule

	for _, line := range lines {
		if line == "" || strings.HasPrefix(strings.TrimSpace(line), "#") {
			currentRule = nil
			continue
		}

		if strings.HasPrefix(line, "\t") && currentRule != nil {
			recipe := strings.TrimPrefix(line, "\t")
			currentRule.Recipes = append(currentRule.Recipes, recipe)

			continue
		}

		handled, includeErr := handleInclude(line, filePath, visited, mf)
		if includeErr != nil {
			return nil, includeErr
		}

		if handled {
			currentRule = nil
			continue
		}

		rules := parseRuleLine(line)
		if len(rules) > 0 {
			for i := range rules {
				mf.Rules = append(mf.Rules, rules[i])
			}

			currentRule = &mf.Rules[len(mf.Rules)-1]

			continue
		}

		currentRule = nil
	}

	mf.Rules = mergeRules(mf.Rules)

	return mf, nil
}

// readLines reads all lines from a reader, joining continuation lines (ending with `\`).
func readLines(r io.Reader) ([]string, error) {
	var lines []string

	scanner := bufio.NewScanner(r)
	var pending string

	for scanner.Scan() {
		text := scanner.Text()
		if strings.HasSuffix(text, "\\") {
			pending += strings.TrimSuffix(text, "\\")
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
func parseRuleLine(line string) []Rule {
	colonIdx := findRuleColon(line)
	if colonIdx < 0 {
		return nil
	}

	before := line[:colonIdx]
	if containsAssignment(before) {
		return nil
	}

	// Handle double-colon rules
	after := line[colonIdx+1:]
	if strings.HasPrefix(after, ":") {
		after = after[1:]
	}

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
		}
	}

	return -1
}

// containsAssignment checks if text before a colon contains an assignment operator.
func containsAssignment(text string) bool {
	for i := 0; i < len(text); i++ {
		switch text[i] {
		case '=':
			return true
		case '?', '+':
			if i+1 < len(text) && text[i+1] == '=' {
				return true
			}
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
