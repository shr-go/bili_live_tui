package tui

import (
	"container/list"
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/shr-go/bili_live_tui/api"
	"golang.org/x/term"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type sessionState uint

const (
	contentView sessionState = iota
	inputView
)

type medalInfo struct {
	level      uint8
	shipLevel  uint8
	name       string
	medalColor string
}

type danmuMsg struct {
	uid          uint64
	uName        string
	chatTime     time.Time
	content      string
	medal        *medalInfo
	nameColor    string
	contentColor string
}

type model struct {
	danmu      *list.List
	room       *api.LiveRoom
	viewport   viewport.Model
	ready      bool
	lockBottom bool
	state      sessionState
}

func InitialModel(room *api.LiveRoom) model {
	return model{
		danmu:      list.New(),
		room:       room,
		viewport:   viewport.Model{},
		ready:      false,
		lockBottom: true,
		state:      contentView,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		headerHeight := lipgloss.Height(m.headerView())
		footerHeight := lipgloss.Height(m.footerView())
		verticalMarginHeight := headerHeight + footerHeight

		if !m.ready {
			m.viewport = viewport.New(msg.Width, msg.Height-verticalMarginHeight)
			m.viewport.YPosition = headerHeight
			m.viewport.HighPerformanceRendering = false
			m.viewport.SetContent(m.renderDanmu())
			m.ready = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - verticalMarginHeight
		}
	case *danmuMsg:
		m.danmu.PushBack(msg)
		//fixme use config to replace hard code
		for m.danmu.Len() > 200 {
			m.danmu.Remove(m.danmu.Front())
		}
		if m.ready {
			m.viewport.SetContent(m.renderDanmu())
		}
	}

	if m.lockBottom {
		m.viewport.GotoBottom()
	}

	scrollPercent := m.viewport.ScrollPercent()
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)
	newScrollPercent := m.viewport.ScrollPercent()

	if scrollPercent != newScrollPercent {
		m.lockBottom = newScrollPercent == 1
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if !m.ready {
		return "\nInitializing..."
	}
	return fmt.Sprintf("%s\n%s\n%s", m.headerView(), m.viewport.View(), m.footerView())
}

func ReceiveMsg(program *tea.Program, room *api.LiveRoom) {
	for msg := range room.MessageChan {
		if msg.Cmd == "DANMU_MSG" {
			if danmu := processDanmuMsg(msg); danmu != nil {
				program.Send(danmu)
			}
		}
	}
}

func PoolWindowSize(program *tea.Program) {
	if runtime.GOOS != "windows" {
		return
	}
	width, height, _ := term.GetSize(int(os.Stdout.Fd()))
	for range time.Tick(50 * time.Millisecond) {
		nowWidth, nowHeight, _ := term.GetSize(int(os.Stdout.Fd()))
		if width != nowWidth || height != nowHeight {
			width = nowWidth
			height = nowHeight
			windowSize := tea.WindowSizeMsg{
				Width:  width,
				Height: height,
			}
			program.Send(windowSize)
		}
	}

}

func (m model) headerView() string {
	b := lipgloss.RoundedBorder()
	b.Right = "├"
	roomID := m.room.ShortID
	if roomID == 0 {
		roomID = m.room.RoomID
	}
	header := fmt.Sprintf("%s - %d(%d)", m.room.Title, roomID, m.room.Hot)
	title := lipgloss.NewStyle().BorderStyle(b).Padding(0, 1).
		Render(header)
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (m model) footerView() string {
	b := lipgloss.RoundedBorder()
	b.Left = "┤"
	info := titleStyle.Copy().BorderStyle(b).
		Render(fmt.Sprintf("%3.f%%", m.viewport.ScrollPercent()*100))
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(info)))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
}

func (m model) renderDanmu() string {
	sb := strings.Builder{}
	viewportHeight := m.viewport.Height
	for n := m.danmu.Len(); n < viewportHeight; n++ {
		sb.WriteRune('\n')
	}
	for danmuElem := m.danmu.Front(); danmuElem != nil; danmuElem = danmuElem.Next() {
		danmu, ok := danmuElem.Value.(*danmuMsg)
		if ok {
			if danmu.medal != nil {
				sb.WriteString(medalStyle(danmu.medal))
				sb.WriteRune(' ')
			}
			sb.WriteString(fmt.Sprintln(nameStyle(danmu.uName, danmu.nameColor),
				contentStyle(danmu.content, danmu.contentColor)))
		}
	}
	return sb.String()
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
