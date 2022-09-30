package live_room

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/google/go-querystring/query"
	"github.com/shr-go/bili_live_tui/api"
	"github.com/shr-go/bili_live_tui/pkg/logging"
	"io"
	"net/http"
	"time"
)

func AuthAndConnect(client *http.Client, roomID uint64) (room *api.LiveRoom, err error) {
	uid := uint64(0)
	if userInfo := GetUserInfo(client); userInfo != nil {
		uid = userInfo.Data.Mid
	}
	roomInfo, err := GetRoomInfo(client, roomID)
	if err != nil {
		return
	}
	realRoomID := uint64(roomInfo.Data.RoomId)
	info, err := GetDanmuInfo(client, realRoomID)
	if err != nil {
		return
	}
	room, err = ConnectDanmuServer(uid, realRoomID, info)
	if err != nil {
		return
	}
	room.Title = roomInfo.Data.Title
	room.ShortID = uint64(roomInfo.Data.ShortId)
	room.Client = client

	if CheckAuth(client) {
		userRoomInfo, err := GetUserRoomInfo(client, realRoomID)
		if err != nil {
			return nil, err
		}
		roomUserInfo := userRoomInfo.Data.Property
		room.RoomUserInfo = &roomUserInfo
		room.CSRF = getCSRF(client)
		// 处理心跳
		go processHeartBeat(room)
	} else {
		room.RoomUserInfo = nil
	}
	return
}

func processHeartBeat(room *api.LiveRoom) {
	nextInterval := 20
	heartBeatTicker := time.NewTicker(time.Duration(nextInterval) * time.Second)
	defer heartBeatTicker.Stop()
Loop:
	for {
		select {
		case <-room.DoneChan:
			break Loop
		case <-heartBeatTicker.C:
			newNextInterval := roomHeartBeatReq(room.Client, nextInterval, room.RoomID)
			if newNextInterval != nextInterval {
				nextInterval = newNextInterval
				heartBeatTicker.Reset(time.Duration(nextInterval) * time.Second)
			}
		}
	}
}

func roomHeartBeatReq(client *http.Client, nextInterval int, realRoomID uint64) int {
	logging.Debugf("roomHeartBeatReq, nextInterval=%d, realRoomID=%d", nextInterval, realRoomID)
	hb := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%d|%d|1|0", nextInterval, realRoomID)))
	params := struct {
		HB string `url:"hb"`
		PF string `url:"pf"`
	}{
		HB: hb,
		PF: "web",
	}
	v, err := query.Values(params)
	if err != nil {
		logging.Errorf("heart beat error, err=%v", err)
	}
	baseURL := "https://live-trace.bilibili.com/xlive/rdata-interface/v1/heartbeat/webHeartBeat"
	realUrl := fmt.Sprintf("%s?%s", baseURL, v.Encode())
	resp, err := client.Get(realUrl)
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logging.Errorf("heart beat error, err=%v", err)
	}
	var data struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Ttl     int    `json:"ttl"`
		Data    struct {
			NextInterval int `json:"next_interval"`
		} `json:"data"`
	}

	if err = json.Unmarshal(body, &data); err != nil || data.Code != 0 {
		logging.Errorf("heart beat error, err=%v, data=%v", err, data)
	}
	return data.Data.NextInterval
}
