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
	g.Expect(cfg.Graphviz.CallEdges.Color).To(Equal("red"))
	g.Expect(cfg.Graphviz.TaskNodes.Color).To(Equal("black"))
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
		Log:    cli.CreateLogger(),
	}

	err := cli.Run(flags)

	g.Expect(err).NotTo(HaveOccurred())

	content, readErr := os.ReadFile(outputFile)
	g.Expect(readErr).NotTo(HaveOccurred())
	g.Expect(string(content)).To(ContainSubstring("digraph"))
	g.Expect(string(content)).To(ContainSubstring("build"))
	g.Expect(string(content)).To(ContainSubstring("test"))
}

func TestRun_SimpleMakefile_ProducesMermaidOutput(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	outputDir := t.TempDir()
	outputFile := filepath.Join(outputDir, "output.mermaid")

	cli := CLI{
		Makefile:  filepath.Join("testdata", "simple.mk"),
		Output:    outputFile,
		GraphType: "mermaid",
	}

	cfg := config.New()
	flags := &Flags{
		Config: cfg,
		Log:    cli.CreateLogger(),
	}

	err := cli.Run(flags)

	g.Expect(err).NotTo(HaveOccurred())

	content, readErr := os.ReadFile(outputFile)
	g.Expect(readErr).NotTo(HaveOccurred())
	g.Expect(string(content)).To(ContainSubstring("flowchart"))
	g.Expect(string(content)).To(ContainSubstring("build"))
}
