package tui

import (
	"github.com/BurntSushi/toml"
	"github.com/shr-go/bili_live_tui/api"
	"github.com/shr-go/bili_live_tui/pkg/logging"
	"golang.org/x/term"
	"os"
)

var (
	windowWidth  int
	windowHeight int
	LiveConfig   api.BiliLiveConfig
)

func init() {
	logging.InitLogConfig()
	windowWidth, windowHeight, _ = term.GetSize(int(os.Stdout.Fd()))
	f := "config.toml"
	_, err := toml.DecodeFile(f, &LiveConfig)
	if err != nil {
		logging.Fatalf("load config error, err=%v", err)
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
