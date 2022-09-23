package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/shr-go/bili_live_tui/internal/live_room"
	"github.com/shr-go/bili_live_tui/internal/tui"
	"github.com/shr-go/bili_live_tui/pkg/logging"
	"net/http"
	"os"
)

func main() {
	client := &http.Client{}
	if !tui.LoadCookie(client) {
		logging.Fatalf("Cookie check failed")
	}
	room, err := live_room.AuthAndConnect(client, 7570705)
	if err != nil {
		logging.Fatalf("Connect server error, err=%v", err)
	}
	p := tea.NewProgram(tui.InitialModel(room), tea.WithAltScreen(), tea.WithMouseCellMotion())
	go tui.ReceiveMsg(p, room)
	go tui.PoolWindowSize(p)
	if err := p.Start(); err != nil {
		logging.Fatalf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
