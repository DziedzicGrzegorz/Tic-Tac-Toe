package main

import (
	"fmt"
	"net"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dziedzicgrzegorz/Tic-Tac-Toe/constants"
	"github.com/spf13/cobra"
)

type mode int

const (
	modeMenu mode = iota
	modeMultiPlayer
	modeMultiTCP
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
	conn      net.Conn
	player    string
}

func initialModel(width, height int) model {
	return model{
		mode:   modeMenu,
		width:  width,
		height: height,
		cursor: 0,
		menuItems: []menuItem{
			{mode: modeMultiPlayer, name: "Multiplayer"},
			{mode: modeMultiTCP, name: "Multiplayer TCP"},
		},
	}
}

func (m model) Init() tea.Cmd {
	return tea.SetWindowTitle("Tic-Tac-Toe Game")
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.mode {
		case modeMenu:
			switch msg.String() {
			case constants.Up:
				if m.cursor > 0 {
					m.cursor--
				}
			case constants.Down:
				if m.cursor < len(m.menuItems)-1 {
					m.cursor++
				}
			case constants.Enter:
				switch m.cursor {
				case 0:
					game := NewGameModel(m.width, m.height)
					return game, nil
				case 1:
					tcpInputModel := NewTCPInputModel(m.width, m.height)
					return tcpInputModel, nil

				}
			case constants.Quit, constants.CtrlC:
				return m, tea.Quit
			}
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		constants.BackgroundStyle = constants.BackgroundStyle.Width(m.width).Height(m.height)
	}
	return m, nil
}

func choicesView(m model) string {
	title := "What to do today?\n"

	var choices []string
	for i, choice := range m.menuItems {
		choiceStr := constants.NormalStyle.Render("[ ]  " + choice.name)
		if m.cursor == i {
			choiceStr = constants.SelectedStyle.Render("[x]  " + choice.name)
		}
		choices = append(choices, choiceStr)
	}

	menuSelect := lipgloss.JoinVertical(lipgloss.Left, choices...)
	footer := constants.SubtleStyle.Render("up ↑ / down ↓ : select | enter: choose | q, esc: quit")

	view := lipgloss.JoinVertical(lipgloss.Center,
		constants.TitleStyle.Render(title),
		menuSelect,
		constants.SubtleStyle.Render(footer),
	)

	fullScreen := constants.BackgroundStyle.Render(
		lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, view, lipgloss.WithWhitespaceChars(" ")),
	)

	return fullScreen
}

func (m model) View() string {
	var view string
	switch m.mode {
	case modeMenu:
		view = choicesView(m)
	}
	fullScreen := lipgloss.NewStyle().Align(lipgloss.Center).Render(
		constants.BackgroundStyle.Render(
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
