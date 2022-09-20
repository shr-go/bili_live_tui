package tui

import (
	"container/list"
	"time"
)

type chat struct {
	uid      uint64
	chatTime time.Time
	content  string ``
}

type userInfo struct {
}

type model struct {
	chatContent *list.List
}
