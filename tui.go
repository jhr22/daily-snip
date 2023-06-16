package main

import (
	"fmt"
	"io"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/jhr22/daily-snip/internal/snips"
)

var (
	lightPurple = lipgloss.Color("#D48DFC")
	darkPurple  = lipgloss.Color("#A657D4")

	selectionItemStyle = lipgloss.NewStyle().Bold(true)
	listViewStyle      = lipgloss.NewStyle().
				Foreground(darkPurple).
				Padding(2)
	selectionViewStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lightPurple).
				Padding(2).
				BorderForeground(darkPurple).
				Border(lipgloss.NormalBorder(), true, false, false, false)
	windowViewStyle = lipgloss.NewStyle().
			BorderForeground(darkPurple).
			Border(lipgloss.NormalBorder(), true)
)

type item struct {
	snip snips.Snip
}

func (i item) FilterValue() string { return i.snip.String() }

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s (%s)", index+1, i.snip.Trigger, i.snip.Hint)

	if index == m.Index() {
		str = selectionItemStyle.Render(str)
	}

	fmt.Fprint(w, str)
}

type model struct {
	list     list.Model
	chosen   bool
	quitting bool
}

func initialModel(snippets []snips.Snip) model {
	items := make([]list.Item, 0, len(snippets))

	for _, snip := range snippets {
		items = append(items, item{snip: snip})
	}

	list := list.New(items, itemDelegate{}, 100, 24)

	list.Title = "Browse vim-go snippets ..."
	list.Styles.Title.Bold(true).Background(lightPurple)

	return model{
		list: list,
	}
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		selectionViewStyle.Width(msg.Width - 4)
		selectionViewStyle.Height(msg.Height - m.list.Height() - 8)
		return m, nil

	// Is it a key press?
	case tea.KeyMsg:
		if m.list.FilterState() != list.Filtering {
			// Cool, what was the actual key pressed?
			switch msg.String() {

			// These keys should exit the program.
			case "ctrl+c", "q":
				m.quitting = true
				return m, tea.Quit

			case "enter":
				m.chosen = true
				return m, tea.Quit
			}
		}
	}

	var cmd tea.Cmd

	m.list, cmd = m.list.Update(msg)

	return m, cmd
}

func (m model) View() string {
	if m.quitting {
		return ""
	}

	selection := ""

	if m.list.SelectedItem() != nil {
		snip := m.list.SelectedItem().(item).snip
		selection = fmt.Sprintf("\n%s\n", snip)
	}

	return windowViewStyle.Render(
		lipgloss.JoinVertical(
			lipgloss.Top,
			listViewStyle.Render(m.list.View()),
			selectionViewStyle.Render(selection),
		),
	)
}

func runTUI(snippets []snips.Snip) {
	p := tea.NewProgram(initialModel(snippets), tea.WithAltScreen())
	mt, err := p.Run()

	if err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}

	if mt == nil {
		fmt.Println("Woops no model returned probably an error")
		os.Exit(1)
	}

	m := mt.(model)

	if m.chosen {
		fmt.Println(m.list.SelectedItem().(item).snip)
	}
}
