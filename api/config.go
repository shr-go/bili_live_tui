package api

type BiliLiveConfig struct {
	RoomID     uint64 `toml:"room_id"`
	ChatBuffer int    `toml:"chat_buffer"`
}
