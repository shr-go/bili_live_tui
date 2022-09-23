package live_room

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/shr-go/bili_live_tui/api"
	"github.com/skip2/go-qrcode"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)

var (
	QRCodeGenerateErr = errors.New("QRCode generate error")
	PollLoginError    = errors.New("poll login failed")
)

func QRCodeLogin(client *http.Client) (data *api.QRCodeLoginData, err error) {
	baseURL := "https://passport.bilibili.com/x/passport-login/web/qrcode/generate"
	resp, err := client.Get(baseURL)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	respData := new(api.QRCodeGenerateResp)
	if err = json.Unmarshal(body, respData); err != nil || respData.Code != 0 {
		err = QRCodeGenerateErr
		return
	}
	q, err := qrcode.New(respData.Data.Url, qrcode.Low)
	if err != nil {
		return
	}
	qrStr := q.ToSmallString(false)
	data = &api.QRCodeLoginData{
		QRString: qrStr,
		QRKey:    respData.Data.QrcodeKey,
		Status:   api.QRLoginNotScan,
	}
	return
}

func PollLogin(client *http.Client, data *api.QRCodeLoginData) (cookie string, err error) {
	baseURL := "https://passport.bilibili.com/x/passport-login/web/qrcode/poll"
	realURL := fmt.Sprintf("%s?qrcode_key=%s", baseURL, data.QRKey)
	resp, err := client.Get(realURL)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	var pollLogin api.PollLoginResp
	err = json.Unmarshal(body, &pollLogin)
	if err != nil || pollLogin.Code != 0 {
		err = PollLoginError
		return
	}
	data.Status = pollLogin.Data.Code
	if data.Status == api.QRLoginSuccess {
		sb := strings.Builder{}
		cookies := resp.Cookies()
		for n, oneCookie := range cookies {
			sb.WriteString(oneCookie.Name)
			sb.WriteRune('=')
			sb.WriteString(oneCookie.Value)
			if n+1 != len(cookies) {
				sb.WriteString("; ")
			}
		}
		cookie = sb.String()
	}
	return
}

func parseCookieStr(client *http.Client, cookies string) {
	defer func() {
		recover()
	}()
	jar, _ := cookiejar.New(nil)
	elements := strings.Split(cookies, ";")
	var cookieSlice []*http.Cookie
	for _, element := range elements {
		element := strings.TrimSpace(element)
		nameValue := strings.Split(element, "=")
		cookie := &http.Cookie{
			Name:   nameValue[0],
			Value:  nameValue[1],
			Path:   "/",
			Domain: ".bilibili.com",
		}
		cookieSlice = append(cookieSlice, cookie)
	}
	u, _ := url.Parse("https://bilibili.com")
	jar.SetCookies(u, cookieSlice)
	client.Jar = jar
}

func CheckCookieValid(client *http.Client, cookie string) bool {
	parseCookieStr(client, cookie)
	baseURL := "https://account.bilibili.com/site/getCoin"
	resp, err := client.Get(baseURL)
	if err != nil {
		return false
	}

	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false
	}
	var data map[string]interface{}
	if err = json.Unmarshal(respBody, &data); err != nil {
		return false
	}
	return data["code"].(float64) == 0
}

func GetUserInfo(client *http.Client) *api.UserInfo {
	baseURL := "https://api.bilibili.com/x/web-interface/nav"
	resp, err := client.Get(baseURL)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil
	}
	var userInfo api.UserInfo
	if err = json.Unmarshal(respBody, &userInfo); err != nil || userInfo.Code != 0 {
		return nil
	}
	return &userInfo
}

func getCSRF(client *http.Client) string {
	u, _ := url.Parse("https://bilibili.com")
	cookies := client.Jar.Cookies(u)
	for _, cookie := range cookies {
		if cookie.Name == "bili_jct" {
			return cookie.Value
		}
	}
	return ""
}
