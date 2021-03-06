package sendcloud

import (
	"bytes"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"WolaiWebservice/config"
	"WolaiWebservice/redis"
)

const (
	SC_SMS_URL              = "http://sendcloud.sohu.com/smsapi/send"
	SC_SMSHOOK_REQUEST      = "request"
	SC_SMSHOOK_DELIVER      = "deliver"
	SC_SMSHOOK_WORKERERROR  = "workererror"
	SC_SMSHOOK_DELIVEREROOR = "delivererror"
)

var (
	ErrMsgRepeatSend   = errors.New("验证码发送过于频繁")
	ErrRandCodeTimeout = errors.New("验证码已失效")
)

type SendcloudResponse struct {
	Message    string `json:"message"`
	Result     bool   `json:"result"`
	StatusCode int64  `json:"statusCode"`
}

type RandInfoCode struct {
	RandCode  string
	Timestamp int64
}

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
		"smsUser":    {config.Env.SendCloud.SmsUser},
		"templateId": {config.Env.SendCloud.TemplateId},
		"phone":      {phone},
		"vars":       {"{'%Code%':'" + randCode + "'}"},
	}
	encodeParams := Signature(config.Env.SendCloud.SmsKey, params)
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
	var sendcloudResponse SendcloudResponse
	json.Unmarshal(bodyByte, &sendcloudResponse)
	if !sendcloudResponse.Result {
		return errors.New(sendcloudResponse.Message)
	}
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
	return signature == signatureCal
}

func SMSHook(token, timestamp, signature, event string, phones, randCodeType string) {
	verify(config.Env.SendCloud.AppKey, token, timestamp, signature)
	phones = phones[1 : len(phones)-1]
	if event == SC_SMSHOOK_DELIVEREROOR || event == SC_SMSHOOK_WORKERERROR {
		for _, phone := range strings.Split(phones, ",") {
			redis.RemoveSendcloudRandCode(phone, randCodeType)
		}
	}
}

func SendMessage(phone, randCodeType string) error {
	oldRandCode, timestamp := redis.GetSendcloudRandCode(phone, randCodeType)
	currTimeUnix := time.Now().Unix()
	if oldRandCode != "" && (currTimeUnix-timestamp <= 60) {
		return ErrMsgRepeatSend
	} else {
		newRandCode := GenerateRandCode()
		err := SCSendMessage(phone, newRandCode)
		redis.SetSendcloudRandCode(phone, newRandCode, randCodeType)
		return err
	}
	return nil
}
