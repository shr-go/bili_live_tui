package live_room

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"container/list"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/andybalholm/brotli"
	"github.com/google/go-querystring/query"
	"github.com/shr-go/bili_live_tui/api"
	"github.com/shr-go/bili_live_tui/pkg/logging"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"reflect"
	"strconv"
	"sync/atomic"
	"time"
)

var (
	invalidMessageErr    = errors.New("message invalid")
	headerNotCompleteErr = errors.New("header not complete")
)

func packMessage(data []byte, protoVer api.DanmuProtol, protoOp api.DanmuOp, seq uint32) []byte {
	b := bytes.Buffer{}
	length := uint32(16 + len(data))
	binary.Write(&b, binary.BigEndian, length)
	binary.Write(&b, binary.BigEndian, uint16(16))
	binary.Write(&b, binary.BigEndian, uint16(protoVer))
	binary.Write(&b, binary.BigEndian, uint32(protoOp))
	binary.Write(&b, binary.BigEndian, seq)

	if len(data) > 0 {
		switch protoVer {
		case api.DanmuProtolNormalZlib:
			w := zlib.NewWriter(&b)
			w.Write(data)
			w.Close()
		case api.DanmuProtolNormalBrotli:
			w := brotli.NewWriter(&b)
			w.Write(data)
			w.Close()
		default:
			b.Write(data)
		}
	}
	return b.Bytes()
}

func parseHeader(data []byte) (header *api.DanmuMessageHeader, err error) {
	if len(data) < 16 {
		err = headerNotCompleteErr
		return
	}
	b := bytes.NewBuffer(data)
	header = new(api.DanmuMessageHeader)

	v := reflect.ValueOf(header).Elem()
	for i := 0; i < v.NumField(); i++ {
		ptr := v.Field(i).Addr().Interface()
		binary.Read(b, binary.BigEndian, ptr)
	}
	if header.HeaderSize != 16 ||
		header.ProtoVer > api.DanmuProtolNormalBrotli ||
		header.OpCode > api.DanmuOpAuthResp {
		err = invalidMessageErr
		return
	}
	return
}

func unpackMessage(room *api.LiveRoom, data []byte) (unpack uint32) {
	for dataLen := len(data); dataLen > 0; dataLen = len(data) {
		header, _ := parseHeader(data)
		if header == nil || (header.Size) > uint32(dataLen) {
			return
		}
		unpack += header.Size
		rawMessage := data[header.HeaderSize:header.Size]
		data = data[header.Size:]
		logging.Debugf("read message, header=%+v", header)
		var normalMessage []byte
		switch header.ProtoVer {
		case api.DanmuProtolNormalZlib:
			b := bytes.NewBuffer(rawMessage)
			if zr, err := gzip.NewReader(b); err != nil {
				logging.Errorf("decompress gzip error, err=%v", err)
				continue
			} else {
				if normalMessage, err = ioutil.ReadAll(zr); err != nil {
					logging.Errorf("decompress gzip error, err=%v", err)
					continue
				}
			}
		case api.DanmuProtolNormalBrotli:
			b := bytes.NewBuffer(rawMessage)
			br := brotli.NewReader(b)
			var err error
			if normalMessage, err = ioutil.ReadAll(br); err != nil {
				logging.Errorf("decompress brotli error, err=%v", err)
				continue
			}
		default:
			normalMessage = rawMessage
		}
		logging.Debugf("read message, header=%+v", header)
		switch header.OpCode {
		case api.DanmuOpHeartBeatResp:
			if len(normalMessage) >= 4 {
				room.Hot = binary.BigEndian.Uint32(normalMessage)
			}
		case api.DanmuOpNormal:
			if header.ProtoVer == api.DanmuProtolNormal {
				danmuMessage := new(api.DanmuMessage)
				if err := json.Unmarshal(normalMessage, danmuMessage); err != nil {
					logging.Errorf("unmarshal normal message error, err=%v", err)
					continue
				}
				room.MessageChan <- danmuMessage
			} else if header.ProtoVer == api.DanmuProtolNormalZlib || header.ProtoVer == api.DanmuProtolNormalBrotli {
				for messagesLen := len(normalMessage); messagesLen > 0; messagesLen = len(normalMessage) {
					messageHeader, err := parseHeader(normalMessage)
					if err != nil {
						logging.Errorf("parse message error, err=%v", err)
						break
					} else if messageHeader.Size > uint32(messagesLen) {
						logging.Errorf("header message size overflow")
						break
					}
					oneNormalMessage := normalMessage[messageHeader.HeaderSize:messageHeader.Size]
					normalMessage = normalMessage[messageHeader.Size:]
					danmuMessage := new(api.DanmuMessage)
					if err := json.Unmarshal(oneNormalMessage, danmuMessage); err != nil {
						logging.Errorf("unmarshal normal message error, err=%v", err)
						continue
					}
					room.MessageChan <- danmuMessage
				}
			}
		}
	}
	return
}

func GetDanmuInfo(client *http.Client, id uint64) (info *api.DanmuInfoResp, err error) {
	danmuInfoReq := api.DanmuInfoReq{ID: id}
	v, err := query.Values(danmuInfoReq)
	if err != nil {
		return
	}
	baseURL := "https://api.live.bilibili.com/xlive/web-room/v1/index/getDanmuInfo"
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
	danmuInfoResp := new(api.DanmuInfoResp)
	err = json.Unmarshal(body, danmuInfoResp)
	if err != nil {
		return
	}
	if danmuInfoResp.Code != 0 {
		err = errors.New(danmuInfoResp.Message)
	}
	info = danmuInfoResp
	return
}

func connectDanmuServer(uid uint64, roomID uint64, info *api.DanmuInfoResp) (conn net.Conn, err error) {
	for _, HostData := range info.Data.HostList {
		timeout := time.Second
		conn, err = net.DialTimeout("tcp", net.JoinHostPort(HostData.Host, strconv.Itoa(HostData.Port)), timeout)
		if err == nil && conn != nil {
			break
		}
	}
	if conn == nil {
		return nil, errors.New("no server can connect")
	}
	danmuAuthPacketReq := api.DanmuAuthPacketReq{
		UID:      uid,
		RoomID:   roomID,
		ProtoVer: 3,
		Platform: "web",
		Type:     2,
		Key:      info.Data.Token,
	}
	jsonReq, err := json.Marshal(danmuAuthPacketReq)
	if err != nil {
		return
	}
	data := packMessage(jsonReq, api.DanmuProtolHeartBeat, api.DanmuOpAuth, 1)
	dataLen := len(data)
	n, err := conn.Write(data)
	if err != nil {
		return
	} else if n != dataLen {
		err = errors.New("connect server failed")
		return
	}
	resp := make([]byte, 8192)
	n, err = conn.Read(resp)
	if err != nil {
		return
	}
	danmuHeader, _ := parseHeader(resp)
	if danmuHeader == nil {
		err = errors.New("parse header failed")
		return
	}
	danmuAuthPacketResp := api.DanmuAuthPacketResp{}
	err = json.Unmarshal(resp[danmuHeader.HeaderSize:n], &danmuAuthPacketResp)
	if err != nil || danmuAuthPacketResp.Code != 0 {
		err = errors.New("connect server auth failed")
		return
	}
	return
}

func ConnectDanmuServer(uid uint64, roomID uint64, info *api.DanmuInfoResp) (room *api.LiveRoom, err error) {
	conn, err := connectDanmuServer(uid, roomID, info)
	if err != nil {
		return
	}
	room = &api.LiveRoom{
		UID:         uid,
		RoomID:      roomID,
		Hot:         0,
		Seq:         1,
		MessageChan: make(chan *api.DanmuMessage, 10),
		ReqChan:     make(chan []byte, 10),
		DoneChan:    make(chan struct{}),
		RetryChan:   make(chan struct{}),
		StreamConn:  conn,
	}

	// process write
	go processWrite(room)

	// process read
	go processRead(room)

	go monitorConn(room)

	return
}

func heartBeatReq(room *api.LiveRoom) {
	body, _ := hex.DecodeString("5b6f626a656374204f626a6563745d")
	seq := atomic.AddUint32(&room.Seq, 1)
	data := packMessage(body, api.DanmuProtolHeartBeat, api.DanmuOpHeartBeat, seq)
	room.ReqChan <- data
}

func processWrite(room *api.LiveRoom) {
	heartBeatReq(room)
	heartBeatTicker := time.NewTicker(30 * time.Second)
	defer heartBeatTicker.Stop()
	dataList := list.New()
	doneChan := room.DoneChan
Loop:
	for {
		select {
		case <-doneChan:
			break Loop
		case <-heartBeatTicker.C:
			heartBeatReq(room)
		case data := <-room.ReqChan:
			for dataList.Len() > 0 {
				preData := dataList.Front().Value.([]byte)
				// todo Add timeout settings
				if _, err := room.StreamConn.Write(preData); err != nil {
					if err != nil {
						logging.Errorf("connection close from write, err=%v", err)
						break Loop
					}
				} else {
					dataList.Remove(dataList.Front())
				}
			}
			if dataList.Len() == 0 {
				_, err := room.StreamConn.Write(data)
				if err == nil {
					continue
				}
			}
			dataList.PushBack(data)
		}
	}
	logging.Infof("write goroutine quit")
}

func processRead(room *api.LiveRoom) {
	var notComplete []byte
	doneChan := room.DoneChan
Loop:
	for {
		select {
		case <-doneChan:
			break Loop
		default:
			data := make([]byte, 64*1024)
			// todo Add timeout settings
			n, err := room.StreamConn.Read(data)
			if err != nil {
				close(room.RetryChan)
				logging.Errorf("connection close from read, err=%v", err)
				break Loop
			}
			data = data[:n]
			if len(notComplete) != 0 {
				data = append(notComplete, data...)
			}
			dataLen := len(data)
			if dataLen > 0 {
				unpackLen := unpackMessage(room, data)
				leftLen := dataLen - int(unpackLen)
				if leftLen > 0 {
					notComplete = data[unpackLen:]
				} else {
					notComplete = nil
				}
			}
		}
	}
	logging.Infof("read goroutine quit")
}

func monitorConn(room *api.LiveRoom) {
Loop:
	for {
		select {
		case <-room.DoneChan:
			break Loop
		case <-room.RetryChan:
			logging.Infof("retry connect danmu server")
			close(room.DoneChan)
			client := room.Client
			realRoomID := room.RoomID
			info, err := GetDanmuInfo(client, realRoomID)
			if err != nil {
				logging.Fatalf("retry get danmu info failed, err=%v", err)
			}
			conn, err := connectDanmuServer(room.UID, realRoomID, info)
			if err != nil {
				logging.Fatalf("retry connect danmu server failed, err=%v", err)
			}
			logging.Infof("retry connect danmu server success")
			room.StreamConn = conn
			room.DoneChan = make(chan struct{})
			room.RetryChan = make(chan struct{})
			go processWrite(room)
			go processRead(room)
		}
	}
	logging.Infof("monitor goroutine quit")
}
