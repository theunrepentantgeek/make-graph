package makegraph

import (
	"testing"

	. "github.com/onsi/gomega"

	"github.com/theunrepentantgeek/make-graph/internal/graph"
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
	ids := make([]string, 0, gr.NodeCount())
	for node := range gr.Nodes() {
		ids = append(ids, node.ID())
	}

	return ids
}
