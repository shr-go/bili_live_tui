package api

import (
	"net"
	"net/http"
)

type LiveRoom struct {
	UID          uint64
	RoomID       uint64
	Hot          uint32
	Seq          uint32
	MessageChan  chan *DanmuMessage
	ReqChan      chan []byte
	DoneChan     chan struct{}
	StreamConn   net.Conn
	Title        string
	ShortID      uint64
	RoomUserInfo UserRoomProperty
	Client       *http.Client
	CSRF         string
}

type DanmuInfoReq struct {
	ID uint64 `url:"id"`
}

type DanmuInfoResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Ttl     int    `json:"ttl"`
	Data    struct {
		Group            string  `json:"group"`
		BusinessId       int     `json:"business_id"`
		RefreshRowFactor float64 `json:"refresh_row_factor"`
		RefreshRate      int     `json:"refresh_rate"`
		MaxDelay         int     `json:"max_delay"`
		Token            string  `json:"token"`
		HostList         []struct {
			Host    string `json:"host"`
			Port    int    `json:"port"`
			WssPort int    `json:"wss_port"`
			WsPort  int    `json:"ws_port"`
		} `json:"host_list"`
	} `json:"data"`
}

type RoomInfoReq struct {
	RoomID uint64 `url:"room_id"`
}

type RoomInfoResp struct {
	Code    int    `json:"code"`
	Msg     string `json:"msg"`
	Message string `json:"message"`
	Data    struct {
		Uid              int      `json:"uid"`
		RoomId           int      `json:"room_id"`
		ShortId          int      `json:"short_id"`
		Attention        int      `json:"attention"`
		Online           int      `json:"online"`
		IsPortrait       bool     `json:"is_portrait"`
		Description      string   `json:"description"`
		LiveStatus       int      `json:"live_status"`
		AreaId           int      `json:"area_id"`
		ParentAreaId     int      `json:"parent_area_id"`
		ParentAreaName   string   `json:"parent_area_name"`
		OldAreaId        int      `json:"old_area_id"`
		Background       string   `json:"background"`
		Title            string   `json:"title"`
		UserCover        string   `json:"user_cover"`
		Keyframe         string   `json:"keyframe"`
		IsStrictRoom     bool     `json:"is_strict_room"`
		LiveTime         string   `json:"live_time"`
		Tags             string   `json:"tags"`
		IsAnchor         int      `json:"is_anchor"`
		RoomSilentType   string   `json:"room_silent_type"`
		RoomSilentLevel  int      `json:"room_silent_level"`
		RoomSilentSecond int      `json:"room_silent_second"`
		AreaName         string   `json:"area_name"`
		Pendants         string   `json:"pendants"`
		AreaPendants     string   `json:"area_pendants"`
		HotWords         []string `json:"hot_words"`
		HotWordsStatus   int      `json:"hot_words_status"`
		Verify           string   `json:"verify"`
		NewPendants      struct {
			Frame struct {
				Name       string `json:"name"`
				Value      string `json:"value"`
				Position   int    `json:"position"`
				Desc       string `json:"desc"`
				Area       int    `json:"area"`
				AreaOld    int    `json:"area_old"`
				BgColor    string `json:"bg_color"`
				BgPic      string `json:"bg_pic"`
				UseOldArea bool   `json:"use_old_area"`
			} `json:"frame"`
			Badge struct {
				Name     string `json:"name"`
				Position int    `json:"position"`
				Value    string `json:"value"`
				Desc     string `json:"desc"`
			} `json:"badge"`
			MobileFrame struct {
				Name       string `json:"name"`
				Value      string `json:"value"`
				Position   int    `json:"position"`
				Desc       string `json:"desc"`
				Area       int    `json:"area"`
				AreaOld    int    `json:"area_old"`
				BgColor    string `json:"bg_color"`
				BgPic      string `json:"bg_pic"`
				UseOldArea bool   `json:"use_old_area"`
			} `json:"mobile_frame"`
			MobileBadge interface{} `json:"mobile_badge"`
		} `json:"new_pendants"`
		UpSession            string `json:"up_session"`
		PkStatus             int    `json:"pk_status"`
		PkId                 int    `json:"pk_id"`
		BattleId             int    `json:"battle_id"`
		AllowChangeAreaTime  int    `json:"allow_change_area_time"`
		AllowUploadCoverTime int    `json:"allow_upload_cover_time"`
		StudioInfo           struct {
			Status     int           `json:"status"`
			MasterList []interface{} `json:"master_list"`
		} `json:"studio_info"`
	} `json:"data"`
}

type DanmuAuthPacketReq struct {
	UID      uint64 `json:"uid"`
	RoomID   uint64 `json:"roomid"`
	ProtoVer uint8  `json:"protover"`
	Platform string `json:"platform"`
	Type     uint8  `json:"type"`
	Key      string `json:"key"`
}

type DanmuProtol uint16

const (
	DanmuProtolNormal DanmuProtol = iota
	DanmuProtolHeartBeat
	DanmuProtolNormalZlib
	DanmuProtolNormalBrotli
)

type DanmuOp uint32

const (
	DanmuOpHeartBeat     DanmuOp = 2
	DanmuOpHeartBeatResp DanmuOp = 3
	DanmuOpNormal        DanmuOp = 5
	DanmuOpAuth          DanmuOp = 7
	DanmuOpAuthResp      DanmuOp = 8
)

type DanmuAuthPacketResp struct {
	Code uint32 `json:"code"`
}

type DanmuMessageHeader struct {
	Size       uint32
	HeaderSize uint16
	ProtoVer   DanmuProtol
	OpCode     DanmuOp
	Sequence   uint32
}

type DanmuMessage struct {
	Cmd  string                 `json:"cmd"`
	Info []interface{}          `json:"info"`
	Data map[string]interface{} `json:"data"`
}

type QRCodeGenerateResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Ttl     int    `json:"ttl"`
	Data    struct {
		Url       string `json:"url"`
		QrcodeKey string `json:"qrcode_key"`
	} `json:"data"`
}

type QRLoginStatus uint32

const (
	QRLoginSuccess    QRLoginStatus = 0
	QRLoginNotConfirm QRLoginStatus = 86090
	QRLoginNotScan    QRLoginStatus = 86101
	QRLoginExpired    QRLoginStatus = 86038
)

type QRCodeLoginData struct {
	QRString string
	QRKey    string
	Status   QRLoginStatus
}

type PollLoginResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Ttl     int    `json:"ttl"`
	Data    struct {
		Url          string        `json:"url"`
		RefreshToken string        `json:"refresh_token"`
		Timestamp    int64         `json:"timestamp"`
		Code         QRLoginStatus `json:"code"`
		Message      string        `json:"message"`
	} `json:"data"`
}

type UserInfo struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Ttl     int    `json:"ttl"`
	Data    struct {
		Mid   uint64 `json:"mid"`
		Uname string `json:"uname"`
	} `json:"data"`
}

type UserRoomInfo struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Ttl     int    `json:"ttl"`
	Data    struct {
		Property UserRoomProperty `json:"property"`
	} `json:"data"`
}

type UserRoomProperty struct {
	UnameColor string `json:"uname_color"`
	Bubble     int    `json:"bubble"`
	Danmu      struct {
		Mode   int `json:"mode"`
		Color  int `json:"color"`
		Length int `json:"length"`
		RoomId int `json:"room_id"`
	} `json:"danmu"`
	BubbleColor string `json:"bubble_color"`
}

type SendMsgReq struct {
	Bubble    int    `url:"bubble"`
	Msg       string `url:"msg"`
	Color     int    `url:"color"`
	Mode      int    `url:"mode"`
	Fontsize  int    `url:"fontsize"`
	Rnd       int64  `url:"rnd"`
	RoomID    uint64 `url:"roomid"`
	CSRF      string `url:"csrf"`
	CSRFToken string `url:"csrf_token"`
}
