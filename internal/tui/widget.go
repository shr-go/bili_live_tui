package tui

import (
	"github.com/charmbracelet/lipgloss"
	"strconv"
)

var (
	shipLevelToString = map[uint8]string{
		0: "",
		1: "总",
		2: "提",
		3: "舰",
	}
	medalStyle = func(medal *medalInfo) string {
		shipString := shipLevelToString[medal.shipLevel]
		if shipString != "" {
			shipString = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#F87299")).
				Render(shipString)
		}
		nameString := lipgloss.NewStyle().
			MaxWidth(10).
			Align(lipgloss.Center).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color(medal.medalColor)).
			Render(medal.name)
		levelString := lipgloss.NewStyle().
			Width(2).
			Align(lipgloss.Right).
			Foreground(lipgloss.Color("#3C3C3C")).
			Background(lipgloss.Color("#FAFAFA")).
			Render(strconv.Itoa(int(medal.level)))
		return shipString + nameString + levelString
	}

	nameStyle = func(name string, nameColor string) string {
		if nameColor == "" {
			nameColor = "#FAFAFA"
		}
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color(nameColor)).
			Render(name + ":")
	}

	contentStyle = func(content string, contentColor string) string {
		if contentColor == "" {
			contentColor = "#FAFAFA"
		}
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color(contentColor)).
			Render(content)
	}
)

var (
	subtle    = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
	highlight = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	special   = lipgloss.AdaptiveColor{Light: "#43BF6D", Dark: "#73F59F"}

	divider = lipgloss.NewStyle().
		SetString("•").
		Padding(0, 1).
		Foreground(subtle).
		String()

	url = lipgloss.NewStyle().Foreground(special).Render

	activeTabBorder = lipgloss.Border{
		Top:         "─",
		Bottom:      " ",
		Left:        "│",
		Right:       "│",
		TopLeft:     "╭",
		TopRight:    "╮",
		BottomLeft:  "┘",
		BottomRight: "└",
	}

	tabBorder = lipgloss.Border{
		Top:         "─",
		Bottom:      "─",
		Left:        "│",
		Right:       "│",
		TopLeft:     "╭",
		TopRight:    "╮",
		BottomLeft:  "┴",
		BottomRight: "┴",
	}

	tab = lipgloss.NewStyle().
		Border(tabBorder, true).
		BorderForeground(highlight).
		Padding(0, 1)

	activeTab = tab.Copy().Border(activeTabBorder, true)

	tabGap = tab.Copy().
		BorderTop(false).
		BorderLeft(false).
		BorderRight(false)

	titleStyle = lipgloss.NewStyle().
			MarginLeft(1).
			MarginRight(5).
			Padding(0, 1).
			Italic(true).
			Foreground(lipgloss.Color("#FFF7DB")).
			SetString("Lip Gloss")

	descStyle = lipgloss.NewStyle().MarginTop(1)

	infoStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderTop(true).
			BorderForeground(subtle)

	dialogBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#874BFD")).
			Padding(1, 0).
			BorderTop(true).
			BorderLeft(true).
			BorderRight(true).
			BorderBottom(true)

	buttonStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFF7DB")).
			Background(lipgloss.Color("#888B7E")).
			Padding(0, 3).
			MarginTop(1)

	activeButtonStyle = buttonStyle.Copy().
				Foreground(lipgloss.Color("#FFF7DB")).
				Background(lipgloss.Color("#F25D94")).
				MarginRight(2).
				Underline(true)

	listStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), false, true, false, false).
			BorderForeground(subtle).
			MarginRight(2)

	listHeader = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderBottom(true).
			BorderForeground(subtle).
			MarginRight(2).
			Render

	listItem = lipgloss.NewStyle().PaddingLeft(2).Render

	checkMark = lipgloss.NewStyle().SetString("✓").
			Foreground(special).
			PaddingRight(1).
			String()

	listDone = func(s string) string {
		return checkMark + lipgloss.NewStyle().
			Strikethrough(true).
			Foreground(lipgloss.AdaptiveColor{Light: "#969B86", Dark: "#696969"}).
			Render(s)
	}

	statusNugget = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Padding(0, 1)

	statusBarStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#343433", Dark: "#C1C6B2"}).
			Background(lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#353533"})

	statusStyle = lipgloss.NewStyle().
			Inherit(statusBarStyle).
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#FF5F87")).
			Padding(0, 1).
			MarginRight(1)

	encodingStyle = statusNugget.Copy().
			Background(lipgloss.Color("#A550DF")).
			Align(lipgloss.Right)

	statusText = lipgloss.NewStyle().Inherit(statusBarStyle)

	fishCakeStyle = statusNugget.Copy().Background(lipgloss.Color("#6124DF"))

	// Page.

	docStyle = lipgloss.NewStyle().Padding(1, 2, 1, 2)

	focusedStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("69"))

	unFocusedStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.HiddenBorder())
)
