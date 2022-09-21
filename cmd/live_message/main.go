package main

import (
	"encoding/json"
	"github.com/shr-go/bili_live_tui/api"
	"github.com/shr-go/bili_live_tui/internal/live_room"
	"github.com/shr-go/bili_live_tui/pkg/logging"
	"io/ioutil"
	"net/http"
	"time"
)

func main() {
	uid := uint64(0)
	roomID := uint64(545068)

	client := &http.Client{}
	info, err := live_room.GetDanmuInfo(client, roomID)

	if err != nil {
		logging.Errorf("GetDanmuInfo Error, %v", err)
		return
	}

	server, err := live_room.ConnectDanmuServer(uid, roomID, info)
	if err != nil {
		logging.Errorf("ConnectDanmuServer Error, %v", err)
		return
	}
	var messages []*api.DanmuMessage
	for n := 0; n < 100; n++ {
		msg := <-server.MessageChan
		messages = append(messages, msg)
		logging.Debugf("%d/%d", n, 100)
	}

	if data, err := json.Marshal(messages); err != nil {
		logging.Error(err)
	} else {
		ioutil.WriteFile("data.json", data, 0o666)
	}
	close(server.DoneChan)
	time.Sleep(10 * time.Second)

	//time.Sleep(600 * time.Second)
	//close(server.DoneChan)
	//logging.Debugf("Close DoneChan")
	//time.Sleep(10 * time.Second)
}
