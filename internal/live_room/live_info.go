package live_room

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/go-querystring/query"
	"github.com/shr-go/bili_live_tui/api"
	"io"
	"net/http"
)

func GetRoomInfo(client *http.Client, roomID uint64) (info *api.RoomInfoResp, err error) {
	roomInfoReq := api.RoomInfoReq{RoomID: roomID}
	v, err := query.Values(roomInfoReq)
	if err != nil {
		return
	}
	baseURL := "https://api.live.bilibili.com/room/v1/Room/get_info"
	realUrl := fmt.Sprintf("%s?%s", baseURL, v.Encode())
	resp, err := client.Get(realUrl)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	roomInfoResp := new(api.RoomInfoResp)
	err = json.Unmarshal(body, roomInfoResp)
	if err != nil {
		return
	}
	if roomInfoResp.Code != 0 {
		err = errors.New(roomInfoResp.Message)
	}
	info = roomInfoResp
	return
}

//GetUserRoomInfo this function trigger user enter room event
func GetUserRoomInfo(client *http.Client, roomID uint64) (info *api.UserRoomInfo, err error) {
	roomInfoReq := api.RoomInfoReq{RoomID: roomID}
	v, err := query.Values(roomInfoReq)
	if err != nil {
		return
	}
	baseURL := "https://api.live.bilibili.com/xlive/web-room/v1/index/getInfoByUser"
	realUrl := fmt.Sprintf("%s?%s", baseURL, v.Encode())
	resp, err := client.Get(realUrl)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	userRoomInfo := new(api.UserRoomInfo)
	err = json.Unmarshal(body, userRoomInfo)
	if err != nil {
		return
	}
	if userRoomInfo.Code != 0 {
		err = errors.New(userRoomInfo.Message)
	}
	info = userRoomInfo
	return
}
