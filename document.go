package main

import (
	s "strings"
)

func isDocumentable(class string) bool {
	var redflags = []string{}

	for _, term := range redflags {
		if s.Contains(class, term) {
			return false
		}
	}

	return true
}

func documentClass(class string, template string) string {
	template = s.TrimSpace(template)
	class = s.Trim(class, ".")
	outputHTML := s.Replace(template, "{{class}}", class, -1)
	return outputHTML
}

func writeHTMLHeader() string {
	return ("<html><head><link rel=stylesheet href='https://www.herokucdn.com/purple3/latest/purple3.min.css'></head><body>")
}

func writeHTMLFooter() string {
	return ("</body></html>")
}
