# make-graph Implementation Plan

> **For agentic workers:** REQUIRED: Use superpowers:subagent-driven-development (if subagents available) or superpowers:executing-plans to implement this plan. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build a Go CLI tool that parses Makefiles and generates Graphviz/Mermaid dependency graphs, mirroring task-graph's architecture.

**Architecture:** Parser reads Makefile → builds intermediate `parser.Makefile` → `makegraph.Builder` converts to `graph.Graph` → output formatters generate `.dot` or Mermaid files. Packages copied from task-graph: graph, graphviz, mermaid, config, dot, indentwriter, safe, namespace, autocolor. Novel packages: parser, makegraph.

**Tech Stack:** Go 1.22+, kong (CLI), eris (errors), gomega (test assertions), goldie (golden file tests), console-slog (logging), yaml.v3 (config)

**Spec:** `docs/superpowers/specs/2026-03-31-make-graph-design.md`

---

## Chunk 1: Project Scaffolding & Core Packages

### Task 1: Initialize Go module and project files

**Files:**
- Create: `go.mod`
- Create: `main.go`
- Create: `DEVELOPMENT.md`
- Create: `TODO.md`

- [ ] **Step 1: Initialize Go module**

Run:
```bash
cd /home/bevan/github/make-graph
go mod init github.com/theunrepentantgeek/make-graph
```

- [ ] **Step 2: Create main.go**

```go
package main

import (
	"os"

	"github.com/alecthomas/kong"
	"github.com/rotisserie/eris"

	"github.com/theunrepentantgeek/make-graph/internal/cmd"
)

func main() {
	var cli cmd.CLI

	ctx := kong.Parse(
		&cli,
		kong.UsageOnError(),
	)

	log := cli.CreateLogger()
	cfg, err := cli.CreateConfig()
	if err != nil {
		log.Error(eris.ToString(err, true))
		ctx.Exit(1)
	}

	flags := &cmd.Flags{
		Verbose: cli.Verbose,
		Log:     log,
		Config:  cfg,
	}

	err = ctx.Run(flags)
	if err != nil {
		log.Error(eris.ToString(err, true))
		ctx.Exit(1)
	}

	os.Exit(0)
}
```

- [ ] **Step 3: Copy DEVELOPMENT.md from task-graph**

Copy `/home/bevan/github/task-graph/DEVELOPMENT.md` to `/home/bevan/github/make-graph/DEVELOPMENT.md`. Replace all references to `task-graph` with `make-graph` and `Taskfile` with `Makefile`.

- [ ] **Step 4: Create TODO.md**

```markdown
# TODO

Items to create as GitHub issues when the repository is published.

## Parser Limitations

- [ ] $(MAKE) detection: target is a variable (`$(MAKE) $(GOALS)`) — no variable expansion attempted
- [ ] $(MAKE) detection: complex commands with pipes or redirects (`$(MAKE) target | grep error`)
- [ ] $(MAKE) detection: loops (`for t in $(TARGETS); do $(MAKE) $$t; done`)
- [ ] $(MAKE) detection: multiple make invocations on one line (only the first is detected)
```

- [ ] **Step 5: Commit**

```bash
git add -A
git commit -m "Initialize Go module and project scaffold"
```

### Task 2: Copy graph package

**Files:**
- Create: `internal/graph/graph.go`
- Create: `internal/graph/node.go`
- Create: `internal/graph/edge.go`
- Create: `internal/graph/graph_test.go`
- Create: `internal/graph/node_test.go`
- Create: `internal/graph/edge_test.go`

- [ ] **Step 1: Copy graph package files from task-graph**

Copy all files from `/home/bevan/github/task-graph/internal/graph/` to `/home/bevan/github/make-graph/internal/graph/`. In each file, replace the module path `github.com/theunrepentantgeek/task-graph` with `github.com/theunrepentantgeek/make-graph`.

- [ ] **Step 2: Run tests to verify**

Run: `go test ./internal/graph/...`
Expected: All tests pass.

- [ ] **Step 3: Commit**

```bash
git add -A
git commit -m "Add graph package (copied from task-graph)"
```

### Task 3: Copy indentwriter package

**Files:**
- Create: `internal/indentwriter/indent_writer.go`
- Create: `internal/indentwriter/line.go`
- Create: `internal/indentwriter/indent_writer_test.go`
- Copy any other files in the indentwriter directory.

- [ ] **Step 1: Copy indentwriter package files from task-graph**

Copy all files from `/home/bevan/github/task-graph/internal/indentwriter/` to `/home/bevan/github/make-graph/internal/indentwriter/`. Replace module paths.

- [ ] **Step 2: Run tests to verify**

Run: `go test ./internal/indentwriter/...`
Expected: All tests pass.

- [ ] **Step 3: Commit**

```bash
git add -A
git commit -m "Add indentwriter package (copied from task-graph)"
```

### Task 4: Copy safe package

**Files:**
- Create: `internal/safe/registry.go`
- Create: `internal/safe/label.go`
- Copy any test files.

- [ ] **Step 1: Copy safe package files from task-graph**

Copy all files from `/home/bevan/github/task-graph/internal/safe/` to `/home/bevan/github/make-graph/internal/safe/`. Replace module paths.

- [ ] **Step 2: Run tests to verify**

Run: `go test ./internal/safe/...`
Expected: All tests pass.

- [ ] **Step 3: Commit**

```bash
git add -A
git commit -m "Add safe package (copied from task-graph)"
```

### Task 5: Copy config package

**Files:**
- Create: `internal/config/config.go`
- Create: `internal/config/graphviz.go`
- Create: `internal/config/mermaid.go`
- Create: `internal/config/node_style.go`

- [ ] **Step 1: Copy config package files from task-graph**

Copy all files from `/home/bevan/github/task-graph/internal/config/` to `/home/bevan/github/make-graph/internal/config/`. Replace module paths. Replace any references to "task" or "Taskfile" in comments with "make" or "Makefile" where appropriate.

- [ ] **Step 2: Verify compilation**

Run: `go build ./internal/config/...`
Expected: Compiles successfully.

- [ ] **Step 3: Commit**

```bash
git add -A
git commit -m "Add config package (copied from task-graph)"
```

### Task 6: Copy namespace package and adapt for Makefile conventions

**Files:**
- Create: `internal/namespace/namespace.go`
- Create: `internal/namespace/namespace_test.go` (if exists in task-graph)

- [ ] **Step 1: Copy namespace package from task-graph**

Copy all files from `/home/bevan/github/task-graph/internal/namespace/` to `/home/bevan/github/make-graph/internal/namespace/`. Replace module paths.

- [ ] **Step 2: Remove colon as formal delimiter**

Modify `namespace.go`: Change the delimiter handling so that colon `:` is NOT used as a namespace delimiter (Makefiles use colons in rule syntax). Also ensure slash `/` is NOT a delimiter (Makefiles use slashes in file paths). Only hyphen `-` and dot `.` should be informal delimiters.

Specifically:
- Remove or empty the `formalDelimiter` constant (set to `""` or remove the tier-1 logic)
- Keep `informalDelimiters` as `"-."` (hyphen and dot only, no slash)

- [ ] **Step 3: Update or add tests for Makefile-specific delimiter behavior**

Add test cases:
- `build-docker` → namespace `build`
- `build.docker` → namespace `build`
- `build/output.o` → namespace `` (empty — slash NOT a delimiter)
- `all` → namespace `` (no delimiter)
- `build:docker` → namespace `` (empty — colon NOT a delimiter)

- [ ] **Step 4: Run tests**

Run: `go test ./internal/namespace/...`
Expected: All tests pass.

- [ ] **Step 5: Commit**

```bash
git add -A
git commit -m "Add namespace package, adapted for Makefile conventions"
```

### Task 7: Copy autocolor package

**Files:**
- Create: `internal/autocolor/autocolor.go`
- Create: `internal/autocolor/autocolor_test.go`

- [ ] **Step 1: Copy autocolor package from task-graph**

Copy all files from `/home/bevan/github/task-graph/internal/autocolor/` to `/home/bevan/github/make-graph/internal/autocolor/`. Replace module paths.

- [ ] **Step 2: Run tests**

Run: `go test ./internal/autocolor/...`
Expected: All tests pass.

- [ ] **Step 3: Commit**

```bash
git add -A
git commit -m "Add autocolor package (copied from task-graph)"
```

### Task 8: Copy dot package

**Files:**
- Create: `internal/dot/dot.go`
- Create: `internal/dot/dot_test.go`

- [ ] **Step 1: Copy dot package from task-graph**

Copy all files from `/home/bevan/github/task-graph/internal/dot/` to `/home/bevan/github/make-graph/internal/dot/`. Replace module paths.

- [ ] **Step 2: Run tests**

Run: `go test ./internal/dot/...`
Expected: All tests pass.

- [ ] **Step 3: Commit**

```bash
git add -A
git commit -m "Add dot package (copied from task-graph)"
```

### Task 9: Copy graphviz package

**Files:**
- Create: all files from `internal/graphviz/` including test files and testdata

- [ ] **Step 1: Copy graphviz package from task-graph**

Copy all files from `/home/bevan/github/task-graph/internal/graphviz/` to `/home/bevan/github/make-graph/internal/graphviz/`. Replace module paths. Replace "task" references in comments with "target" where appropriate (e.g., "task nodes" → "target nodes").

- [ ] **Step 2: Run tests**

Run: `go test ./internal/graphviz/...`
Expected: All tests pass.

- [ ] **Step 3: Commit**

```bash
git add -A
git commit -m "Add graphviz package (copied from task-graph)"
```

### Task 10: Copy mermaid package

**Files:**
- Create: all files from `internal/mermaid/` including test files and testdata

- [ ] **Step 1: Copy mermaid package from task-graph**

Copy all files from `/home/bevan/github/task-graph/internal/mermaid/` to `/home/bevan/github/make-graph/internal/mermaid/`. Replace module paths. Replace "task" references in comments with "target" where appropriate.

- [ ] **Step 2: Run tests**

Run: `go test ./internal/mermaid/...`
Expected: All tests pass.

- [ ] **Step 3: Commit**

```bash
git add -A
git commit -m "Add mermaid package (copied from task-graph)"
```

### Task 11: Resolve dependencies and verify full build

- [ ] **Step 1: Add all dependencies**

Run:
```bash
cd /home/bevan/github/make-graph
go mod tidy
```

- [ ] **Step 2: Run all tests**

Run: `go test ./...`
Expected: All copied package tests pass.

- [ ] **Step 3: Commit**

```bash
git add -A
git commit -m "Resolve dependencies, all copied packages pass tests"
```

---

## Chunk 2: Makefile Parser

### Task 12: Create parsed model types

**Files:**
- Create: `internal/parser/makefile.go`
- Create: `internal/parser/rule.go`

- [ ] **Step 1: Create makefile.go with the Makefile model**

```go
package parser

// Makefile represents a parsed Makefile.
type Makefile struct {
	Rules []Rule
}
```

- [ ] **Step 2: Create rule.go with the Rule model**

```go
package parser

// Rule represents a single Makefile target rule.
type Rule struct {
	Target        string
	Prerequisites []string
	Description   string
	Recipes       []string
}
```

- [ ] **Step 3: Verify compilation**

Run: `go build ./internal/parser/...`
Expected: Compiles.

- [ ] **Step 4: Commit**

```bash
git add -A
git commit -m "Add parser model types: Makefile and Rule"
```

### Task 13: Write parser tests for basic rule parsing

**Files:**
- Create: `internal/parser/parser.go` (minimal stub)
- Create: `internal/parser/parser_test.go`
- Create: `internal/parser/testdata/basic-rules.mk`

- [ ] **Step 1: Create test Makefile for basic rules**

Create `internal/parser/testdata/basic-rules.mk`:
```makefile
all: build test

build:
	go build ./...

test: build
	go test ./...

clean:
	rm -rf build/
```

- [ ] **Step 2: Write failing test for basic rule parsing**

```go
package parser

import (
	"path/filepath"
	"testing"

	. "github.com/onsi/gomega"
)

func TestParse_BasicRules_ReturnsExpectedTargets(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	// Act
	mf, err := Parse(filepath.Join("testdata", "basic-rules.mk"))

	// Assert
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(mf.Rules).To(HaveLen(4))

	targets := make([]string, len(mf.Rules))
	for i, r := range mf.Rules {
		targets[i] = r.Target
	}

	g.Expect(targets).To(ConsistOf("all", "build", "test", "clean"))
}

func TestParse_BasicRules_ReturnsExpectedPrerequisites(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	// Act
	mf, err := Parse(filepath.Join("testdata", "basic-rules.mk"))

	// Assert
	g.Expect(err).NotTo(HaveOccurred())

	ruleMap := makeRuleMap(mf)
	g.Expect(ruleMap["all"].Prerequisites).To(Equal([]string{"build", "test"}))
	g.Expect(ruleMap["build"].Prerequisites).To(BeEmpty())
	g.Expect(ruleMap["test"].Prerequisites).To(Equal([]string{"build"}))
	g.Expect(ruleMap["clean"].Prerequisites).To(BeEmpty())
}

func TestParse_BasicRules_CollectsRecipes(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	// Act
	mf, err := Parse(filepath.Join("testdata", "basic-rules.mk"))

	// Assert
	g.Expect(err).NotTo(HaveOccurred())

	ruleMap := makeRuleMap(mf)
	g.Expect(ruleMap["build"].Recipes).To(HaveLen(1))
	g.Expect(ruleMap["build"].Recipes[0]).To(ContainSubstring("go build"))
	g.Expect(ruleMap["all"].Recipes).To(BeEmpty())
}

// makeRuleMap creates a map from target name to Rule for easy lookup in tests.
func makeRuleMap(mf *Makefile) map[string]Rule {
	m := make(map[string]Rule, len(mf.Rules))
	for _, r := range mf.Rules {
		m[r.Target] = r
	}

	return m
}
```

- [ ] **Step 3: Create minimal parser stub**

Create `internal/parser/parser.go`:
```go
package parser

import (
	"bufio"
	"os"
	"strings"

	"github.com/rotisserie/eris"
)

// Parse reads a Makefile and returns a parsed representation.
func Parse(filename string) (*Makefile, error) {
	return parseFile(filename, nil)
}

func parseFile(filename string, visited map[string]bool) (*Makefile, error) {
	if visited == nil {
		visited = make(map[string]bool)
	}

	absPath, err := resolveAbsPath(filename)
	if err != nil {
		return nil, eris.Wrapf(err, "resolving path for %s", filename)
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

func parseReader(r *os.File, filePath string, visited map[string]bool) (*Makefile, error) {
	scanner := bufio.NewScanner(r)
	mf := &Makefile{}

	var currentRule *Rule
	var continuation strings.Builder

	for scanner.Scan() {
		line := scanner.Text()

		// Handle continuation lines
		if continuation.Len() > 0 || strings.HasSuffix(line, "\\") {
			if strings.HasSuffix(line, "\\") {
				continuation.WriteString(strings.TrimSuffix(line, "\\"))
				continuation.WriteString(" ")

				continue
			}

			continuation.WriteString(line)
			line = continuation.String()
			continuation.Reset()
		}

		// Skip blank lines
		if strings.TrimSpace(line) == "" {
			currentRule = nil

			continue
		}

		// Skip full-line comments
		if strings.HasPrefix(strings.TrimSpace(line), "#") {
			continue
		}

		// Recipe lines (tab-indented)
		if line[0] == '\t' && currentRule != nil {
			currentRule.Recipes = append(currentRule.Recipes, strings.TrimPrefix(line, "\t"))

			continue
		}

		// Not a recipe line — clear current rule
		if line[0] == '\t' {
			continue
		}

		// Try to parse as a rule
		if rule, ok := parseRuleLine(line); ok {
			// Filter dot-prefixed targets
			if strings.HasPrefix(rule.Target, ".") {
				continue
			}

			mf.Rules = append(mf.Rules, rule...)
			if len(rule) > 0 {
				currentRule = &mf.Rules[len(mf.Rules)-1]
			}

			continue
		}

		// Everything else is skipped (variable assignments, conditionals, etc.)
		currentRule = nil
	}

	if err := scanner.Err(); err != nil {
		return nil, eris.Wrapf(err, "reading %s", filePath)
	}

	// Handle trailing continuation
	if continuation.Len() > 0 {
		line := strings.TrimSpace(continuation.String())
		if rule, ok := parseRuleLine(line); ok {
			for _, r := range rule {
				if !strings.HasPrefix(r.Target, ".") {
					mf.Rules = append(mf.Rules, r)
				}
			}
		}
	}

	return mf, nil
}

func parseRuleLine(line string) ([]Rule, bool) {
	// Skip lines that look like variable assignments
	// Check for = before any : to detect VAR = value, VAR := value, VAR ?= value, VAR += value
	eqIdx := strings.Index(line, "=")
	colonIdx := findRuleColon(line)

	if colonIdx < 0 {
		return nil, false
	}

	// If there's an = before the colon, this is likely a variable assignment
	if eqIdx >= 0 && eqIdx < colonIdx {
		return nil, false
	}

	// Check for := immediately (e.g., "VAR:=value")
	if colonIdx+1 < len(line) && line[colonIdx+1] == '=' {
		return nil, false
	}

	// Check for double-colon
	isDoubleColon := colonIdx+1 < len(line) && line[colonIdx+1] == ':'
	afterColon := colonIdx + 1
	if isDoubleColon {
		afterColon = colonIdx + 2
	}

	targetsPart := strings.TrimSpace(line[:colonIdx])
	restPart := ""
	if afterColon < len(line) {
		restPart = strings.TrimSpace(line[afterColon:])
	}

	// Split targets by whitespace
	targets := strings.Fields(targetsPart)
	if len(targets) == 0 {
		return nil, false
	}

	// Extract description (## comment) and prerequisites
	var description string
	var prereqStr string

	if descIdx := strings.Index(restPart, "##"); descIdx >= 0 {
		description = strings.TrimSpace(restPart[descIdx+2:])
		prereqStr = strings.TrimSpace(restPart[:descIdx])
	} else if commentIdx := strings.Index(restPart, "#"); commentIdx >= 0 {
		// Single # is just a comment, not a description
		prereqStr = strings.TrimSpace(restPart[:commentIdx])
	} else {
		prereqStr = restPart
	}

	// Split prerequisites by whitespace
	var prereqs []string
	if prereqStr != "" {
		prereqs = strings.Fields(prereqStr)
	}

	// Create a rule for each target
	rules := make([]Rule, 0, len(targets))
	for _, target := range targets {
		rules = append(rules, Rule{
			Target:        target,
			Prerequisites: prereqs,
			Description:   description,
		})
	}

	return rules, true
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
				i++ // skip the ( or {
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

func resolveAbsPath(filename string) (string, error) {
	absPath, err := filepath.Abs(filename)
	if err != nil {
		return "", eris.Wrapf(err, "resolving absolute path for %s", filename)
	}

	return absPath, nil
}
```

Note: This stub needs `"path/filepath"` added to its imports.

- [ ] **Step 4: Run tests to verify they pass**

Run: `go test ./internal/parser/... -v`
Expected: All 3 tests pass.

- [ ] **Step 5: Commit**

```bash
git add -A
git commit -m "Add parser with basic rule parsing and tests"
```

### Task 14: Add parser tests for descriptions, comments, dot-prefix filtering

**Files:**
- Modify: `internal/parser/parser_test.go`
- Create: `internal/parser/testdata/descriptions.mk`
- Create: `internal/parser/testdata/dot-prefix.mk`

- [ ] **Step 1: Create test Makefile for descriptions**

Create `internal/parser/testdata/descriptions.mk`:
```makefile
build: deps ## Build the binary
test-unit: ## Run unit tests
	go test ./...
clean: # Not a description
	rm -rf build/
deploy: deps ## Deploy to production # inline note
```

- [ ] **Step 2: Write description extraction tests**

Add to `parser_test.go`:
```go
func TestParse_Descriptions_ExtractsDoubleHashComments(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	mf, err := Parse(filepath.Join("testdata", "descriptions.mk"))

	g.Expect(err).NotTo(HaveOccurred())

	ruleMap := makeRuleMap(mf)
	g.Expect(ruleMap["build"].Description).To(Equal("Build the binary"))
	g.Expect(ruleMap["test-unit"].Description).To(Equal("Run unit tests"))
	g.Expect(ruleMap["clean"].Description).To(BeEmpty())
	g.Expect(ruleMap["deploy"].Description).To(Equal("Deploy to production # inline note"))
}

func TestParse_InlineComment_DoesNotIncludeCommentInPrerequisites(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	mf, err := Parse(filepath.Join("testdata", "descriptions.mk"))

	g.Expect(err).NotTo(HaveOccurred())

	ruleMap := makeRuleMap(mf)
	g.Expect(ruleMap["clean"].Prerequisites).To(BeEmpty())
}
```

- [ ] **Step 3: Create test Makefile for dot-prefix filtering**

Create `internal/parser/testdata/dot-prefix.mk`:
```makefile
.PHONY: build test clean

build: main.o
	gcc -o build main.o

test: build
	./run-tests.sh

clean:
	rm -rf build/

.SUFFIXES: .c .o

.DEFAULT:
	echo "No rule for $@"
```

- [ ] **Step 4: Write dot-prefix filtering tests**

Add to `parser_test.go`:
```go
func TestParse_DotPrefix_ExcludesSpecialTargets(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	mf, err := Parse(filepath.Join("testdata", "dot-prefix.mk"))

	g.Expect(err).NotTo(HaveOccurred())

	targets := make([]string, len(mf.Rules))
	for i, r := range mf.Rules {
		targets[i] = r.Target
	}

	g.Expect(targets).To(ConsistOf("build", "test", "clean"))
	g.Expect(targets).NotTo(ContainElement(".PHONY"))
	g.Expect(targets).NotTo(ContainElement(".SUFFIXES"))
	g.Expect(targets).NotTo(ContainElement(".DEFAULT"))
}
```

- [ ] **Step 5: Run tests**

Run: `go test ./internal/parser/... -v`
Expected: All tests pass.

- [ ] **Step 6: Commit**

```bash
git add -A
git commit -m "Add parser tests for descriptions and dot-prefix filtering"
```

### Task 15: Add parser tests for continuation lines, multiple targets, double-colon rules

**Files:**
- Modify: `internal/parser/parser_test.go`
- Create: `internal/parser/testdata/continuation.mk`
- Create: `internal/parser/testdata/multiple-targets.mk`
- Create: `internal/parser/testdata/double-colon.mk`

- [ ] **Step 1: Create test Makefile for continuation lines**

Create `internal/parser/testdata/continuation.mk`:
```makefile
build: dep1 \
       dep2 \
       dep3
	gcc -o build dep1.o dep2.o dep3.o
```

- [ ] **Step 2: Write continuation line tests**

```go
func TestParse_ContinuationLines_JoinsLines(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	mf, err := Parse(filepath.Join("testdata", "continuation.mk"))

	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(mf.Rules).To(HaveLen(1))
	g.Expect(mf.Rules[0].Target).To(Equal("build"))
	g.Expect(mf.Rules[0].Prerequisites).To(Equal([]string{"dep1", "dep2", "dep3"}))
}
```

- [ ] **Step 3: Create test Makefile for multiple targets**

Create `internal/parser/testdata/multiple-targets.mk`:
```makefile
all clean distclean: setup
	echo "running $@"

setup:
	echo "setup"
```

- [ ] **Step 4: Write multiple target tests**

```go
func TestParse_MultipleTargets_CreatesSeparateRules(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	mf, err := Parse(filepath.Join("testdata", "multiple-targets.mk"))

	g.Expect(err).NotTo(HaveOccurred())

	targets := make([]string, len(mf.Rules))
	for i, r := range mf.Rules {
		targets[i] = r.Target
	}

	g.Expect(targets).To(ConsistOf("all", "clean", "distclean", "setup"))

	ruleMap := makeRuleMap(mf)
	g.Expect(ruleMap["all"].Prerequisites).To(Equal([]string{"setup"}))
	g.Expect(ruleMap["clean"].Prerequisites).To(Equal([]string{"setup"}))
	g.Expect(ruleMap["distclean"].Prerequisites).To(Equal([]string{"setup"}))
}
```

- [ ] **Step 5: Create test Makefile for double-colon rules**

Create `internal/parser/testdata/double-colon.mk`:
```makefile
build:: compile
	echo "compiling"

build:: link
	echo "linking"
```

- [ ] **Step 6: Write double-colon rule tests**

```go
func TestParse_DoubleColon_MergesRules(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	mf, err := Parse(filepath.Join("testdata", "double-colon.mk"))

	g.Expect(err).NotTo(HaveOccurred())

	ruleMap := makeRuleMap(mf)
	g.Expect(ruleMap["build"].Prerequisites).To(ConsistOf("compile", "link"))
	g.Expect(ruleMap["build"].Recipes).To(HaveLen(2))
}
```

Note: This test will fail until the parser merges double-colon rules. The parser currently creates separate rules for the same target. We need to add merging logic.

- [ ] **Step 7: Add rule merging to parser**

In `parseReader`, after parsing all lines, merge rules with the same target:

```go
func mergeRules(rules []Rule) []Rule {
	seen := make(map[string]int) // target → index in result
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
```

Call `mf.Rules = mergeRules(mf.Rules)` before returning from `parseReader`. Also fix recipe collection — when a double-colon rule re-opens a target, `currentRule` should point to the correct merged rule. The simplest approach: do a post-processing merge step rather than tracking during parsing. Collect all rules during parsing (allowing duplicates), then merge at the end.

- [ ] **Step 8: Run tests**

Run: `go test ./internal/parser/... -v`
Expected: All tests pass.

- [ ] **Step 9: Commit**

```bash
git add -A
git commit -m "Add continuation lines, multiple targets, double-colon rule merging"
```

### Task 16: Add parser tests for variable assignments, conditionals, edge cases

**Files:**
- Modify: `internal/parser/parser_test.go`
- Create: `internal/parser/testdata/skip-variables.mk`
- Create: `internal/parser/testdata/empty.mk`
- Create: `internal/parser/testdata/variable-prereqs.mk`

- [ ] **Step 1: Create Makefile with variable assignments that should be skipped**

Create `internal/parser/testdata/skip-variables.mk`:
```makefile
CC := gcc
CFLAGS = -Wall -O2
SRCS ?= main.c
OBJS += main.o

ifeq ($(DEBUG),1)
CFLAGS += -g
endif

build: main.o
	$(CC) $(CFLAGS) -o build main.o

clean:
	rm -f build
```

- [ ] **Step 2: Write variable-skipping tests**

```go
func TestParse_VariableAssignments_SkippedCorrectly(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	mf, err := Parse(filepath.Join("testdata", "skip-variables.mk"))

	g.Expect(err).NotTo(HaveOccurred())

	targets := make([]string, len(mf.Rules))
	for i, r := range mf.Rules {
		targets[i] = r.Target
	}

	g.Expect(targets).To(ConsistOf("build", "clean"))
}
```

- [ ] **Step 3: Create empty Makefile**

Create `internal/parser/testdata/empty.mk`:
(empty file)

- [ ] **Step 4: Write empty file test**

```go
func TestParse_EmptyFile_ReturnsEmptyMakefile(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	mf, err := Parse(filepath.Join("testdata", "empty.mk"))

	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(mf.Rules).To(BeEmpty())
}
```

- [ ] **Step 5: Create Makefile with variable prerequisites**

Create `internal/parser/testdata/variable-prereqs.mk`:
```makefile
build: $(OBJS)
	gcc -o build $(OBJS)
```

- [ ] **Step 6: Write variable prerequisite test**

```go
func TestParse_VariablePrerequisites_TreatedAsLiteral(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	mf, err := Parse(filepath.Join("testdata", "variable-prereqs.mk"))

	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(mf.Rules).To(HaveLen(1))
	g.Expect(mf.Rules[0].Prerequisites).To(Equal([]string{"$(OBJS)"}))
}
```

- [ ] **Step 7: Run tests**

Run: `go test ./internal/parser/... -v`
Expected: All tests pass.

- [ ] **Step 8: Commit**

```bash
git add -A
git commit -m "Add tests for variable skipping, empty files, variable prerequisites"
```

### Task 17: Add include directive handling

**Files:**
- Modify: `internal/parser/parser.go`
- Modify: `internal/parser/parser_test.go`
- Create: `internal/parser/testdata/includes/main.mk`
- Create: `internal/parser/testdata/includes/lib.mk`
- Create: `internal/parser/testdata/includes/cycle-a.mk`
- Create: `internal/parser/testdata/includes/cycle-b.mk`

- [ ] **Step 1: Create test Makefiles for includes**

Create `internal/parser/testdata/includes/main.mk`:
```makefile
include lib.mk

all: build lib
	echo "all done"
```

Create `internal/parser/testdata/includes/lib.mk`:
```makefile
lib: lib-util
	echo "building lib"

lib-util:
	echo "building lib-util"
```

- [ ] **Step 2: Write include parsing tests**

```go
func TestParse_Include_ParsesIncludedFile(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	mf, err := Parse(filepath.Join("testdata", "includes", "main.mk"))

	g.Expect(err).NotTo(HaveOccurred())

	targets := make([]string, len(mf.Rules))
	for i, r := range mf.Rules {
		targets[i] = r.Target
	}

	g.Expect(targets).To(ConsistOf("all", "lib", "lib-util"))
}
```

- [ ] **Step 3: Add include parsing logic to parser**

In `parseReader`, after checking for comments and before checking for rule lines, add:

```go
trimmed := strings.TrimSpace(line)

// Handle include directives
if handleInclude(trimmed, filePath, visited, mf) {
    continue
}
```

Implement `handleInclude`:
```go
func handleInclude(line, currentFile string, visited map[string]bool, mf *Makefile) bool {
	silent := false
	var paths string

	switch {
	case strings.HasPrefix(line, "include "):
		paths = strings.TrimPrefix(line, "include ")
	case strings.HasPrefix(line, "-include "):
		paths = strings.TrimPrefix(line, "-include ")
		silent = true
	case strings.HasPrefix(line, "sinclude "):
		paths = strings.TrimPrefix(line, "sinclude ")
		silent = true
	default:
		return false
	}

	dir := filepath.Dir(currentFile)
	for _, p := range strings.Fields(paths) {
		incPath := p
		if !filepath.IsAbs(incPath) {
			incPath = filepath.Join(dir, incPath)
		}

		included, err := parseFile(incPath, visited)
		if err != nil {
			if silent {
				continue
			}

			// For non-silent includes, we should propagate the error.
			// This requires changing handleInclude to return an error.
			continue
		}

		mf.Rules = append(mf.Rules, included.Rules...)
	}

	return true
}
```

Note: This needs refinement — `handleInclude` should return an error for non-silent includes. Refactor to return `(bool, error)` and handle the error in the caller.

- [ ] **Step 4: Create cycle detection test files**

Create `internal/parser/testdata/includes/cycle-a.mk`:
```makefile
include cycle-b.mk

a-target:
	echo "a"
```

Create `internal/parser/testdata/includes/cycle-b.mk`:
```makefile
include cycle-a.mk

b-target:
	echo "b"
```

- [ ] **Step 5: Write cycle detection test**

```go
func TestParse_IncludeCycle_DoesNotLoop(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	mf, err := Parse(filepath.Join("testdata", "includes", "cycle-a.mk"))

	g.Expect(err).NotTo(HaveOccurred())

	targets := make([]string, len(mf.Rules))
	for i, r := range mf.Rules {
		targets[i] = r.Target
	}

	g.Expect(targets).To(ConsistOf("a-target", "b-target"))
}
```

- [ ] **Step 6: Write missing include test**

```go
func TestParse_SilentIncludeMissing_ContinuesWithoutError(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	// Create a temp file with -include of nonexistent file
	dir := t.TempDir()
	content := "-include nonexistent.mk\n\nbuild:\n\techo build\n"
	err := os.WriteFile(filepath.Join(dir, "Makefile"), []byte(content), 0o644)
	g.Expect(err).NotTo(HaveOccurred())

	mf, parseErr := Parse(filepath.Join(dir, "Makefile"))

	g.Expect(parseErr).NotTo(HaveOccurred())
	g.Expect(mf.Rules).To(HaveLen(1))
	g.Expect(mf.Rules[0].Target).To(Equal("build"))
}
```

- [ ] **Step 7: Run tests**

Run: `go test ./internal/parser/... -v`
Expected: All tests pass.

- [ ] **Step 8: Commit**

```bash
git add -A
git commit -m "Add include directive handling with cycle detection"
```

### Task 18: Add $(MAKE) detection in recipes

**Files:**
- Create: `internal/parser/recipe.go`
- Create: `internal/parser/recipe_test.go`
- Create: `internal/parser/testdata/recursive-make.mk`

- [ ] **Step 1: Create recipe.go with $(MAKE) detection**

```go
package parser

import (
	"regexp"
	"strings"
)

// makeCallPattern matches common $(MAKE), ${MAKE}, and bare make invocations.
// Captures the target name.
var makeCallPattern = regexp.MustCompile(
	`(?:` +
		`\$\(MAKE\)|\$\{MAKE\}|make` +
		`)` +
		`(?:\s+(?:-[sSkkinBwWjCl]\b|--[a-z-]+(?:=\S+)?))*` + // optional flags
		`(?:\s+-C\s+\S+)?` + // optional -C dir
		`(?:\s+(?:-[sSkkinBwWjl]\b|--[a-z-]+(?:=\S+)?))*` + // more optional flags after -C
		`\s+([a-zA-Z0-9_][a-zA-Z0-9_./-]*)`, // target name (must not start with -)
)

// DetectMakeCall scans a recipe line for a recursive make invocation
// and returns the target name if found.
func DetectMakeCall(recipeLine string) (string, bool) {
	// Strip leading recipe prefixes (@, -, +)
	line := strings.TrimLeft(recipeLine, "@-+ \t")

	matches := makeCallPattern.FindStringSubmatch(line)
	if len(matches) < 2 {
		return "", false
	}

	target := matches[1]

	// Skip if target looks like a variable
	if strings.HasPrefix(target, "$") {
		return "", false
	}

	return target, true
}
```

- [ ] **Step 2: Write recipe detection tests**

Create `internal/parser/recipe_test.go`:
```go
package parser

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestDetectMakeCall_VariousPatterns_DetectsExpectedTargets(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		recipe   string
		expected string
		found    bool
	}{
		"$(MAKE) target": {
			recipe: "$(MAKE) test", expected: "test", found: true,
		},
		"${MAKE} target": {
			recipe: "${MAKE} build", expected: "build", found: true,
		},
		"make target": {
			recipe: "make clean", expected: "clean", found: true,
		},
		"$(MAKE) -C dir target": {
			recipe: "$(MAKE) -C subdir all", expected: "all", found: true,
		},
		"$(MAKE) with flags": {
			recipe: "$(MAKE) -s -k test", expected: "test", found: true,
		},
		"@$(MAKE) target": {
			recipe: "@$(MAKE) build", expected: "build", found: true,
		},
		"+$(MAKE) target": {
			recipe: "+$(MAKE) build", expected: "build", found: true,
		},
		"variable target skipped": {
			recipe: "$(MAKE) $(TARGET)", expected: "", found: false,
		},
		"no make invocation": {
			recipe: "gcc -o main main.c", expected: "", found: false,
		},
		"echo with make word": {
			recipe: "echo \"running make\"", expected: "", found: false,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			g := NewWithT(t)

			target, found := DetectMakeCall(c.recipe)

			g.Expect(found).To(Equal(c.found))
			if c.found {
				g.Expect(target).To(Equal(c.expected))
			}
		})
	}
}
```

- [ ] **Step 3: Run recipe tests**

Run: `go test ./internal/parser/... -run TestDetectMakeCall -v`
Expected: All tests pass.

- [ ] **Step 4: Commit**

```bash
git add -A
git commit -m "Add $(MAKE) detection in recipe lines"
```

---

## Chunk 3: Graph Builder & CLI

### Task 19: Create makegraph builder

**Files:**
- Create: `internal/makegraph/makegraph.go`
- Create: `internal/makegraph/makegraph_test.go`

- [ ] **Step 1: Write failing tests for the graph builder**

```go
package makegraph

import (
	"slices"
	"testing"

	. "github.com/onsi/gomega"

	"github.com/theunrepentantgeek/make-graph/internal/parser"
)

func TestBuild_BasicRules_CreatesNodesForEachTarget(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	mf := &parser.Makefile{
		Rules: []parser.Rule{
			{Target: "all", Prerequisites: []string{"build", "test"}},
			{Target: "build"},
			{Target: "test", Prerequisites: []string{"build"}},
		},
	}

	gr := New(mf).Build()

	nodeIDs := collectNodeIDs(gr)
	g.Expect(nodeIDs).To(ConsistOf("all", "build", "test"))
}

func TestBuild_Prerequisites_CreatesDependencyEdges(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	mf := &parser.Makefile{
		Rules: []parser.Rule{
			{Target: "all", Prerequisites: []string{"build", "test"}},
			{Target: "build"},
			{Target: "test"},
		},
	}

	gr := New(mf).Build()

	allNode, ok := gr.Node("all")
	g.Expect(ok).To(BeTrue())
	g.Expect(allNode.Edges()).To(HaveLen(2))

	for _, edge := range allNode.Edges() {
		g.Expect(edge.Class()).To(Equal("dep"))
	}
}

func TestBuild_MissingPrerequisite_CreatesNodeForIt(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	mf := &parser.Makefile{
		Rules: []parser.Rule{
			{Target: "build", Prerequisites: []string{"main.o"}},
		},
	}

	gr := New(mf).Build()

	_, ok := gr.Node("main.o")
	g.Expect(ok).To(BeTrue())
}

func TestBuild_RecipeWithMakeCall_CreatesCallEdge(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	mf := &parser.Makefile{
		Rules: []parser.Rule{
			{Target: "all", Recipes: []string{"$(MAKE) build"}},
			{Target: "build"},
		},
	}

	gr := New(mf).Build()

	allNode, ok := gr.Node("all")
	g.Expect(ok).To(BeTrue())
	g.Expect(allNode.Edges()).To(HaveLen(1))
	g.Expect(allNode.Edges()[0].Class()).To(Equal("call"))
	g.Expect(allNode.Edges()[0].To().ID()).To(Equal("build"))
}

func TestBuild_Description_SetsNodeDescription(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	mf := &parser.Makefile{
		Rules: []parser.Rule{
			{Target: "build", Description: "Build the binary"},
		},
	}

	gr := New(mf).Build()

	buildNode, ok := gr.Node("build")
	g.Expect(ok).To(BeTrue())
	g.Expect(buildNode.Description).To(Equal("Build the binary"))
}

func collectNodeIDs(gr *graph.Graph) []string {
	var ids []string
	for node := range gr.Nodes() {
		ids = append(ids, node.ID())
	}

	slices.Sort(ids)

	return ids
}
```

Note: Need to import `graph` package. Add: `"github.com/theunrepentantgeek/make-graph/internal/graph"`

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test ./internal/makegraph/... -v`
Expected: Fails (package not found / functions not defined).

- [ ] **Step 3: Implement makegraph builder**

Create `internal/makegraph/makegraph.go`:
```go
package makegraph

import (
	"slices"

	"github.com/theunrepentantgeek/make-graph/internal/graph"
	"github.com/theunrepentantgeek/make-graph/internal/parser"
)

// Builder constructs a graph.Graph from a parsed Makefile.
type Builder struct {
	makefile *parser.Makefile
}

// New creates a new Builder for the given parsed Makefile.
func New(mf *parser.Makefile) *Builder {
	return &Builder{makefile: mf}
}

// Build constructs the dependency graph.
func (b *Builder) Build() *graph.Graph {
	g := graph.New()

	// Sort rules alphabetically for deterministic output
	rules := make([]parser.Rule, len(b.makefile.Rules))
	copy(rules, b.makefile.Rules)
	slices.SortFunc(rules, func(a, b parser.Rule) int {
		if a.Target < b.Target {
			return -1
		}

		if a.Target > b.Target {
			return 1
		}

		return 0
	})

	// Create nodes
	for _, rule := range rules {
		node := g.AddNode(rule.Target)
		node.Description = rule.Description
	}

	// Create edges
	for _, rule := range rules {
		b.addEdgesForPrerequisites(rule, g)
		b.addEdgesForCalls(rule, g)
	}

	return g
}

func (b *Builder) addEdgesForPrerequisites(rule parser.Rule, g *graph.Graph) {
	fromNode, ok := g.Node(rule.Target)
	if !ok {
		return
	}

	for _, prereq := range rule.Prerequisites {
		toNode, exists := g.Node(prereq)
		if !exists {
			toNode = g.AddNode(prereq)
		}

		edge := fromNode.AddEdge(toNode)
		edge.SetClass("dep")
	}
}

func (b *Builder) addEdgesForCalls(rule parser.Rule, g *graph.Graph) {
	fromNode, ok := g.Node(rule.Target)
	if !ok {
		return
	}

	for _, recipe := range rule.Recipes {
		target, found := parser.DetectMakeCall(recipe)
		if !found {
			continue
		}

		toNode, exists := g.Node(target)
		if !exists {
			toNode = g.AddNode(target)
		}

		edge := fromNode.AddEdge(toNode)
		edge.SetClass("call")
	}
}
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `go test ./internal/makegraph/... -v`
Expected: All tests pass.

- [ ] **Step 5: Commit**

```bash
git add -A
git commit -m "Add makegraph builder with tests"
```

### Task 20: Create CLI (cmd package)

**Files:**
- Create: `internal/cmd/cli.go`
- Create: `internal/cmd/context.go`

- [ ] **Step 1: Create context.go**

```go
package cmd

import (
	"log/slog"

	"github.com/theunrepentantgeek/make-graph/internal/config"
)

// Flags contains shared state passed through kong commands.
type Flags struct {
	Verbose bool
	Log     *slog.Logger
	Config  *config.Config
}
```

- [ ] **Step 2: Create cli.go**

Copy the structure from task-graph's `cli.go`, adapting:
- Change `Taskfile` field to `Makefile` with default `"Makefile"`
- Replace `loader.Load()` with `parser.Parse()`
- Replace `taskgraph.New()` with `makegraph.New()`
- Update all imports to use `make-graph` module path
- Replace "TaskNodes" references with appropriate config references
- Keep all config loading, export, auto-color, highlight, and render-image logic

The file should follow the same structure: CLI struct → Run → CreateLogger → CreateConfig → ExportConfigToFile → private helpers.

- [ ] **Step 3: Verify compilation**

Run: `go build ./...`
Expected: Compiles successfully.

- [ ] **Step 4: Commit**

```bash
git add -A
git commit -m "Add CLI command with kong integration"
```

### Task 21: Add CLI tests

**Files:**
- Create: `internal/cmd/cli_test.go`
- Create: `internal/cmd/testdata/config.yaml`
- Create: `internal/cmd/testdata/simple.mk`

- [ ] **Step 1: Create test config file**

Create `internal/cmd/testdata/config.yaml`:
```yaml
graphviz:
  font: "Fira Code"
  fontSize: 12
  callEdges:
    color: "red"
  taskNodes:
    color: "black"
```

- [ ] **Step 2: Create test Makefile**

Create `internal/cmd/testdata/simple.mk`:
```makefile
all: build test ## Build and test everything

build: ## Build the binary
	go build ./...

test: build ## Run tests
	go test ./...

clean: ## Clean build artifacts
	rm -rf build/
```

- [ ] **Step 3: Write CLI tests**

```go
package cmd

import (
	"os"
	"path/filepath"
	"testing"

	. "github.com/onsi/gomega"

	"github.com/theunrepentantgeek/make-graph/internal/config"
)

func TestCreateConfig_DefaultsWhenNoFileProvided(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	cli := CLI{}
	cfg, err := cli.CreateConfig()

	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(cfg).To(Equal(config.New()))
}

func TestCreateConfig_LoadsYAMLConfig(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	cli := CLI{
		Config: filepath.Join("testdata", "config.yaml"),
	}

	cfg, err := cli.CreateConfig()

	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(cfg.Graphviz.Font).To(Equal("Fira Code"))
	g.Expect(cfg.Graphviz.FontSize).To(Equal(12))
}

func TestRun_SimpleMakefile_ProducesDotOutput(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	outputDir := t.TempDir()
	outputFile := filepath.Join(outputDir, "output.dot")

	cli := CLI{
		Makefile: filepath.Join("testdata", "simple.mk"),
		Output:   outputFile,
	}

	cfg := config.New()
	flags := &Flags{
		Config: cfg,
	}

	err := cli.Run(flags)

	g.Expect(err).NotTo(HaveOccurred())

	content, readErr := os.ReadFile(outputFile)
	g.Expect(readErr).NotTo(HaveOccurred())
	g.Expect(string(content)).To(ContainSubstring("digraph"))
	g.Expect(string(content)).To(ContainSubstring("build"))
	g.Expect(string(content)).To(ContainSubstring("test"))
}
```

- [ ] **Step 4: Run CLI tests**

Run: `go test ./internal/cmd/... -v`
Expected: All tests pass.

- [ ] **Step 5: Commit**

```bash
git add -A
git commit -m "Add CLI tests"
```

### Task 22: End-to-end build and manual test

- [ ] **Step 1: Build the binary**

Run:
```bash
cd /home/bevan/github/make-graph
go build -o build/make-graph
```
Expected: Binary created at `build/make-graph`.

- [ ] **Step 2: Create a sample Makefile for manual testing**

Create `samples/simple/Makefile`:
```makefile
.PHONY: all build test clean deploy

all: build test ## Build and test everything

build: deps ## Build the binary
	go build -o bin/app ./cmd/app

deps: ## Install dependencies
	go mod download

test: build ## Run tests
	go test ./...

test-integration: test ## Run integration tests
	go test -tags=integration ./...

clean: ## Clean build artifacts
	rm -rf bin/

deploy: build test ## Deploy to production
	@$(MAKE) -C deploy all

lint: ## Run linter
	golangci-lint run
```

- [ ] **Step 3: Run make-graph on the sample**

Run:
```bash
./build/make-graph samples/simple/Makefile -o samples/simple/Makefile.dot
```
Expected: Generates `samples/simple/Makefile.dot` with a valid Graphviz digraph.

- [ ] **Step 4: Check output content**

Run: `cat samples/simple/Makefile.dot`
Expected: Valid `.dot` file with nodes for all, build, deps, test, test-integration, clean, deploy, lint and edges for prerequisites and the $(MAKE) call edge.

- [ ] **Step 5: Test Mermaid output**

Run:
```bash
./build/make-graph samples/simple/Makefile -o samples/simple/Makefile.mermaid --graph-type mermaid
```
Expected: Valid Mermaid flowchart output.

- [ ] **Step 6: Commit samples**

```bash
git add -A
git commit -m "Add sample Makefile and generated outputs"
```

---

## Chunk 4: Build System, Linting & Cleanup

### Task 23: Create Taskfile.yml

**Files:**
- Create: `Taskfile.yml`

- [ ] **Step 1: Create Taskfile.yml**

```yaml
version: '3'

tasks:
  default:
    deps: [build, unit-test]

  build:
    aliases: [b]
    cmds:
      - go build -o build/make-graph
    sources:
      - '**/*.go'
      - go.mod
      - go.sum
    generates:
      - build/make-graph

  unit-test:
    aliases: [t]
    cmds:
      - go test ./...

  lint:
    aliases: [l]
    cmds:
      - golangci-lint-custom run --verbose

  tidy:
    cmds:
      - task: tidy:gofumpt
      - task: tidy:mod
      - task: tidy:lint

  tidy:gofumpt:
    cmds:
      - gofumpt -w .

  tidy:mod:
    cmds:
      - go mod tidy

  tidy:lint:
    cmds:
      - golangci-lint-custom run --fix

  ci:
    cmds:
      - task: build
      - task: unit-test
      - task: lint

  update-golden-files:
    cmds:
      - go test ./... -update
```

- [ ] **Step 2: Commit**

```bash
git add Taskfile.yml
git commit -m "Add Taskfile.yml for build automation"
```

### Task 24: Copy linter configuration

**Files:**
- Create: `.golangci.yml`

- [ ] **Step 1: Copy .golangci.yml from task-graph**

Copy `/home/bevan/github/task-graph/.golangci.yml` to `/home/bevan/github/make-graph/.golangci.yml`.

- [ ] **Step 2: Commit**

```bash
git add .golangci.yml
git commit -m "Add golangci-lint configuration (copied from task-graph)"
```

### Task 25: Run all tests and verify

- [ ] **Step 1: Run go mod tidy**

Run: `go mod tidy`

- [ ] **Step 2: Run all tests**

Run: `go test ./...`
Expected: All tests pass.

- [ ] **Step 3: Run gofumpt**

Run: `gofumpt -w .`

- [ ] **Step 4: Commit any formatting changes**

```bash
git add -A
git commit -m "Apply gofumpt formatting" --allow-empty
```

### Task 26: Final verification

- [ ] **Step 1: Build the binary**

Run: `go build -o build/make-graph`
Expected: Compiles.

- [ ] **Step 2: Run all tests**

Run: `go test ./... -count=1`
Expected: All tests pass.

- [ ] **Step 3: Run the binary on a sample**

Run: `./build/make-graph samples/simple/Makefile -o /tmp/test.dot`
Expected: Valid output generated.

- [ ] **Step 4: Final commit**

```bash
git add -A
git commit -m "make-graph v0.1.0: initial working version"
```
