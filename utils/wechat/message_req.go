// POIWeChatMessage
package wechat

import (
	//"bytes"
	//	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	//"WolaiWebservice/config"

	"github.com/cihub/seelog"
)

const (
	WC_SEND_MSG = "http://wx.wolai.me/weixin/open/summary/message/template/send"
)

type WCField struct {
	Value string `json:"value"`
}

type WCMessage struct {
	ToUser  string `json:"touser"`
	Data    string `json:"data"`
	Url     string `json:"url"`
	MsgType string `json:"type"`
}

func WCSendTypedMessage(wcMsg *WCMessage) error {
	return wcSendMessage(wcMsg)
}

func wcSendMessage(wcMsg *WCMessage) error {
	var err error

	send_url := WC_SEND_MSG

	//query, _ := json.Marshal(wcMsg)
	//seelog.Trace("[wcSendMessage]: ", string(query))

	//req, err := http.NewRequest("POST", url, bytes.NewBuffer(query))

	vv := url.Values{}

	vv.Set("data", wcMsg.Data)
	vv.Set("url", wcMsg.Url)
	vv.Set("type", wcMsg.MsgType)
	vv.Set("touser", wcMsg.ToUser)

	body := ioutil.NopCloser(strings.NewReader(vv.Encode()))

	req, err := http.NewRequest("POST", send_url, body)

	if err != nil {
		seelog.Error(err.Error())
		return errors.New("创建消息请求失败")
	}
	//	req.Header.Set("X-AVOSCloud-Application-Id", config.Env.LeanCloud.AppId)
	//	req.Header.Set("X-AVOSCloud-Master-Key", config.Env.LeanCloud.MasterKey)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		seelog.Error(err.Error())
		return errors.New("发送消息请求失败")
	}

	b, err := ioutil.ReadAll(resp.Body)
	seelog.Debugf("[wcSendMessage] response: %s", b)

	defer func() {
		if x := recover(); x != nil {
			seelog.Error(x)
		}
	}()
	defer resp.Body.Close()
	return nil
}
