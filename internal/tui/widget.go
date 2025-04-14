package tui

import (
	"strconv"

	"github.com/charmbracelet/lipgloss"
)

var (
	getColor = func(color string, overrideAsBlack ...bool) lipgloss.Color {
		if !LiveConfig.ColorMode {
			if len(overrideAsBlack) > 0 && overrideAsBlack[0] {
				return lipgloss.Color("#3C3C3C")
			}
			return lipgloss.Color("#FAFAFA")
		}
		return lipgloss.Color(color)
	}

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
				Foreground(getColor("#F87299")).
				Render(shipString)
		}
		nameString := lipgloss.NewStyle().
			MaxWidth(10).
			Align(lipgloss.Center).
			Foreground(getColor("#FAFAFA")).
			Background(getColor(medal.medalColor, true)).
			Render(medal.name)
		levelString := lipgloss.NewStyle().
			Width(2).
			Align(lipgloss.Right).
			Foreground(getColor("#3C3C3C")).
			Background(getColor("#FAFAFA", true)).
			Render(strconv.Itoa(int(medal.level)))
		if !LiveConfig.ShowMedalLevel {
			levelString = ""
		}
		if !LiveConfig.ShowMedalName {
			nameString = ""
		}
		if !LiveConfig.ShowShipLevel {
			shipString = ""
		}
		if LiveConfig.ShowMedalLevel || LiveConfig.ShowMedalName || LiveConfig.ShowShipLevel {
			return shipString + nameString + levelString + " "
		} else {
			return shipString + nameString + levelString
		}
	}

	nameStyle = func(name string, nameColor string) string {
		if nameColor == "" {
			nameColor = "#FAFAFA"
		}
		return lipgloss.NewStyle().
			Foreground(getColor(nameColor)).
			Render(name + ":")
	}

	contentStyle = func(content string, contentColor string) string {
		if contentColor == "" {
			contentColor = "#FAFAFA"
		}
		return lipgloss.NewStyle().
			Foreground(getColor(contentColor)).
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

	urlStyle = lipgloss.NewStyle().Foreground(special).Render

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
			Foreground(getColor("#FFF7DB")).
			SetString("Lip Gloss")

	descStyle = lipgloss.NewStyle().MarginTop(1)

	infoStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderTop(true).
			BorderForeground(subtle)

	dialogBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(getColor("#874BFD")).
			Padding(1, 0).
			BorderTop(true).
			BorderLeft(true).
			BorderRight(true).
			BorderBottom(true)

	buttonStyle = lipgloss.NewStyle().
			Foreground(getColor("#FFF7DB")).
			Background(getColor("#888B7E")).
			Padding(0, 3).
			MarginTop(1)

	activeButtonStyle = buttonStyle.Copy().
				Foreground(getColor("#FFF7DB")).
				Background(getColor("#F25D94")).
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
			Foreground(getColor("#FFFDF5")).
			Padding(0, 1)

	statusBarStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#343433", Dark: "#C1C6B2"}).
			Background(lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#353533"})

	statusStyle = lipgloss.NewStyle().
			Inherit(statusBarStyle).
			Foreground(getColor("#FFFDF5")).
			Background(getColor("#FF5F87")).
			Padding(0, 1).
			MarginRight(1)

	encodingStyle = statusNugget.Copy().
			Background(getColor("#A550DF")).
			Align(lipgloss.Right)

	statusText = lipgloss.NewStyle().Inherit(statusBarStyle)

	fishCakeStyle = statusNugget.Copy().Background(getColor("#6124DF"))

	// Page.

	docStyle = lipgloss.NewStyle().Padding(1, 2, 1, 2)

	focusedStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(getColor("69"))

	unFocusedStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.HiddenBorder())
)
