package main

import (
	"fmt"
	s "strings"

	"github.com/almonk/css"
)

func isDocumentable(class string) bool {
	var redflags = []string{
		":",
		"@yank",
		"hk-button-group",
		"--",
	}

	for _, term := range redflags {
		if s.Contains(class, term) {
			return false
		}
	}

	return true
}

func documentClass(class string, rule map[string]*css.CSSStyleDeclaration, template string) {
	template = s.TrimSpace(template)
	class = s.Trim(class, ".")
	outputHTML := s.Replace(template, "{{class}}", class, 1)
	fmt.Printf(outputHTML)
}
