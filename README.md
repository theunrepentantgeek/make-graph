# make-graph

A command-line tool that reads a `Makefile` and generates a dependency graph in
[Graphviz](https://graphviz.org/) `.dot` or [Mermaid](https://mermaid.js.org/)
format. Optionally renders the graph directly to an image (PNG, SVG, etc.) using
the `dot` executable.

## Quick start

```bash
# Install
go install github.com/theunrepentantgeek/make-graph@latest

# Generate a Graphviz dot file
make-graph Makefile -o Makefile.dot

# Generate a Mermaid flowchart instead
make-graph Makefile -o Makefile.mermaid --graph-type mermaid

# Generate dot and render to PNG in one step (requires graphviz installed)
make-graph Makefile -o Makefile.dot --render-image png
```

## Example

Given a simple `Makefile`:

```makefile
all: build test

build: deps
	go build -o bin/app ./cmd/app

deps:
	go mod download

test: build
	go test ./...

deploy: build test
	@$(MAKE) -C deploy all
```

Running `make-graph Makefile -o Makefile.dot` produces a Graphviz digraph showing
the dependency relationships between targets:

```dot
digraph {
  "all" -> "build"
  "all" -> "test"
  "build" -> "deps"
  "test" -> "build"
  "deploy" -> "build"
  "deploy" -> "test"
  "deploy" -.-> "all"   // $(MAKE) call shown as dashed edge
  ...
}
```

Solid edges represent prerequisite dependencies; dashed edges represent recursive
`$(MAKE)` calls.

## Usage

```
Usage: make-graph --output=STRING [<makefile>] [flags]

Arguments:
  [<makefile>]    Path to the Makefile to process.

Flags:
  -h, --help                    Show context-sensitive help.
  -o, --output=STRING           Path to the output file.
  -c, --config=STRING           Path to a config file (YAML or JSON).
      --group-by-namespace      Group targets in the same namespace together
                                in the output.
      --auto-color              Automatically color nodes by namespace using
                                a built-in palette.
      --graph-type=STRING       Type of graph to generate (dot or mermaid).
                                Defaults to dot.
      --highlight=STRING        Highlight specific targets in the graph.
                                Accepts target names or glob patterns,
                                separated by commas or semicolons.
      --render-image=STRING     Render the graph as an image using graphviz
                                dot. Specify the file type (e.g. png, svg).
      --export-config=STRING    Export the effective configuration to a file
                                (YAML or JSON based on file extension).
      --verbose                 Enable verbose logging.
```

## Installation

### From source

Requires Go 1.24 or later.

```bash
go install github.com/theunrepentantgeek/make-graph@latest
```

### Build locally

```bash
git clone https://github.com/theunrepentantgeek/make-graph.git
cd make-graph
go build -o make-graph
```

## Configuration

Styling and behaviour can be customised via a YAML or JSON config file passed
with `--config`:

```yaml
# Graph type: "dot" or "mermaid"
graphType: dot

# Group targets sharing a namespace prefix (before ':') together
groupByNamespace: true

# Automatically assign colors to namespaces
autoColor: true

# Fill color used by --highlight (default: yellow)
highlightColor: gold

# Path to the dot executable (if not on PATH)
dotPath: /usr/local/bin/dot

# Graphviz-specific settings
graphviz:
  font: Verdana
  fontSize: 16
  taskNodes:
    color: black
    fillColor: lightyellow
    style: filled
    fontColor: black
  dependencyEdges:
    color: black
    width: 1
    style: solid
  callEdges:
    color: blue
    width: 1
    style: dashed

# Mermaid-specific settings
mermaid:
  direction: TD    # TD, LR, BT, or RL

# Pattern-matched style rules (applied in order; last match wins)
nodeStyleRules:
  - match: "test*"
    fillColor: lightblue
    style: filled
  - match: "deploy"
    fillColor: salmon
    style: filled
```

Export the effective configuration (including defaults) with:

```bash
make-graph Makefile -o graph.dot --export-config effective.yaml
```

## Features

- **Graphviz dot output** — generates a `.dot` digraph with styled nodes and
  edges, suitable for rendering with Graphviz.
- **Mermaid output** — generates a Mermaid flowchart for embedding in Markdown
  or documentation.
- **Image rendering** — renders the dot graph directly to PNG, SVG, or other
  formats supported by Graphviz (requires `dot` installed).
- **Namespace grouping** — targets sharing a common prefix (before `:`) can be
  visually grouped together with `--group-by-namespace`.
- **Auto-coloring** — `--auto-color` assigns distinct fill colors to each
  namespace automatically.
- **Highlighting** — `--highlight` accepts target names or glob patterns to
  visually emphasise specific targets.
- **Style rules** — pattern-matched style rules in the config file allow
  fine-grained control over node appearance.
- **Recursive make detection** — `$(MAKE)` calls are detected and shown as
  dashed call edges, distinct from prerequisite dependencies.

## Contributing

See [DEVELOPMENT.md](DEVELOPMENT.md) for coding conventions, testing practices,
and build instructions.
