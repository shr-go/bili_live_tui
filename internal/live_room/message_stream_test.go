package live_room

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
	"time"
)

func AssertEqual(t *testing.T, a interface{}, b interface{}) {
	if a == b {
		return
	}
	t.Errorf("Received %v (type %v), expected %v (type %v)", a, reflect.TypeOf(a), b, reflect.TypeOf(b))
}

func TestDanmuInfo(t *testing.T) {
	client := &http.Client{}
	info, err := GetDanmuInfo(client, 3)
	if err != nil {
		t.Logf("GetDanmuInfo Error, %v\n", err)
		return
	}
	t.Logf("%+v", info)
}

func TestConnect(t *testing.T) {
	client := &http.Client{}
	info, err := GetDanmuInfo(client, 3)
	if err != nil {
		t.Logf("GetDanmuInfo Error, %v\n", err)
		return
	}
	_, err = ConnectDanmuServer(0, 3, info)
	if err != nil {
		t.Logf("Connect Error, %v\n", err)
		return
	}
}

func TestConnectDanmuServer(t *testing.T) {
	uid := uint64(0)
	roomID := uint64(545068)

	client := &http.Client{}
	info, err := GetDanmuInfo(client, roomID)
	if err != nil {
		t.Logf("GetDanmuInfo Error, %v\n", err)
		return
	}

	server, err := ConnectDanmuServer(uid, roomID, info)
	if err != nil {
		t.Logf("ConnectDanmuServer Error, %v\n", err)
		return
	}
	time.Sleep(600 * time.Second)
	close(server.DoneChan)
	fmt.Println("Close DoneChan")
	time.Sleep(10 * time.Second)
}
