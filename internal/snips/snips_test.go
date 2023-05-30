package snips_test

import (
	"errors"
	"testing"

	"github.com/jhr22/daily-snip/internal/snips"
)

type FakeScanner struct {
	CurrentIndex int
	Lines        []string
}

func (fs *FakeScanner) Scan() bool {
	fs.CurrentIndex += 1
	return len(fs.Lines) > fs.CurrentIndex
}

func (fs *FakeScanner) Text() string {
	return fs.Lines[fs.CurrentIndex]
}

func NewFakeScanner(lines []string) *FakeScanner {
	return &FakeScanner{
		CurrentIndex: -1,
		Lines:        lines,
	}
}

var (
	parseLinesCases = []struct {
		Error error
		Lines []string
		Found int
	}{
		{
			Lines: []string{
				"// A comment",
				"",
				"",
				"snippet a \"Hint for A\"",
				"code {",
				"\t\tfor(a)",
				"}",
				"endsnippet",
			},
			Found: 1,
		},
		{
			Lines: []string{
				"// A comment",
				"",
				"",
				"snippet a \"Hint for A\"",
				"code {",
				"\t\tfor(a)",
				"}",
				"endsnippet",
				"",
				"snippet b \"Hint for B\"",
				"code {",
				"\t\tfor(b)",
				"}",
				"endsnippet",
			},
			Found: 2,
		},
		{
			Lines: []string{},
			Found: 0,
		},
		{
			Error: snips.UnclosedSnip,
			Lines: []string{
				"snippet unclosed \"Hint for unclosed\"",
			},
		},
		{
			Error: snips.MalformedSnip,
			Lines: []string{
				"snippet \"Hint for malformed\"",
				"endsnippet",
			},
		},
		{
			Error: snips.MalformedSnip,
			Lines: []string{
				"snippet ",
				"endsnippet",
			},
		},
	}
)

func TestParseLines(t *testing.T) {
	for _, c := range parseLinesCases {
		fs := NewFakeScanner(c.Lines)

		snippets, err := snips.ParseLines(fs)
		if c.Error != nil && !errors.Is(err, c.Error) {
			t.Errorf("expected error Got: %s, Want: %s", err, c.Error)
		} else if c.Error == nil && err != nil {
			t.Errorf("unexpected err: %s", err)
		}

		if len(snippets) != c.Found {
			t.Errorf("len(snippets) Got: %d, Want: %d", len(snippets), c.Found)
		}
	}
}
