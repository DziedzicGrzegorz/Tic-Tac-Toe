package constants

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	InfoStyle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#1F618D")).Align(lipgloss.Center).Margin(0, 0, 1, 0)
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFFFF")).
			Margin(1, 0, 2, 0).
			Padding(0, 1).
			Align(lipgloss.Center)

	BackgroundStyle = lipgloss.NewStyle().
			Align(lipgloss.Center).
			Width(80).
			Height(24)

	NormalStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	SelectedStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFD700")).Bold(true)
	SubtleStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Margin(2, 0, 0, 0)
	CellStyle      = lipgloss.NewStyle().Align(lipgloss.Center, lipgloss.Center).Width(5).Height(3).Border(lipgloss.NormalBorder())
	BoardStyle     = lipgloss.NewStyle().Border(lipgloss.RoundedBorder())
	HeaderStyle    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FFD700")).Align(lipgloss.Center).Margin(0, 0, 2, 0)
	ErrorStyle     = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FF0000")).Align(lipgloss.Center)
	TCPErrStyle    = ErrorStyle.Margin(2, 0, 0, 0)
	BlinkingStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).Bold(true).Blink(true)
	HighlightStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("241")).Align(lipgloss.Center)
	DrawMsgStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#F1C40F")).Align(lipgloss.Center)
	WinMsgStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#2ECC71 ")).Align(lipgloss.Center)
	LoseMsgStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FF0000")).Align(lipgloss.Center)
	FocusedStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFD700"))
	BlurredStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	NoStyle        = lipgloss.NewStyle()
)
