package main

import (
	"fmt"
	"net"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dziedzicgrzegorz/Tic-Tac-Toe/constants"
)

type moveMessage struct{ command string }
type errMsg struct{ err error }

func (e errMsg) Error() string {
	return e.err.Error()
}

type TCPmodel struct {
	board          [constants.BoardSize][constants.BoardSize]int
	selectedRow    int
	selectedColumn int
	winner         string
	conn           *net.Conn
	player         int
	playerTurn     int
	width          int
	height         int
	errorMessage   string
	infoMessage    string
}

func newTCPModel(width, height int, conn *net.Conn, player int) TCPmodel {
	return TCPmodel{
		board:          [constants.BoardSize][constants.BoardSize]int{},
		selectedRow:    0,
		selectedColumn: 0,
		winner:         "",
		conn:           conn,
		player:         player,
		playerTurn:     constants.PlayerX,
		width:          width,
		height:         height,
	}
}

func (m TCPmodel) Init() tea.Cmd {
	return createReceiveMove(*m.conn)
}

func (m TCPmodel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var err error

	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case constants.CtrlC, constants.Esc:
			return m, tea.Quit

		case constants.Up:
			if m.selectedRow > 0 {
				m.selectedRow--
			}
		case constants.Down:
			if m.selectedRow < constants.BoardSize-1 {
				m.selectedRow++
			}
		case constants.Left:
			if m.selectedColumn > 0 {
				m.selectedColumn--
			}
		case constants.Right:
			if m.selectedColumn < constants.BoardSize-1 {
				m.selectedColumn++
			}
		case constants.Enter:
			m, err = m.HandleMyEnter()
			if err != nil {
				m.errorMessage = err.Error()
				return NewEndGameModel(m.width, m.height, m.errorMessage), nil
			}
			if val := m.checkWinner(); val != 0 {
				winMsg := fmt.Sprintf("Player %s wins!", mapValueToMarker(val))
				return NewEndGameModel(m.width, m.height, constants.WinMsgStyle.Render(winMsg)), nil
			}
			if m.isDraw() {
				drawMsg := "It's a draw!"
				return NewEndGameModel(m.width, m.height, constants.DrawMsgStyle.Render(drawMsg)), nil
			}
		}

		return m, createReceiveMove(*m.conn)

	case moveMessage:
		commandParts := strings.Split(msg.command, ",")

		if commandParts[0] == constants.Enter {
			m, err = m.HandleOpponentEnter(commandParts[0], commandParts[1], commandParts[2], commandParts[3])
			if err != nil {
				m.errorMessage = err.Error()
				return NewEndGameModel(m.width, m.height, m.errorMessage), nil
			}
			if val := m.checkWinner(); val != 0 {
				loseMsg := fmt.Sprintf("Player %s wins!", mapValueToMarker(val))
				return NewEndGameModel(m.width, m.height, constants.LoseMsgStyle.Render(loseMsg)), nil
			}
			if m.isDraw() {
				drawMsg := "It's a draw!"
				return NewEndGameModel(m.width, m.height, constants.DrawMsgStyle.Render(drawMsg)), nil
			}
		}

		return m, nil

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	return m, nil
}

func (m TCPmodel) View() string {
	var cells []string
	for i := 0; i < constants.BoardSize*constants.BoardSize; i++ {
		row, col := i/constants.BoardSize, i%constants.BoardSize
		cell := m.board[row][col]
		cellStr := " "

		style := constants.CellStyle

		if row == m.selectedRow && col == m.selectedColumn {
			cellStr = constants.BlinkingStyle.Render(m.getCurrentUser())
		} else if cell == constants.PlayerX {
			cellStr = "X"
		} else if cell == constants.PlayerO {
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

	currentPlayer := fmt.Sprintf("I am a %s player: \n", m.getCurrentUser())

	header := constants.HeaderStyle.Render(currentPlayer)

	// Instructions
	footer := constants.SubtleStyle.Render("arrow keys: move | enter: select | ctrl+c or Esc: quit")

	// Joining all cells horizontally
	board := lipgloss.JoinVertical(lipgloss.Left,
		lipgloss.JoinHorizontal(lipgloss.Center, cells[0], cells[1], cells[2]),
		lipgloss.JoinHorizontal(lipgloss.Left, cells[3], cells[4], cells[5]),
		lipgloss.JoinHorizontal(lipgloss.Left, cells[6], cells[7], cells[8]))

	whoseTurn := fmt.Sprintf("It's %s's turn.\n", m.getCurrentMarker())

	infoMsg := constants.InfoStyle.Render(m.infoMessage)

	// Joining all elements vertically
	view := lipgloss.JoinVertical(lipgloss.Center,
		header,
		infoMsg,
		whoseTurn,
		constants.BoardStyle.Render(board),
		m.errorMessage,
		footer,
	)

	centeredFullView := lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, view)

	return centeredFullView
}

func (m TCPmodel) getCurrentMarker() string {
	if m.playerTurn == constants.PlayerX {
		return "X"
	}
	return "O"
}
func (m TCPmodel) getCurrentUser() string {
	if m.player == constants.PlayerX {
		return "X"
	}
	return "O"
}
func mapValueToMarker(value int) string {
	if value == constants.PlayerX {
		return "X"
	}
	if value == constants.PlayerO {
		return "O"
	}
	return " "
}

func (m TCPmodel) HandleMyEnter() (TCPmodel, error) {
	m = m.handlePlayerEnter(m.player, m.selectedRow, m.selectedColumn)
	err := m.sendMove(constants.Enter)
	return m, err
}

func (m TCPmodel) HandleOpponentEnter(key, opponentStr, selectedRowStr, selectedColStr string) (TCPmodel, error) {
	selectedRow, err := strconv.Atoi(selectedRowStr)
	if err != nil {
		return m, err
	}

	selectedCol, err := strconv.Atoi(selectedColStr)
	if err != nil {
		return m, err
	}

	opponent, err := strconv.Atoi(opponentStr)
	if err != nil {
		return m, err
	}

	return m.handlePlayerEnter(opponent, selectedRow, selectedCol), nil
}

func (m TCPmodel) handlePlayerEnter(player, row, col int) TCPmodel {
	if player != m.playerTurn {
		m.infoMessage = fmt.Sprintf("Ignoring %d's move as it's %d's turn.", player, m.playerTurn)
		return m
	}

	if m.board[row][col] != constants.Empty {
		m.infoMessage = fmt.Sprintf("Ignoring %s's move as cell [%d, %d] is not Empty.", m.getCurrentMarker(), row, col)
		return m
	}

	m.board[row][col] = player
	m.infoMessage = fmt.Sprintf("%s marked cell [%d, %d]", m.getCurrentMarker(), row, col)

	m.switchPlayer()

	return m
}

func (m *TCPmodel) checkWinner() int {
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
		// Check if all cells in the line are the same and not Empty
		if m.board[a[0]][a[1]] != constants.Empty &&
			m.board[a[0]][a[1]] == m.board[b[0]][b[1]] &&
			m.board[a[0]][a[1]] == m.board[c[0]][c[1]] {
			// Return the winning marker
			return m.board[a[0]][a[1]]
		}
	}

	return 0
}

func (m *TCPmodel) isDraw() bool {
	boardFull := true
	for i := 0; i < constants.BoardSize; i++ {
		for j := 0; j < constants.BoardSize; j++ {
			if m.board[i][j] == constants.Empty {
				boardFull = false
				break
			}
		}
	}

	return boardFull
}

func (m *TCPmodel) switchPlayer() {
	m.playerTurn = -m.playerTurn
}

func (m *TCPmodel) sendMove(key string) error {
	move := fmt.Sprintf("%s,%d,%d,%d", key, m.player, m.selectedRow, m.selectedColumn)
	_, err := (*m.conn).Write([]byte((move)))
	return err
}
