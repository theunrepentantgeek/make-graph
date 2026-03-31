# make-graph Design Spec

## Overview

**make-graph** is a Go command-line tool that reads a Makefile and generates a Graphviz `.dot` or Mermaid flowchart file representing the dependency graph of targets. It can optionally render the graph as an image using the `dot` executable.

The architecture mirrors [task-graph](https://github.com/theunrepentantgeek/task-graph) exactly. The two projects are fully independent — no shared code, no module dependencies between them. make-graph replaces task-graph's Taskfile loader with a custom Makefile parser, while reusing the same intermediate graph model, output generators, and supporting packages.

## Data Flow

```
Makefile → parser.Parse() → *parser.Makefile
    ↓
makegraph.Builder.Build() → *graph.Graph
    ↓
applyAutoColor() [if enabled]
    ↓
graphviz.SaveTo() or mermaid.SaveTo()
    ↓
dot.RenderImage() [optional]
```

## Project Layout

```
make-graph/
  main.go                          # Entry point, kong CLI wiring
  go.mod                           # Module: github.com/theunrepentantgeek/make-graph
  go.sum
  Taskfile.yml                     # Build/test/lint tasks
  DEVELOPMENT.md                   # Coding conventions (copied from task-graph)
  .golangci.yml                    # Linter config (copied from task-graph)
  internal/
    cmd/
      cli.go                       # CLI struct (kong), Run(), config loading
      cli_test.go
      context.go                   # Flags struct
      testdata/                    # Sample Makefiles for CLI tests
    parser/
      parser.go                    # Makefile parser: reads file, returns Makefile struct
      parser_test.go
      makefile.go                  # Parsed Makefile model (targets, includes)
      rule.go                      # Rule struct: target, prerequisites, description
      recipe.go                    # Recipe line scanning ($(MAKE) detection)
      recipe_test.go
      testdata/                    # Sample Makefiles for parser tests
    makegraph/
      makegraph.go                 # Builds graph.Graph from parsed Makefile
      makegraph_test.go
    graph/
      graph.go                     # Graph container (map of nodes)
      graph_test.go
      node.go                      # Node: ID, Label, Description, edges
      node_test.go
      edge.go                      # Edge: from, to, label, class
      edge_test.go
    graphviz/
      graphviz.go                  # .dot file generation from graph model
      graphviz_test.go
      record.go                    # Mrecord node formatting
      record_test.go
      node_properties.go           # Node attribute handling
      node_properties_test.go
      edge_properties.go           # Edge attribute handling
      edge_properties_test.go
      properties.go                # Shared property map
      properties_test.go
      testdata/                    # Golden files
    mermaid/
      mermaid.go                   # Mermaid flowchart generation
      mermaid_test.go
      testdata/                    # Golden files
    config/
      config.go                    # Config struct, loading, validation
      graphviz.go                  # Graphviz-specific config types
      mermaid.go                   # Mermaid-specific config types
      node_style.go                # NodeStyleRule struct
    dot/
      dot.go                       # dot executable discovery and image rendering
      dot_test.go
    indentwriter/
      indent_writer.go             # Hierarchical text builder
      indent_writer_test.go
      line.go                      # Line struct with nesting
    safe/
      registry.go                  # Safe identifier generation
      label.go                     # Label escaping
    namespace/
      namespace.go                 # Namespace extraction and pattern matching
    autocolor/
      autocolor.go                 # Auto-color assignment by namespace
      autocolor_test.go
  samples/                         # Sample Makefiles and generated output
  docs/
  .devcontainer/
    install-dependencies.sh        # Tool installation
  .github/
    workflows/
      pr-validation.yml            # CI pipeline
```

## Package Responsibilities

| Package | Responsibility | Source |
|---|---|---|
| `cmd` | CLI struct, config loading, orchestration | Adapted from task-graph |
| `parser` | Makefile parsing: targets, prerequisites, includes, recipes | **New** |
| `makegraph` | Build `graph.Graph` from `parser.Makefile` | **New** (mirrors `taskgraph`) |
| `graph` | Minimal graph data structure (Node, Edge, Graph) | Copied from task-graph |
| `graphviz` | Generate `.dot` output with Mrecord nodes, namespace clusters | Copied from task-graph |
| `mermaid` | Generate Mermaid flowchart output | Copied from task-graph |
| `config` | Config structs, loading, validation | Copied from task-graph |
| `dot` | Find `dot` executable, render images | Copied from task-graph |
| `namespace` | Extract/match namespaces, glob patterns | Copied from task-graph, tweaked |
| `autocolor` | Auto-assign colors by namespace | Copied from task-graph |
| `indentwriter` | Build hierarchical text structures | Copied from task-graph |
| `safe` | Sanitize identifiers, avoid collisions | Copied from task-graph |

## Makefile Parser

### Parsed Model

```go
// Makefile represents a parsed Makefile.
type Makefile struct {
    Rules    []Rule
    Includes []string
}

// Rule represents a single Makefile target rule.
type Rule struct {
    Target        string
    Prerequisites []string
    Description   string   // Extracted from ## comment
    Recipes       []string // Raw recipe lines (for $(MAKE) scanning)
}
```

### Parsing Strategy

Single-pass, line-by-line:

1. **Join continuation lines** — lines ending with `\` are concatenated with the next line before processing.
2. **Skip blank lines and full-line comments** — lines that are empty or start with `#`.
3. **Handle `include` directives** — `include path/to/file` and `-include path/to/file` (silent variant). Resolve paths relative to the directory of the current Makefile. Recursively parse included files.
4. **Detect rule lines** — match `target: prerequisites` pattern. A line containing `:` where the left side has no `=` (to avoid variable assignments like `VAR := value`). Also extract `## description` from the end of the line.
5. **Collect recipe lines** — lines starting with a tab character after a rule line belong to that rule's recipe. Store raw text for `$(MAKE)` scanning.
6. **Skip everything else** — variable assignments (`VAR = ...`, `VAR := ...`, `VAR ?= ...`, `VAR += ...`), conditionals (`ifeq`, `ifdef`, `ifndef`, `else`, `endif`), directives (`.EXPORT_ALL_VARIABLES`, `define`, `endef`), and any other unrecognized lines.

### Dot-Prefix Filtering

Any target starting with `.` is excluded at parse time. This covers all special targets (`.PHONY`, `.SUFFIXES`, `.DEFAULT`, `.PRECIOUS`, `.INTERMEDIATE`, `.SECONDARY`, `.DELETE_ON_ERROR`, `.SILENT`, `.IGNORE`, `.NOTPARALLEL`, etc.) as well as any user-defined hidden targets.

### Include Handling

- Both `include` and `-include` (also `sinclude`) are recognized.
- Paths resolved relative to the directory of the file containing the directive.
- **Cycle detection**: maintain a set of absolute file paths in the current parse chain. If a file is encountered that's already being parsed, skip it and continue.
- Missing files referenced by `include` produce an error; missing files referenced by `-include` are silently skipped.

### $(MAKE) Detection

Recipe lines are scanned with regex for recursive make invocations:

- `$(MAKE) target` / `${MAKE} target`
- `$(MAKE) -C dir target` / `${MAKE} -C dir target`
- Also detect bare `make target` and `make -C dir target`

Extraction is best-effort. If the target cannot be determined (e.g., the target is a variable, the command is too complex), the line is silently skipped. No variable expansion is attempted.

### Description Extraction

The `## description` convention (commonly used for `make help` targets):

```makefile
build: deps ## Build the binary
test-unit: ## Run unit tests
```

A `##` followed by text on a rule line is captured as the rule's `Description`. Single `#` comments are not captured (they're often internal notes, not user-facing descriptions).

## Graph Builder (makegraph)

### Interface

```go
type Builder struct {
    makefile *parser.Makefile
}

func New(mf *parser.Makefile) *Builder
func (b *Builder) Build() *graph.Graph
```

### Build Process

1. **Create nodes** — iterate all rules alphabetically, create a `graph.Node` for each target. Set `Label` to the target name, `Description` from the `## comment`.
2. **Create prerequisite edges** — for each rule, create an edge from the target to each prerequisite with class `"dep"`. If a prerequisite has no corresponding rule, create a node for it anyway (it may be a file target).
3. **Create call edges** — for each `$(MAKE)` invocation detected in recipes, create an edge from the calling target to the invoked target with class `"call"`.
4. **Handle missing targets** — prerequisites or `$(MAKE)` targets not defined as rules still get nodes. No error — they appear as leaf nodes with no description.

### Edge Classes

- `"dep"` — prerequisite relationships. Styled as solid lines (default).
- `"call"` — `$(MAKE)` invocations. Styled as dashed lines (default).

### Determinism

Rules and prerequisites processed in alphabetical order for deterministic output across runs.

## CLI

### Kong Struct

```go
type CLI struct {
    Makefile         string `arg:"" help:"Path to Makefile" default:"Makefile"`
    Output           string `long:"output" short:"o" required:"true" help:"Output file path"`
    Config           string `long:"config" short:"c" help:"Config file (YAML or JSON)"`
    GroupByNamespace bool   `long:"group-by-namespace" help:"Group targets by namespace"`
    AutoColor        bool   `long:"auto-color" help:"Auto-color by namespace"`
    GraphType        string `long:"graph-type" default:"dot" enum:"dot,mermaid" help:"Output format"`
    Highlight        string `long:"highlight" help:"Comma-separated patterns to highlight"`
    RenderImage      string `long:"render-image" help:"Render image file type (png, svg, etc.)"`
    ExportConfig     string `long:"export-config" help:"Export effective config to file"`
    Verbose          bool   `long:"verbose" short:"v"`
}
```

### Run Method

```go
func (c *CLI) Run(flags *Flags) error {
    // 1. Parse Makefile
    mf, err := parser.Parse(c.Makefile)

    // 2. Build graph
    gr := makegraph.New(mf).Build()

    // 3. Apply auto-color if enabled
    if flags.Config.AutoColor {
        autocolor.Apply(gr, flags.Config)
    }

    // 4. Generate output
    switch flags.Config.GraphType {
    case "dot":
        graphviz.SaveTo(c.Output, gr, flags.Config)
    case "mermaid":
        mermaid.SaveTo(c.Output, gr, flags.Config)
    }

    // 5. Optionally render image
    if c.RenderImage != "" {
        dot.RenderImage(ctx, dotExe, c.Output, imageFile, c.RenderImage)
    }
}
```

## Configuration

Identical to task-graph's config system:

```go
type Config struct {
    GroupByNamespace bool
    GraphType        string
    HighlightColor   string
    AutoColor        bool
    NodeStyleRules   []NodeStyleRule
    Graphviz         *Graphviz
    Mermaid          *Mermaid
    DotPath          string
}
```

Same YAML/JSON loading, same CLI-overrides-file-overrides-defaults priority, same glob-pattern style rules.

## Namespace Handling

Copied from task-graph with one adjustment: **colon `:` is removed as a namespace delimiter**. Makefiles use colons in rule syntax, not as namespace separators.

Supported delimiters for informal namespaces:
- Hyphen `-`: `build-docker` → namespace `build`
- Dot `.`: `build.docker` → namespace `build`
- Slash `/`: `build/docker` → namespace `build`

## Output Formats

### Graphviz (.dot)

Copied from task-graph. Nodes rendered as Mrecord shapes with target name and description. Supports namespace clustering via `subgraph cluster_*` blocks. Edge styling by class (dep vs call).

### Mermaid

Copied from task-graph. Flowchart output with configurable direction (TD, LR, BT, RL). Style rules applied as separate `style` directives.

## Testing Strategy

### Libraries

- `gomega` — fluent assertions
- `goldie/v2` — golden file tests
- Test naming: `Test<Subject>_<Scenario>_<Expectation>`
- Table-driven with `t.Parallel()`

### Parser Tests (most critical — novel code)

| Area | Tests |
|---|---|
| Basic rules | Single target, multiple prerequisites, no prerequisites |
| Multiple targets | Multiple targets on one line (`all clean: deps`) |
| Descriptions | `## comment` extraction, no `##`, edge cases |
| Continuation lines | `\` line joining across 2+ lines |
| Includes | `include`, `-include`, relative paths, cycle detection, missing files |
| Recipe lines | Tab-indented capture, non-tab ends recipe |
| Dot-prefix filtering | `.PHONY`, `.SUFFIXES`, etc. excluded |
| $(MAKE) detection | `$(MAKE) target`, `${MAKE} -C dir target`, bare `make target`, unrecognized patterns |
| Comments | Full-line `#` skipped, inline `#` vs `##` distinction |
| Edge cases | Empty file, only comments, target with no recipe, blank lines between rules |
| Variable assignments | `VAR = ...`, `VAR := ...`, `VAR ?= ...` all skipped |
| Conditionals | `ifeq`, `ifdef`, `else`, `endif` all skipped |

### Graph Builder Tests

- Node creation for each target
- Prerequisite edges (class `"dep"`)
- Call edges (class `"call"`)
- Missing prerequisite targets get nodes
- Alphabetical ordering

### Output Tests

Copied from task-graph's graphviz and mermaid test suites. These operate on `graph.Graph` and should pass with minimal changes.

### CLI Integration Tests

Invoke `CLI.Run()` with sample Makefiles, verify output files match golden expectations.

### Golden File Updates

```bash
go test ./... -update
# or
task update-golden-files
```

## Build & CI

### Taskfile.yml Tasks

| Task | Command |
|---|---|
| `build` | `go build -o build/make-graph` |
| `unit-test` | `go test ./...` |
| `lint` | Custom `golangci-lint-custom` |
| `tidy` | `gofumpt` + `go mod tidy` + lint fix |
| `ci` | build + test + lint |
| `update-golden-files` | `go test ./... -update` |

### Dependencies

| Dependency | Purpose |
|---|---|
| `alecthomas/kong` | CLI parsing |
| `onsi/gomega` | Test assertions |
| `sebdah/goldie/v2` | Golden file tests |
| `rotisserie/eris` | Error wrapping |
| `phsym/console-slog` | Structured logging |
| `gopkg.in/yaml.v3` | Config loading |

No `go-task/task/v3` dependency — the Makefile parser is our own code.

### Linting

Same `.golangci.yml` and custom-built `golangci-lint-custom` with nilaway as task-graph. Must use `task lint`, never run `golangci-lint` directly.

### Devcontainer & CI Workflows

Copied from task-graph, references updated from `task-graph` to `make-graph`.

## Coding Conventions

All conventions from task-graph's `DEVELOPMENT.md` apply:

- Small, focused packages with one clear responsibility
- Explicit over implicit — dependencies passed as parameters
- Errors wrapped with `eris.Wrap`/`eris.Wrapf` — never bare `fmt.Errorf`
- Max 120 character lines, max 60 line functions
- `gofumpt` formatting
- Strict linting as safety net
