package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

type mode int

const (
	modeMenu mode = iota
	modeSinglePlayer
	modeMultiPlayer
	modeStats
)

type menuItem struct {
	mode mode
	name string
}

type model struct {
	width     int
	height    int
	mode      mode
	cursor    int
	menuItems []menuItem
}

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFFFF")).
			Margin(1, 0, 2, 0).
			Padding(0, 1).
			Align(lipgloss.Center)

	messageStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFD700")). // Gold
			Align(lipgloss.Center).
			Background(lipgloss.Color("#0000FF")).
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#FFD700"))

	backgroundStyle = lipgloss.NewStyle().
			Align(lipgloss.Center).
			Width(80).
			Height(24)

	normalStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	selectedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFD700")).Bold(true)
	subtleStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Margin(2, 0, 0, 0)
)

func initialModel(width, height int) model {
	return model{
		mode:   modeMenu,
		width:  width,
		height: height,
		cursor: 0,
		menuItems: []menuItem{
			{mode: modeSinglePlayer, name: "Single Player"},
			{mode: modeMultiPlayer, name: "Multiplayer"},
			{mode: modeStats, name: "View Stats"},
		},
	}
}

func (m model) Init() tea.Cmd {
	return tea.EnterAltScreen
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.mode {
		case modeMenu:
			switch msg.String() {
			case "up", "k":
				if m.cursor > 0 {
					m.cursor--
				}
			case "down", "j":
				if m.cursor < len(m.menuItems)-1 {
					m.cursor++
				}
			case "enter":
				m.mode = m.menuItems[m.cursor].mode
			case "ctrl+c", "q", "esc":
				return m, tea.Quit
			}
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		backgroundStyle = backgroundStyle.Width(m.width).Height(m.height)
	}
	return m, nil
}

func choicesView(m model) string {
	title := "What to do today?\n"

	var choices []string
	for i, choice := range m.menuItems {

		choiceStr := normalStyle.Render("[ ]  " + choice.name)

		if m.cursor == i {
			choiceStr = selectedStyle.Render("[x]  " + choice.name)
		}

		choices = append(choices, choiceStr)
	}

	menuSelect := lipgloss.JoinVertical(lipgloss.Left, choices...)

	footer := subtleStyle.Render("j/k, up/down: select | enter: choose | q, esc: quit")

	view := lipgloss.JoinVertical(lipgloss.Center,
		titleStyle.Render(title),
		menuSelect,
		subtleStyle.Render(footer),
	)

	fullScreen := backgroundStyle.Render(
		lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, view, lipgloss.WithWhitespaceChars(" ")),
	)

	return fullScreen
}

func (m model) View() string {
	var view string

	switch m.mode {
	case modeMenu:
		view = choicesView(m)
	case modeSinglePlayer:
		view = messageStyle.Render("Single Player Mode Selected")
	case modeMultiPlayer:
		view = messageStyle.Render("Multiplayer Mode Selected")
	case modeStats:
		view = messageStyle.Render("Game Stats: \n\n- Wins: 10\n- Losses: 5\n- Draws: 2") // Example stats
	}

	fullScreen := lipgloss.NewStyle().Align(lipgloss.Center).Render(
		backgroundStyle.Render(
			lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, view),
		),
	)

	return fullScreen
}

func main() {
	var rootCmd = &cobra.Command{
		Use:   "game",
		Short: "Tic-Tac-Toe game",
		Long:  "A simple Tic-Tac-Toe game written in Go using the Bubble Tea library and the Lip Gloss library.",
		Run: func(cmd *cobra.Command, args []string) {
			p := tea.NewProgram(initialModel(0, 0), tea.WithAltScreen())
			p.SetWindowTitle("Welcome to the game")
			if _, err := p.Run(); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
		},
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
