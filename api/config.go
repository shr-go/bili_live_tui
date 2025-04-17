package api

type BiliLiveConfig struct {
	RoomID         uint64 `toml:"room_id"`
	ChatBuffer     int    `toml:"chat_buffer"`
	ShowShipLevel  bool   `toml:"show_ship_level"`
	ShowMedalName  bool   `toml:"show_medal_name"`
	ShowMedalLevel bool   `toml:"show_medal_level"`
	ColorMode      bool   `toml:"color_mode"`
	ShowRoomTitle  bool   `toml:"show_room_title"`
	ShowRoomNumber bool   `toml:"show_room_number"`
}
