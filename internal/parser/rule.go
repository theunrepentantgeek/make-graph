package parser

// Rule represents a single Makefile target rule.
type Rule struct {
	Target        string
	Prerequisites []string
	Description   string
	Recipes       []string
}
