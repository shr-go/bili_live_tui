package tui

import (
	"fmt"
	"github.com/shr-go/bili_live_tui/api"
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
		medal.medalColor = fmt.Sprintf("#%06X", int64(rawMedalInfo[3].(float64)))
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
		uName:        "测试弹幕机器人",
		chatTime:     time.Now(),
		content:      content,
		medal:        nil,
		nameColor:    "",
		contentColor: "",
	}
	return danmu
}
