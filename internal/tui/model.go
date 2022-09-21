package tui

import (
	"container/list"
	"fmt"
	"github.com/shr-go/bili_live_tui/api"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
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
	danmu *list.List
	room  *api.LiveRoom
}

func InitialModel(room *api.LiveRoom) model {
	return model{
		danmu: list.New(),
		room:  room,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
	case *danmuMsg:
		m.danmu.PushBack(msg)
		//fixme use config to replace hard code
		for m.danmu.Len() > 10 {
			m.danmu.Remove(m.danmu.Front())
		}
	}
	return m, nil
}

func (m model) View() string {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("Chat Content - %d - %d\n\n", m.room.RoomID, m.room.Hot))
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

func ReceiveMsg(program *tea.Program, room *api.LiveRoom) {
	for msg := range room.MessageChan {
		if msg.Cmd == "DANMU_MSG" {
			if danmu := processDanmuMsg(msg); danmu != nil {
				program.Send(danmu)
			}
		}
	}
}
