package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Confirmation TUI component
type confirmationModel struct {
	prompt   string
	choice   bool
	quitting bool
}

func (m confirmationModel) Init() tea.Cmd { return nil }
func (m confirmationModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "y", "Y":
			m.choice = true
			m.quitting = true
			return m, tea.Quit
		case "n", "N", "enter":
			m.choice = false
			m.quitting = true
			return m, tea.Quit
		case "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		}
	}
	return m, nil
}
func (m confirmationModel) View() string {
	if m.quitting {
		return ""
	}
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("3")).Bold(true)
	return fmt.Sprintf("%s %s (y/N) ", style.Render("?"), m.prompt)
}

func PromptConfirmation(prompt string) (bool, error) {
	m := confirmationModel{prompt: prompt}
	p := tea.NewProgram(m)
	finalModel, err := p.Run()
	if err != nil {
		return false, err
	}
	return finalModel.(confirmationModel).choice, nil
}

// Menu Selection component
type menuModel struct {
	title    string
	options  []string
	cursor   int
	choice   int
	quitting bool
}

func (m menuModel) Init() tea.Cmd { return nil }
func (m menuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			m.choice = -1
			m.quitting = true
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.options)-1 {
				m.cursor++
			}
		case "enter":
			m.choice = m.cursor
			m.quitting = true
			return m, tea.Quit
		}
	}
	return m, nil
}
func (m menuModel) View() string {
	if m.quitting {
		return ""
	}
	s := fmt.Sprintf("%s\n\n", lipgloss.NewStyle().Bold(true).Render(m.title))
	for i, option := range m.options {
		cursorStr := " "
		if m.cursor == i {
			cursorStr = lipgloss.NewStyle().Foreground(lipgloss.Color("2")).Render(">")
			s += fmt.Sprintf("%s %s\n", cursorStr, lipgloss.NewStyle().Foreground(lipgloss.Color("2")).Render(option))
		} else {
			s += fmt.Sprintf("%s %s\n", cursorStr, option)
		}
	}
	return s
}

func PromptMenu(title string, options []string) (int, error) {
	m := menuModel{title: title, options: options, cursor: 0, choice: -1}
	p := tea.NewProgram(m)
	finalModel, err := p.Run()
	if err != nil {
		return -1, err
	}
	return finalModel.(menuModel).choice, nil
}

// Checkbox Multi-Select component
type checkboxModel struct {
	title    string
	options  []string
	selected map[int]bool
	cursor   int
	quitting bool
}

func (m checkboxModel) Init() tea.Cmd { return nil }
func (m checkboxModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.options)-1 {
				m.cursor++
			}
		case "space":
			m.selected[m.cursor] = !m.selected[m.cursor]
		case "enter":
			m.quitting = true
			return m, tea.Quit
		}
	}
	return m, nil
}
func (m checkboxModel) View() string {
	if m.quitting {
		return ""
	}
	s := fmt.Sprintf("%s (Space to toggle, Enter to confirm)\n\n", lipgloss.NewStyle().Bold(true).Render(m.title))
	for i, option := range m.options {
		cursorStr := " "
		if m.cursor == i {
			cursorStr = ">"
		}
		checked := " "
		if m.selected[i] {
			checked = "x"
		}
		s += fmt.Sprintf("%s [%s] %s\n", cursorStr, checked, option)
	}
	return s
}

func PromptCheckboxes(title string, options []string) ([]int, error) {
	m := checkboxModel{title: title, options: options, selected: make(map[int]bool), cursor: 0}
	p := tea.NewProgram(m)
	finalModel, err := p.Run()
	if err != nil {
		return nil, err
	}
	var choices []int
	for k, v := range finalModel.(checkboxModel).selected {
		if v {
			choices = append(choices, k)
		}
	}
	return choices, nil
}

// Spinner component wrapper
type spinnerModel struct {
	spinner  spinner.Model
	text     string
	quitting bool
}

func (m spinnerModel) Init() tea.Cmd { return m.spinner.Tick }
func (m spinnerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case tea.KeyMsg:
		m.quitting = true
		return m, tea.Quit
	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
}
func (m spinnerModel) View() string {
	if m.quitting {
		return ""
	}
	return fmt.Sprintf("\n %s %s\n\n", m.spinner.View(), m.text)
}
