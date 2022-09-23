package live_room

import (
	"github.com/shr-go/bili_live_tui/api"
	"net/http"
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
	userRoomInfo, err := GetUserRoomInfo(client, realRoomID)
	if err != nil {
		return
	}
	room.Title = roomInfo.Data.Title
	room.ShortID = uint64(roomInfo.Data.ShortId)
	room.RoomUserInfo = userRoomInfo.Data.Property
	return
}
