package leancloud

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/cihub/seelog"

	"WolaiWebservice/config"
)

const LC_INSTALL_BASE = "https://api.leancloud.cn/1.1/installations/"

func LCGetIntallation(objectId string) {
	url := LC_INSTALL_BASE + objectId

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		seelog.Error("LCGetIntallation:", err.Error())
	}
	req.Header.Set("X-LC-Id", config.Env.LeanCloud.AppId)
	req.Header.Set("X-LC-Key", config.Env.LeanCloud.AppKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		seelog.Error("LCGetIntallation:", err.Error())
	}
	defer func() {
		if x := recover(); x != nil {
			seelog.Error(x)
		}
	}()
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	seelog.Trace(string(body))
	var respMap map[string]interface{}
	err = json.Unmarshal(body, &respMap)
	if err != nil {
		seelog.Error(err.Error())
	}
	seelog.Trace(respMap["objectId"])
}
