package main

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"strings"
)

var (
	activeBorder   = lipgloss.Color("12")
	inactiveBorder = lipgloss.Color("8")
)

type model struct {
	viewport          viewport.Model
	currentChoices    []string
	choiceCursor      int
	windowWidth       int
	windowHeight      int
	choiceAreaFocused bool
}

func (m model) Init() tea.Cmd {
	return nil
}

const sidebarWidth = 30

func (m model) getMainWidth() int {
	return m.windowWidth - (sidebarWidth + 1)
}
func (m model) getChoiceHeight() int {
	return len(m.currentChoices) + 2
}
func (m model) getMainHeight() int {
	return m.windowHeight - (m.getChoiceHeight() + 2)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// Window sizing
		m.windowWidth = msg.Width - 2
		m.windowHeight = msg.Height - 2

		m.viewport.Width = m.getMainWidth()
		m.viewport.Height = m.getMainHeight()

	case tea.KeyMsg:
		// Overall key inputs
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "tab":
			m.choiceAreaFocused = !m.choiceAreaFocused
		}
	}

	var cmd tea.Cmd
	if m.choiceAreaFocused {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "up", "k":
				if m.choiceCursor > 0 {
					m.choiceCursor--
				}
			case "down", "j":
				if m.choiceCursor < len(m.currentChoices)-1 {
					m.choiceCursor++
				}
			case "enter":
				// Handle selection
			}
		}
	} else {
		m.viewport, cmd = m.viewport.Update(msg)
	}

	return m, cmd
}

func condColor(cond bool, yes, no lipgloss.TerminalColor) lipgloss.TerminalColor {
	if cond {
		return yes
	}
	return no
}

func (m model) View() string {
	lw := m.getMainWidth()
	ch := m.getChoiceHeight()
	vh := m.getMainHeight()

	rightStyle := lipgloss.NewStyle().
		Width(sidebarWidth).
		Height(m.windowHeight).
		PaddingLeft(1)

	viewportStyle := lipgloss.NewStyle().
		Width(lw).
		Height(vh - 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(condColor(m.choiceAreaFocused, inactiveBorder, activeBorder)).
		PaddingLeft(6)

	choiceStyle := lipgloss.NewStyle().
		Width(lw).
		Border(lipgloss.RoundedBorder()).
		Height(ch - 2).
		BorderForeground(condColor(m.choiceAreaFocused, activeBorder, inactiveBorder))

	// Build choices
	var choiceContent string
	for i, c := range m.currentChoices {
		cursor := "  "
		if i == m.choiceCursor {
			cursor = "> "
		}
		choiceContent += cursor + c + "\n"
	}

	// Left column: viewport stacked on choices
	leftColumn := lipgloss.JoinVertical(
		lipgloss.Left,
		viewportStyle.Render(m.viewport.View()),
		choiceStyle.Render(choiceContent),
	)

	// Right column
	rightColumn := rightStyle.Render("lorem ipsum")

	full := lipgloss.JoinHorizontal(lipgloss.Top, leftColumn, rightColumn)
	return lipgloss.NewStyle().Padding(1).Render(full)
}

func main() {
	m := model{
		viewport:          viewport.New(80, 20),
		currentChoices:    []string{"Option A", "Option B", "Option C"},
		choiceAreaFocused: true,
	}
	m.viewport.SetContent("Your scrollable content here...\n" + strings.Repeat("Line\n", 50))

	p := tea.NewProgram(m, tea.WithAltScreen())
	p.Run()
}
