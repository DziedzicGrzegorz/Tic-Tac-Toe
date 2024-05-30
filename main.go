package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var (
	messageStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFD700")). // Gold
			Align(lipgloss.Center)

	backgroundStyle = lipgloss.NewStyle().
			Align(lipgloss.Center).
			Width(80).
			Height(24)
)

type model struct {
	width  int
	height int
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m, tea.Quit
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		backgroundStyle = backgroundStyle.Width(m.width).Height(m.height)
	}
	return m, nil
}

func (m model) View() string {
	message := "GAME COMING SOON"
	view := messageStyle.Render(message)

	fullScreen := backgroundStyle.Render(
		lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, view))

	return fullScreen
}

func main() {
	var rootCmd = &cobra.Command{
		Use:   "game",
		Short: "Tic-Tac-Toe game",
		Long:  "A simple Tic-Tac-Toe game written in Go using the Bubble Tea library and the Lip Gloss library.",
		Run: func(cmd *cobra.Command, args []string) {
			p := tea.NewProgram(model{}, tea.WithAltScreen())
			//TODO in bubbletea v026.3 nil pointer to SetWindowTitle make issue on GitHub
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
