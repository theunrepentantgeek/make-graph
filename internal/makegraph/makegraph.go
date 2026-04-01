package makegraph

import (
	"cmp"
	"slices"

	"github.com/theunrepentantgeek/make-graph/internal/graph"
	"github.com/theunrepentantgeek/make-graph/internal/parser"
)

type Builder struct {
	makefile *parser.Makefile
}

func New(mf *parser.Makefile) *Builder {
	return &Builder{makefile: mf}
}

func (b *Builder) Build() *graph.Graph {
	g := graph.New()

	rules := make([]parser.Rule, len(b.makefile.Rules))
	copy(rules, b.makefile.Rules)
	slices.SortFunc(rules, func(a, b parser.Rule) int {
		return cmp.Compare(a.Target, b.Target)
	})

	for _, rule := range rules {
		node := g.AddNode(rule.Target)
		node.Description = rule.Description
	}

	for _, rule := range rules {
		b.addEdgesForPrerequisites(rule, g)
		b.addEdgesForCalls(rule, g)
	}

	return g
}

func (*Builder) addEdgesForPrerequisites(rule parser.Rule, g *graph.Graph) {
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

func (*Builder) addEdgesForCalls(rule parser.Rule, g *graph.Graph) {
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
