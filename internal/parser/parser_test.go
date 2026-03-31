package parser

import (
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

// makeRuleMap creates a map from target name to Rule for convenient test lookups.
func makeRuleMap(rules []Rule) map[string]Rule {
	m := make(map[string]Rule, len(rules))
	for _, r := range rules {
		m[r.Target] = r
	}

	return m
}
