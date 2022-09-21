package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/shr-go/bili_live_tui/internal/live_room"
	"github.com/shr-go/bili_live_tui/internal/tui"
	"github.com/shr-go/bili_live_tui/pkg/logging"
	"net/http"
	"os"
)

func main() {
	uid := uint64(0)
	roomID := uint64(545068)

	client := &http.Client{}
	info, err := live_room.GetDanmuInfo(client, roomID)

	if err != nil {
		logging.Fatalf("GetDanmuInfo Error, %v", err)
	}

	room, err := live_room.ConnectDanmuServer(uid, roomID, info)
	if err != nil {
		logging.Fatalf("ConnectDanmuServer Error, %v", err)
	}

	p := tea.NewProgram(tui.InitialModel(room))
	go tui.ReceiveMsg(p, room)
	if err := p.Start(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
