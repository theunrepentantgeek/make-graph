package parser

import (
	"regexp"
	"strings"
)

// makeCallPattern matches common $(MAKE), ${MAKE}, and bare make invocations.
var makeCallPattern = regexp.MustCompile(
	`(?:` +
		`\$\(MAKE\)|\$\{MAKE\}|(?:^|\s)make` +
		`)` +
		`(?:\s+(?:-[sSkkinBwWjCl]\s+\S+|-[sSkkinBwWjl]\b|--[a-z-]+(?:=\S+|\s+\S+)?))*` +
		`\s+([a-zA-Z0-9_][a-zA-Z0-9_./-]*)`,
)

// DetectMakeCall scans a recipe line for a recursive make invocation
// and returns the target name if found.
func DetectMakeCall(recipeLine string) (string, bool) {
	line := strings.TrimLeft(recipeLine, "@-+ \t")

	matches := makeCallPattern.FindStringSubmatch(line)
	if len(matches) < 2 {
		return "", false
	}

	target := matches[1]
	if strings.HasPrefix(target, "$") {
		return "", false
	}

	return target, true
}
