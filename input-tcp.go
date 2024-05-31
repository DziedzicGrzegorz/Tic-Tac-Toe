package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dziedzicgrzegorz/Tic-Tac-Toe/constants"
)

const (
	waitPlaceholder    = "Are you Client? Type C otherwise !C"
	ipPlaceholder      = "IP Address"
	portPlaceholder    = "Port"
	submitButtonText   = "Submit"
	cursorModeHelp     = "cursor mode is %s (ctrl+r to change style)"
	defaultIPAddress   = "localhost"
	defaultPort        = "8080"
	maxErrorMessageLen = 60 // Maximum length of the error message
)

var (
	focusedButton = constants.FocusedStyle.Render("[ " + submitButtonText + " ]")
	blurredButton = fmt.Sprintf("[ %s ]", constants.BlurredStyle.Render(submitButtonText))
)

type TcpInputModel struct {
	focusIndex   int
	inputs       []textinput.Model
	cursorMode   cursor.Mode
	height       int
	width        int
	errorMessage string
}

func NewTCPInputModel(width, height int) TcpInputModel {
	m := TcpInputModel{
		inputs: make([]textinput.Model, 3),
		width:  width,
		height: height,
	}

	var t textinput.Model
	for i := range m.inputs {
		t = textinput.New()
		t.Cursor.Style = constants.FocusedStyle
		t.CharLimit = 32

		switch i {
		case 0:
			t.Placeholder = waitPlaceholder
			t.Focus()
			t.Width = len(waitPlaceholder)
			t.PromptStyle = constants.FocusedStyle
			t.TextStyle = constants.FocusedStyle
		case 1:
			t.Placeholder = ipPlaceholder
			t.Width = len(ipPlaceholder)
			t.CharLimit = 15
		case 2:
			t.Placeholder = portPlaceholder
			t.Width = len(portPlaceholder)
			t.CharLimit = 5
		}

		m.inputs[i] = t
	}

	return m
}

func (m TcpInputModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m TcpInputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {

		case constants.CtrlC, constants.Esc:
			return m, tea.Quit

		// Change cursor mode
		case constants.CtrlR:
			m.cursorMode++
			if m.cursorMode > cursor.CursorHide {
				m.cursorMode = cursor.CursorBlink
			}
			cmds := make([]tea.Cmd, len(m.inputs))
			for i := range m.inputs {
				cmds[i] = m.inputs[i].Cursor.SetMode(m.cursorMode)
			}
			return m, tea.Batch(cmds...)

		// Set focus to next input
		case constants.Tab, constants.ShiftTab, constants.Enter, constants.Up, constants.Down:
			s := msg.String()

			// Did the user press enter while the submit button was focused?
			// If so, attempt to start the game.
			if s == constants.Enter && m.focusIndex == len(m.inputs) {
				wait := strings.ToUpper(m.inputs[0].Value()) != "C"
				ip := m.inputs[1].Value()
				port := m.inputs[2].Value()

				// Use default values if inputs are empty
				if ip == "" {
					ip = defaultIPAddress
				}
				if port == "" {
					port = defaultPort
				}
				//after submit button freeze because waiting for connection
				conn, player, err := setupConnection(wait, ip, port)
				if err != nil {
					m.errorMessage = formatErrorMessage(err.Error())
					return m, nil
				}
				game := newTCPModel(m.width, m.height, &conn, player)
				return game, nil
			}

			// Cycle indexes
			if s == constants.Up || s == constants.ShiftTab {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > len(m.inputs) {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = len(m.inputs)
			}

			cmds := make([]tea.Cmd, len(m.inputs))
			for i := 0; i <= len(m.inputs)-1; i++ {
				if i == m.focusIndex {
					// Set focused state
					cmds[i] = m.inputs[i].Focus()
					m.inputs[i].PromptStyle = constants.FocusedStyle
					m.inputs[i].TextStyle = constants.FocusedStyle
					continue
				}
				// Remove focused state
				m.inputs[i].Blur()
				m.inputs[i].PromptStyle = constants.NoStyle
				m.inputs[i].TextStyle = constants.NoStyle
			}

			return m, tea.Batch(cmds...)
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	// Handle character input and blinking
	cmd := m.updateInputs(msg)

	return m, cmd
}

func (m *TcpInputModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m TcpInputModel) View() string {
	// Collecting all inputs into a vertical layout
	var inputs []string
	for i := range m.inputs {
		inputs = append(inputs, m.inputs[i].View())
	}
	inputsView := lipgloss.JoinVertical(lipgloss.Left, inputs...)

	// Setting the button view
	button := blurredButton
	if m.focusIndex == len(m.inputs) {
		button = focusedButton
	}
	buttonView := fmt.Sprintf("\n\n%s\n\n", button)

	// Combining inputs and button
	mainView := lipgloss.JoinVertical(lipgloss.Center, inputsView, buttonView)

	// Adding the help text
	helpText := fmt.Sprintf(cursorModeHelp, m.cursorMode.String())
	helpView := lipgloss.JoinVertical(lipgloss.Center, mainView, constants.BlurredStyle.Render(helpText))

	// Error message view
	errorMsg := ""
	if m.errorMessage != "" {
		errorMsg = constants.ErrorStyle.Margin(2, 0, 0, 0).Render(m.errorMessage)
	}

	// Placing everything in the center of the screen
	fullScreen := lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, lipgloss.JoinVertical(lipgloss.Center, helpView, errorMsg), lipgloss.WithWhitespaceChars(" "))

	return fullScreen
}

func formatErrorMessage(msg string) string {
	if len(msg) <= maxErrorMessageLen {
		return msg
	}
	var formattedMsg strings.Builder
	for i, r := range msg {
		if i > 0 && i%maxErrorMessageLen == 0 {
			formattedMsg.WriteString("\n")
		}
		formattedMsg.WriteRune(r)
	}
	return formattedMsg.String()
}
