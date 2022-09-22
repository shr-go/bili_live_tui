package live_room

import (
	"github.com/shr-go/bili_live_tui/api"
	"net/http"
)

func AuthAndConnect(roomID uint64) (room *api.LiveRoom, err error) {
	//fixme set user id
	uid := uint64(0)
	client := &http.Client{}
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
	return
}
