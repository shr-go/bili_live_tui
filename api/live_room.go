package api

import (
	"net"
)

type LiveRoom struct {
	UID         uint64
	RoomID      uint64
	Hot         uint32
	Seq         uint32
	MessageChan chan *DanmuMessage
	ReqChan     chan []byte
	DoneChan    chan struct{}
	Client      net.Conn
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
