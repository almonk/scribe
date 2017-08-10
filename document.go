package main

import (
	"fmt"
	s "strings"

	"github.com/almonk/css"
)

func isDocumentable(class string) bool {
	var redflags = []string{
		":",
		",",
		"@yank",
		"hk-button-group",
	}

	for _, term := range redflags {
		if s.Contains(class, term) {
			return false
		}
	}

	return true
}

func documentClass(class string, rule map[string]*css.CSSStyleDeclaration) {
	class = s.Trim(class, ".")

	if s.HasPrefix(class, "b--") {
		documentBorderColor(class)
	}

	if s.HasPrefix(class, "hk-button") {
		documentHkButton(class)
	}
}

func sanitizeModuleToHumanName(moduleName string) string {
	r := s.NewReplacer(
		"_", "",
		"-", " ",
		".css", "",
	)

	humanName := r.Replace(moduleName)
	humanName = s.Title(humanName)
	return humanName
}

func documentBorderColor(class string) {
	fmt.Printf("<div class='pa2 ba %s'></div>", class)
}

func documentHkButton(class string) {
	fmt.Printf("<button class='%s'>Lorem</button>", class)
}
