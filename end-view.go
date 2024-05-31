package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dziedzicgrzegorz/Tic-Tac-Toe/constants"
)

type EndGameModel struct {
	width   int
	height  int
	message string
}

func NewEndGameModel(width, height int, endGameMessage string) *EndGameModel {
	return &EndGameModel{
		width:   width,
		height:  height,
		message: endGameMessage,
	}
}

func (m *EndGameModel) Init() tea.Cmd {
	return nil
}

func (m *EndGameModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case constants.M:
			return initialModel(m.width, m.height), nil
		case constants.Quit, constants.CtrlC, constants.Esc:
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}
	return m, nil
}

func (m *EndGameModel) View() string {
	message := fmt.Sprintf("%s\n\nPress 'm' to return to menu.", m.message)
	styledMessage := constants.HighlightStyle.Render(message)

	centeredMessage := lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, styledMessage)
	return centeredMessage
}
