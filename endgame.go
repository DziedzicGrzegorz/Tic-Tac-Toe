package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	highlightStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#2ECC71")).Align(lipgloss.Center)
	drawStyle      = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#F4D03F")).Align(lipgloss.Center)
)

type EndGameModel struct {
	width   int
	height  int
	message string
	isDraw  bool
}

func NewEndGameModel(width, height int, endGameMessage string, isDraw bool) *EndGameModel {
	return &EndGameModel{
		width:   width,
		height:  height,
		message: endGameMessage,
		isDraw:  isDraw,
	}
}

func (m *EndGameModel) Init() tea.Cmd {
	return nil
}

func (m *EndGameModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "r":
			return NewGameModel(m.width, m.height), nil
		case "m":
			return initialModel(m.width, m.height), nil
		case "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m *EndGameModel) View() string {
	message := fmt.Sprintf("%s\n\nPress 'r' to restart or 'm' to return to menu.", m.message)
	styledMessage := highlightStyle.Render(message)
	if m.isDraw {
		styledMessage = drawStyle.Render(message)
	}
	centeredMessage := lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, styledMessage)
	return centeredMessage
}
