package snips

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

var (
	UnclosedSnip  = errors.New("snippet was not closed")
	MalformedSnip = errors.New("snippet was malformed")
)

type Scanner interface {
	Scan() bool
	Text() string
}

type Snip struct {
	Trigger  string
	Hint     string
	Template string
}

func (s Snip) String() string {
	return fmt.Sprintf(
		"Snippet Trigger: %s\nHint: %s\nTemplate:\n%s",
		s.Trigger,
		s.Hint,
		s.Template,
	)
}

func ParseLines(s Scanner) ([]Snip, error) {
	var snippets []Snip
	var snip Snip

	isSnip := false

	for s.Scan() {
		line := s.Text()

		if isSnip {
			if strings.Index(line, "endsnippet") == 0 {
				isSnip = false
				snippets = append(snippets, snip)
				snip = Snip{}
				continue
			}

			snip.Template += strings.ReplaceAll(line, "\t", "  ") + "\n"
			continue
		}

		if strings.Index(line, "snippet ") == 0 {
			isSnip = true
			parts := strings.Split(line, "\"")

			switch len(parts) {
			case 1:
				// Hint is missing, but that's okay
				// noop
			case 3:
				snip.Hint = parts[1]
			default:
				sParts, _ := json.Marshal(parts)
				err := fmt.Errorf(
					"unexpected file format found during parsing: %s, parts: %s, err: %w",
					line,
					sParts,
					MalformedSnip,
				)
				return nil, err
			}

			frontParts := strings.Split(strings.TrimSpace(parts[0]), " ")

			if len(frontParts) != 2 {
				sFrontParts, _ := json.Marshal(frontParts)
				err := fmt.Errorf(
					"unexpected file format found during parsing: %s, frontParts: %s, err: %w",
					line,
					sFrontParts,
					MalformedSnip,
				)
				return nil, err
			}

			snip.Trigger = frontParts[1]
		}
	}

	if isSnip {
		return nil, UnclosedSnip
	}

	return snippets, nil
}
