package main

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dziedzicgrzegorz/Tic-Tac-Toe/constants"
)

type GameModel struct {
	width        int
	height       int
	cursor       int
	board        [3][3]int
	current      int
	winner       int
	blink        bool
	turnCount    int
	errorMessage string
}

func NewGameModel(width, height int) *GameModel {
	return &GameModel{
		width:     width,
		height:    height,
		cursor:    0,
		current:   1, // X starts
		board:     [3][3]int{},
		turnCount: 9,
	}
}

func (m *GameModel) Init() tea.Cmd {
	return tea.Batch(tea.EnterAltScreen)
}

func (m *GameModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			m.moveCursor(-3)
		case "down", "j":
			m.moveCursor(3)
		case "left", "h":
			m.moveCursor(-1)
		case "right", "l":
			m.moveCursor(1)
		case "enter":
			if !m.placeMarker() {
				m.errorMessage = "Cannot overwrite existing marker!"
			} else {
				m.errorMessage = ""
				m.winner = m.checkWinner()
				if m.winner != 0 {
					endMessage := fmt.Sprintf("Player %s wins!", m.currentMarker())
					endGameModel := NewEndGameModel(m.width, m.height, constants.WinMsgStyle.Render(endMessage))
					//sleep for 500 ms for better UX
					time.Sleep(500 * time.Millisecond)
					return endGameModel, endGameModel.Init()
				}
				m.turnCount--
				if m.turnCount == 0 {
					endMessage := "It's a draw!"

					endGameModel := NewEndGameModel(m.width, m.height, constants.DrawMsgStyle.Render(endMessage))
					//sleep for 500 ms for better UX
					time.Sleep(500 * time.Millisecond)
					return endGameModel, endGameModel.Init()
				}
				m.switchPlayer()
			}
		case "ctrl+c", "q", "esc":
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}
	return m, nil
}

func (m *GameModel) View() string {
	return m.renderBoard()
}

func (m *GameModel) moveCursor(delta int) {
	// Clear the error message when move is made
	m.errorMessage = ""

	newCursor := m.cursor + delta
	// Get the column
	col := m.cursor % 3

	// Calculate the new row and column
	newRow, newCol := newCursor/3, newCursor%3

	// Check for boundary conditions
	if newRow < 0 || newRow >= 3 || newCol < 0 || newCol >= 3 {
		newCursor = m.cursor
	}

	// Handle wrapping around the edges
	if (col == 2 && delta == 1) || (col == 0 && delta == -1) {
		newCursor = m.cursor
	}

	m.cursor = newCursor
}

func (m *GameModel) placeMarker() bool {
	row, col := m.cursor/3, m.cursor%3
	// Check if the cell is empty
	if m.board[row][col] != 0 {
		return false
	}
	// Place the marker
	m.board[row][col] = m.current
	return true
}

func (m *GameModel) switchPlayer() {
	// Switch player if 1 (X) or -1 (O)
	m.current = -m.current
}

func (m *GameModel) currentMarker() string {
	if m.current == 1 {
		return "X"
	}
	return "O"
}

func (m *GameModel) checkWinner() int {
	winningLines := [][3][2]int{
		{{0, 0}, {0, 1}, {0, 2}}, // first row
		{{1, 0}, {1, 1}, {1, 2}}, // second row
		{{2, 0}, {2, 1}, {2, 2}}, // third row
		{{0, 0}, {1, 0}, {2, 0}}, // first column
		{{0, 1}, {1, 1}, {2, 1}}, // second column
		{{0, 2}, {1, 2}, {2, 2}}, // third column
		{{0, 0}, {1, 1}, {2, 2}}, // diagonal top-left to bottom-right
		{{0, 2}, {1, 1}, {2, 0}}, // diagonal top-right to bottom-left
	}

	for _, line := range winningLines {
		a, b, c := line[0], line[1], line[2]
		// Check if all cells in the line are the same and not empty
		if m.board[a[0]][a[1]] != 0 &&
			m.board[a[0]][a[1]] == m.board[b[0]][b[1]] &&
			m.board[a[0]][a[1]] == m.board[c[0]][c[1]] {
			// Return the winning marker
			return m.board[a[0]][a[1]]
		}
	}

	return 0
}

func (m *GameModel) renderBoard() string {
	var cells []string
	for i := 0; i < 9; i++ {
		row, col := i/3, i%3
		cell := m.board[row][col]
		cellStr := " "

		style := constants.CellStyle

		if i == m.cursor {
			cellStr = constants.BlinkingStyle.Render(m.currentMarker())
		} else if cell == 1 {
			cellStr = "X"
		} else if cell == -1 {
			cellStr = "O"
		}

		// Border styling for the cells because if we set a border for each side it has a margin
		if col == 0 || col == 2 {
			style = style.UnsetBorderLeft().UnsetBorderRight().UnsetBorderBottom()
			if row == 0 {
				style = style.UnsetBorderStyle()
			}
		}
		if col == 1 {
			style = style.UnsetBorderBottom()
			if row == 0 {
				style = style.UnsetBorderTop()
			}
		}

		cells = append(cells, style.Render(cellStr))
	}

	currentPlayer := fmt.Sprintf("Current player: %s\n", m.currentMarker())

	header := constants.HeaderStyle.Render(currentPlayer)

	// Quick help
	footer := constants.SubtleStyle.Render("j/k, up/down: move | h/l, left/right: move | enter: select | ctrl+c: quit")

	// Joining all cells horizontally
	board := lipgloss.JoinVertical(lipgloss.Left,
		lipgloss.JoinHorizontal(lipgloss.Center, cells[0], cells[1], cells[2]),
		lipgloss.JoinHorizontal(lipgloss.Left, cells[3], cells[4], cells[5]),
		lipgloss.JoinHorizontal(lipgloss.Left, cells[6], cells[7], cells[8]))

	errorMsg := ""
	if m.errorMessage != "" {
		errorMsg = constants.ErrorStyle.Render(m.errorMessage)
	}

	// Joining all elements vertically
	view := lipgloss.JoinVertical(lipgloss.Center,
		header,
		constants.BoardStyle.Render(board),
		errorMsg,
		footer,
	)

	centeredFullView := lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, view)

	return centeredFullView
}
