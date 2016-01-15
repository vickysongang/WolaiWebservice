package leancloud

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/cihub/seelog"

	"WolaiWebservice/config"
)

const LC_INSTALL_BASE = "https://api.leancloud.cn/1.1/installations/"

type LCInstallation struct {
	Valid          bool     `json:"valid"`
	ObjectId       string   `json:"objectId"`
	DeviceType     string   `json:"deviceType"`
	DeviceToken    string   `json:"deviceToken"`
	DeviceProfile  string   `jsoN:"deviceProfile"`
	InstallationId string   `json:"installationId"`
	TimeZone       string   `json:"timeZone"`
	UpdateAt       string   `json:"updateAt"`
	CreateAt       string   `json:"createAt"`
	Channels       []string `json:"channels"`
	Badge          int64    `json:"badge"`
}

func LCGetIntallation(objectId string) (*LCInstallation, error) {
	url := LC_INSTALL_BASE + objectId

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		seelog.Error("LCGetIntallation:", err.Error())
		return nil, errors.New("创建网络请求失败")
	}
	req.Header.Set("X-LC-Id", config.Env.LeanCloud.AppId)
	req.Header.Set("X-LC-Key", config.Env.LeanCloud.AppKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		seelog.Error("LCGetIntallation:", err.Error())
		return nil, errors.New("网络请求失败")
	}
	defer func() {
		if x := recover(); x != nil {
			seelog.Error(x)
		}
	}()
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	seelog.Trace(string(body))

	var inst LCInstallation
	err = json.Unmarshal(body, &inst)
	if err != nil {
		seelog.Error(err.Error())
		return nil, errors.New("无效的返回值")
	}

	return &inst, nil
}
