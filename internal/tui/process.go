package tui

import (
	"bytes"
	"fmt"
	"github.com/shr-go/bili_live_tui/api"
	"mime/multipart"
	"reflect"
	"strconv"
	"time"
)

func processDanmuMsg(msg *api.DanmuMessage) (danmu *danmuMsg) {
	defer func() {
		if r := recover(); r != nil {
			danmu = nil
		}
	}()
	rawBasicInfo := msg.Info[0].([]interface{})
	content := msg.Info[1].(string)
	rawUserInfo := msg.Info[2].([]interface{})
	rawMedalInfo := msg.Info[3].([]interface{})

	var medal *medalInfo
	if len(rawMedalInfo) > 10 {
		medal = new(medalInfo)
		medal.level = uint8(rawMedalInfo[0].(float64))
		medal.shipLevel = uint8(rawMedalInfo[10].(float64))
		medal.name = rawMedalInfo[1].(string)
		medal.medalColor = fmt.Sprintf("#%06X", int64(rawMedalInfo[4].(float64)))
	}
	danmu = &danmuMsg{
		uid:          uint64(rawUserInfo[0].(float64)),
		uName:        rawUserInfo[1].(string),
		chatTime:     time.UnixMilli(int64(rawBasicInfo[4].(float64))),
		content:      content,
		medal:        medal,
		nameColor:    rawUserInfo[7].(string),
		contentColor: fmt.Sprintf("#%06X", int64(rawBasicInfo[3].(float64))),
	}
	return
}

func generateFakeDanmuMsg(content string) (danmu *danmuMsg) {
	danmu = &danmuMsg{
		uid:          10000,
		uName:        "【未登录 这是一条假弹幕】",
		chatTime:     time.Now(),
		content:      content,
		medal:        nil,
		nameColor:    "#DC143C",
		contentColor: "#DC143C",
	}
	return danmu
}

func generateDanmuMsg(content string, room *api.LiveRoom) (danmu *api.SendMsgReq) {
	property := room.RoomUserInfo
	return &api.SendMsgReq{
		Bubble:    property.Bubble,
		Msg:       content,
		Color:     property.Danmu.Color,
		Mode:      property.Danmu.Mode,
		Fontsize:  25,
		Rnd:       time.Now().Unix(),
		RoomID:    room.RoomID,
		CSRF:      room.CSRF,
		CSRFToken: room.CSRF,
	}
}

func packDanmuMsgForm(danmu *api.SendMsgReq) (contentType string, form *bytes.Buffer) {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
	v := reflect.ValueOf(danmu).Elem()
	t := reflect.TypeOf(danmu).Elem()
	for i := 0; i < v.NumField(); i++ {
		key := t.Field(i).Tag.Get("url")
		vi := v.Field(i).Interface()
		switch value := vi.(type) {
		case int:
			bodyWriter.WriteField(key, strconv.Itoa(value))
		case int64:
			bodyWriter.WriteField(key, strconv.FormatInt(value, 10))
		case uint64:
			bodyWriter.WriteField(key, strconv.FormatUint(value, 10))
		case string:
			bodyWriter.WriteField(key, value)
		}
	}
	bodyWriter.Close()
	contentType = bodyWriter.FormDataContentType()
	form = bodyBuf
	return
}
