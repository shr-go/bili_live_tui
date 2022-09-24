package tui

import (
	"github.com/BurntSushi/toml"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/shr-go/bili_live_tui/api"
	"github.com/shr-go/bili_live_tui/internal/live_room"
	"github.com/shr-go/bili_live_tui/pkg/logging"
	"golang.org/x/term"
	"net/http"
	"os"
)

var (
	windowWidth  int
	windowHeight int
	LiveConfig   api.BiliLiveConfig
)

func init() {
	logging.InitLogConfig()
	windowWidth, windowHeight, _ = term.GetSize(int(os.Stdout.Fd()))
	f := "config.toml"
	_, err := toml.DecodeFile(f, &LiveConfig)
	if err != nil {
		logging.Fatalf("load config error, err=%v", err)
	}
}

func PrepareEnterRoom(client *http.Client) (room *api.LiveRoom, err error) {
	if cookieBytes, err := os.ReadFile("COOKIE.DAT"); err == nil {
		cookies := string(cookieBytes)
		if live_room.CheckCookieValid(client, cookies) {
			return live_room.AuthAndConnect(client, LiveConfig.RoomID)
		}
	}
	loginModel := newLoginModel(client)
	p := tea.NewProgram(&loginModel, tea.WithAltScreen(), tea.WithMouseCellMotion())
	if err := p.Start(); err != nil {
		logging.Fatalf("PrepareEnterRoom ui error: %v", err)
		os.Exit(1)
	}
	if loginModel.quit {
		os.Exit(0)
	}
	if !loginModel.chooseLogin {
		return live_room.AuthAndConnect(client, LiveConfig.RoomID)
	}
	if !live_room.CheckCookieValid(client, loginModel.cookies) {
		logging.Fatalf("PrepareEnterRoom cookies check failed, program exit")
		os.Exit(1)
	}
	os.WriteFile("COOKIE.DAT", []byte(loginModel.cookies), 0o660)
	return live_room.AuthAndConnect(client, LiveConfig.RoomID)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
