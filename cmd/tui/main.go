package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/shr-go/bili_live_tui/internal/live_room"
	"github.com/shr-go/bili_live_tui/internal/tui"
	"github.com/shr-go/bili_live_tui/pkg/logging"
	"os"
)

func main() {
	room, err := live_room.AuthAndConnect(545068)
	if err != nil {
		logging.Fatalf("Connect server error, err=%v", err)
	}
	p := tea.NewProgram(tui.InitialModel(room), tea.WithAltScreen(), tea.WithMouseCellMotion())
	go tui.ReceiveMsg(p, room)
	go tui.PoolWindowSize(p)
	if err := p.Start(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
