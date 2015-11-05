// POISendcloudController
package sendcloud

import (
	"POIWolaiWebService/utils"
	"bytes"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"time"

	seelog "github.com/cihub/seelog"
)

const (
	SC_SMS_URL = "http://sendcloud.sohu.com/smsapi/send"
)

type valSorter struct {
	Keys []string
	Vals []string
}

func mapSorter(m map[string]string) *valSorter {
	vs := &valSorter{
		Keys: make([]string, 0, len(m)),
		Vals: make([]string, 0, len(m)),
	}
	for k, v := range m {
		vs.Keys = append(vs.Keys, k)
		vs.Vals = append(vs.Vals, v)
	}
	return vs
}

func (vs *valSorter) Sort() {
	sort.Sort(vs)
}

func (vs *valSorter) Len() int           { return len(vs.Keys) }
func (vs *valSorter) Less(i, j int) bool { return vs.Keys[i] < vs.Keys[j] }
func (vs *valSorter) Swap(i, j int) {
	vs.Vals[i], vs.Vals[j] = vs.Vals[j], vs.Vals[i]
	vs.Keys[i], vs.Keys[j] = vs.Keys[j], vs.Keys[i]
}

func Signature(smsKey string, params url.Values) (result string) {
	var query string
	pa := make(map[string]string)
	for k, v := range params {
		pa[k] = v[0]
	}
	vs := mapSorter(pa)
	vs.Sort()
	for i := 0; i < vs.Len(); i++ {
		if vs.Keys[i] == "signature" {
			continue
		}
		if vs.Keys[i] != "" && vs.Vals[i] != "" {
			query = fmt.Sprintf("%v&%v=%v", query, vs.Keys[i], vs.Vals[i])
		}
	}
	string_to_sign := fmt.Sprintf("%v%v&%v", smsKey, query, smsKey)
	md5New := md5.New()
	md5New.Write([]byte(string_to_sign))
	return hex.EncodeToString(md5New.Sum(nil))
}

func SCSendMessage(phone string, randCode string) error {
	params := url.Values{
		"smsUser":    {utils.Config.SendCloud.SmsUser},
		"templateId": {utils.Config.SendCloud.TemplateId},
		"phone":      {phone},
		"vars":       {"{'%Code%':'" + randCode + "'}"},
	}
	encodeParams := Signature(utils.Config.SendCloud.SmsKey, params)
	params.Add("signature", encodeParams)
	postBoby := bytes.NewBufferString(params.Encode())
	resp, err := http.Post(SC_SMS_URL, "application/x-www-form-urlencoded", postBoby)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	bodyByte, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	seelog.Info("send message to ", phone, " ", string(bodyByte))
	return nil
}

func GenerateRandCode() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	rs := strconv.Itoa(r.Int())
	randCode := rs[0:4]
	return randCode
}

func verify(appKey, token, timestamp, signature string) bool {
	sha := sha256.New()
	sha.Write([]byte(appKey))
	result := sha.Sum([]byte(timestamp + token))
	signatureCal := hex.EncodeToString(result)
	fmt.Println("appkey:", appKey, " token:", token, " timestamp:", timestamp, " signature:", signature)
	fmt.Println("signatureCal:", signatureCal)
	return signature == signatureCal
}

func SMSHook(token, timestamp, signature, message string, phones []string) {
	verify(utils.Config.SendCloud.AppKey, token, timestamp, signature)
	fmt.Println("message:", message)
}

func SendMessage(phone string) error {
	randCode := GenerateRandCode()
	err := SCSendMessage(phone, randCode)
	return err
}
