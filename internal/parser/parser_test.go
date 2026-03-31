package parser

import (
	"os"
	"path/filepath"
	"testing"

	. "github.com/onsi/gomega"
)

func TestParse_BasicRules_ReturnsExpectedTargets(t *testing.T) {
	t.Parallel()
	g := NewGomegaWithT(t)

	// Act
	mf, err := Parse("testdata/basic-rules.mk")

	// Assert
	g.Expect(err).ToNot(HaveOccurred())
	g.Expect(mf.Rules).To(HaveLen(4))

	rules := makeRuleMap(mf.Rules)
	g.Expect(rules).To(HaveKey("all"))
	g.Expect(rules).To(HaveKey("build"))
	g.Expect(rules).To(HaveKey("test"))
	g.Expect(rules).To(HaveKey("clean"))
}

func TestParse_BasicRules_ReturnsExpectedPrerequisites(t *testing.T) {
	t.Parallel()
	g := NewGomegaWithT(t)

	// Act
	mf, err := Parse("testdata/basic-rules.mk")

	// Assert
	g.Expect(err).ToNot(HaveOccurred())

	rules := makeRuleMap(mf.Rules)
	g.Expect(rules["all"].Prerequisites).To(ConsistOf("build", "test"))
	g.Expect(rules["build"].Prerequisites).To(BeEmpty())
	g.Expect(rules["test"].Prerequisites).To(ConsistOf("build"))
	g.Expect(rules["clean"].Prerequisites).To(BeEmpty())
}

func TestParse_BasicRules_CollectsRecipes(t *testing.T) {
	t.Parallel()
	g := NewGomegaWithT(t)

	// Act
	mf, err := Parse("testdata/basic-rules.mk")

	// Assert
	g.Expect(err).ToNot(HaveOccurred())

	rules := makeRuleMap(mf.Rules)
	g.Expect(rules["all"].Recipes).To(BeEmpty())
	g.Expect(rules["build"].Recipes).To(ConsistOf("go build ./..."))
	g.Expect(rules["test"].Recipes).To(ConsistOf("go test ./..."))
	g.Expect(rules["clean"].Recipes).To(ConsistOf("rm -rf build/"))
}

func TestParse_Descriptions_ExtractsDoubleHashComments(t *testing.T) {
	t.Parallel()
	g := NewGomegaWithT(t)

	// Act
	mf, err := Parse("testdata/descriptions.mk")

	// Assert
	g.Expect(err).ToNot(HaveOccurred())

	rules := makeRuleMap(mf.Rules)
	g.Expect(rules["build"].Description).To(Equal("Build the binary"))
	g.Expect(rules["test-unit"].Description).To(Equal("Run unit tests"))
	g.Expect(rules["clean"].Description).To(BeEmpty())
	g.Expect(rules["deploy"].Description).To(Equal("Deploy to production # inline note"))
}

func TestParse_InlineComment_DoesNotIncludeCommentInPrerequisites(t *testing.T) {
	t.Parallel()
	g := NewGomegaWithT(t)

	// Act
	mf, err := Parse("testdata/descriptions.mk")

	// Assert
	g.Expect(err).ToNot(HaveOccurred())

	rules := makeRuleMap(mf.Rules)
	g.Expect(rules["clean"].Prerequisites).To(BeEmpty())
}

func TestParse_DotPrefix_ExcludesSpecialTargets(t *testing.T) {
	t.Parallel()
	g := NewGomegaWithT(t)

	// Act
	mf, err := Parse("testdata/dot-prefix.mk")

	// Assert
	g.Expect(err).ToNot(HaveOccurred())
	g.Expect(mf.Rules).To(HaveLen(3))

	rules := makeRuleMap(mf.Rules)
	g.Expect(rules).To(HaveKey("build"))
	g.Expect(rules).To(HaveKey("test"))
	g.Expect(rules).To(HaveKey("clean"))
}

func TestParse_ContinuationLines_JoinsLines(t *testing.T) {
	t.Parallel()
	g := NewGomegaWithT(t)

	// Act
	mf, err := Parse("testdata/continuation.mk")

	// Assert
	g.Expect(err).ToNot(HaveOccurred())
	g.Expect(mf.Rules).To(HaveLen(1))

	rules := makeRuleMap(mf.Rules)
	g.Expect(rules["build"].Prerequisites).To(ConsistOf("dep1", "dep2", "dep3"))
}

func TestParse_MultipleTargets_CreatesSeparateRules(t *testing.T) {
	t.Parallel()
	g := NewGomegaWithT(t)

	// Act
	mf, err := Parse("testdata/multiple-targets.mk")

	// Assert
	g.Expect(err).ToNot(HaveOccurred())
	g.Expect(mf.Rules).To(HaveLen(4))

	rules := makeRuleMap(mf.Rules)
	g.Expect(rules).To(HaveKey("all"))
	g.Expect(rules).To(HaveKey("clean"))
	g.Expect(rules).To(HaveKey("distclean"))
	g.Expect(rules).To(HaveKey("setup"))
	g.Expect(rules["all"].Prerequisites).To(ConsistOf("setup"))
	g.Expect(rules["clean"].Prerequisites).To(ConsistOf("setup"))
	g.Expect(rules["distclean"].Prerequisites).To(ConsistOf("setup"))
}

func TestParse_DoubleColon_MergesRules(t *testing.T) {
	t.Parallel()
	g := NewGomegaWithT(t)

	// Act
	mf, err := Parse("testdata/double-colon.mk")

	// Assert
	g.Expect(err).ToNot(HaveOccurred())
	g.Expect(mf.Rules).To(HaveLen(1))

	rules := makeRuleMap(mf.Rules)
	g.Expect(rules["build"].Prerequisites).To(ConsistOf("compile", "link"))
	g.Expect(rules["build"].Recipes).To(HaveLen(2))
}

func TestParse_VariableAssignments_SkippedCorrectly(t *testing.T) {
	t.Parallel()
	g := NewGomegaWithT(t)

	// Act
	mf, err := Parse("testdata/skip-variables.mk")

	// Assert
	g.Expect(err).ToNot(HaveOccurred())
	g.Expect(mf.Rules).To(HaveLen(2))

	rules := makeRuleMap(mf.Rules)
	g.Expect(rules).To(HaveKey("build"))
	g.Expect(rules).To(HaveKey("clean"))
}

func TestParse_EmptyFile_ReturnsEmptyMakefile(t *testing.T) {
	t.Parallel()
	g := NewGomegaWithT(t)

	// Act
	mf, err := Parse("testdata/empty.mk")

	// Assert
	g.Expect(err).ToNot(HaveOccurred())
	g.Expect(mf.Rules).To(BeEmpty())
}

func TestParse_VariablePrerequisites_TreatedAsLiteral(t *testing.T) {
	t.Parallel()
	g := NewGomegaWithT(t)

	// Act
	mf, err := Parse("testdata/variable-prereqs.mk")

	// Assert
	g.Expect(err).ToNot(HaveOccurred())
	g.Expect(mf.Rules).To(HaveLen(1))

	rules := makeRuleMap(mf.Rules)
	g.Expect(rules["build"].Prerequisites).To(ConsistOf("$(OBJS)"))
}

func TestParse_Include_ParsesIncludedFile(t *testing.T) {
	t.Parallel()
	g := NewGomegaWithT(t)

	// Act
	mf, err := Parse("testdata/includes/main.mk")

	// Assert
	g.Expect(err).ToNot(HaveOccurred())

	rules := makeRuleMap(mf.Rules)
	g.Expect(rules).To(HaveKey("all"))
	g.Expect(rules).To(HaveKey("lib"))
	g.Expect(rules).To(HaveKey("lib-util"))
}

func TestParse_IncludeCycle_DoesNotLoop(t *testing.T) {
	t.Parallel()
	g := NewGomegaWithT(t)

	// Act
	mf, err := Parse("testdata/includes/cycle-a.mk")

	// Assert
	g.Expect(err).ToNot(HaveOccurred())

	rules := makeRuleMap(mf.Rules)
	g.Expect(rules).To(HaveKey("a-target"))
	g.Expect(rules).To(HaveKey("b-target"))
}

func TestParse_SilentIncludeMissing_ContinuesWithoutError(t *testing.T) {
	t.Parallel()
	g := NewGomegaWithT(t)

	// Arrange
	dir := t.TempDir()
	mkFile := filepath.Join(dir, "Makefile")
	content := "-include nonexistent.mk\n\nbuild:\n\techo build\n"
	err := os.WriteFile(mkFile, []byte(content), 0o644)
	g.Expect(err).ToNot(HaveOccurred())

	// Act
	mf, parseErr := Parse(mkFile)

	// Assert
	g.Expect(parseErr).ToNot(HaveOccurred())
	g.Expect(mf.Rules).To(HaveLen(1))

	rules := makeRuleMap(mf.Rules)
	g.Expect(rules).To(HaveKey("build"))
}

// makeRuleMap creates a map from target name to Rule for convenient test lookups.
func makeRuleMap(rules []Rule) map[string]Rule {
	m := make(map[string]Rule, len(rules))
	for _, r := range rules {
		m[r.Target] = r
	}

	return m
}
