package yara

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var ruleNameRegex = regexp.MustCompile(`^[a-zA-Z0-9_]{1,64}$`)
var ruleDeclRegex = regexp.MustCompile(`\brule\s+([a-zA-Z0-9_]+)`)

// ValidateRule checks if ruleName is safe and if content contains valid YARA syntax structure.
func ValidateRule(ruleName, content string) error {
	trimmedName := strings.TrimSpace(ruleName)
	if trimmedName == "" {
		return errors.New("rule name cannot be empty")
	}

	if !ruleNameRegex.MatchString(trimmedName) {
		return fmt.Errorf("invalid rule name '%s': must contain only alphanumeric characters and underscores (max 64 chars)", trimmedName)
	}

	trimmedContent := strings.TrimSpace(content)
	if trimmedContent == "" {
		return errors.New("rule content cannot be empty")
	}

	// Must contain 'rule <name>' declaration
	matches := ruleDeclRegex.FindStringSubmatch(trimmedContent)
	if len(matches) < 2 {
		return errors.New("rule content must contain a valid 'rule <name>' declaration")
	}

	// Must contain condition keyword
	if !strings.Contains(trimmedContent, "condition:") {
		return errors.New("rule content must contain a 'condition:' section")
	}

	// Check balanced braces outside string literals and comments
	openBraces, closeBraces := countBraces(trimmedContent)
	if openBraces == 0 || openBraces != closeBraces {
		return errors.New("unbalanced braces '{' and '}' in YARA rule")
	}

	return nil
}

func countBraces(s string) (int, int) {
	openCount := 0
	closeCount := 0

	inString := false
	inLineComment := false
	inBlockComment := false
	escaped := false

	n := len(s)
	for i := 0; i < n; i++ {
		c := s[i]

		if inLineComment {
			if c == '\n' {
				inLineComment = false
			}
			continue
		}

		if inBlockComment {
			if c == '*' && i+1 < n && s[i+1] == '/' {
				inBlockComment = false
				i++ // skip '/'
			}
			continue
		}

		if inString {
			if escaped {
				escaped = false
			} else if c == '\\' {
				escaped = true
			} else if c == '"' {
				inString = false
			}
			continue
		}

		// Not in string or comment
		if c == '/' && i+1 < n {
			if s[i+1] == '/' {
				inLineComment = true
				i++
				continue
			} else if s[i+1] == '*' {
				inBlockComment = true
				i++
				continue
			}
		}

		if c == '"' {
			inString = true
			escaped = false
			continue
		}

		if c == '{' {
			openCount++
		} else if c == '}' {
			closeCount++
		}
	}

	return openCount, closeCount
}
