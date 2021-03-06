package apnsprovider

import (
	"errors"

	"github.com/anachronistic/apns"
	"github.com/cihub/seelog"

	"WolaiWebservice/config"
	"WolaiWebservice/models"
)

const (
	APNS_GATEWAY_DEV  = "gateway.sandbox.push.apple.com:2195"
	APNS_GATEWAY_PROD = "gateway.push.apple.com:2195"

	APNS_ENV_DEV  = "develop"
	APNS_ENV_PROD = "production"
)

var appStoreClient *apns.Client
var inHouseClient *apns.Client
var voipClient *apns.Client

func init() {
	var gateway string
	if config.Env.APNS.Env == APNS_ENV_DEV {
		gateway = APNS_GATEWAY_DEV
	} else if config.Env.APNS.Env == APNS_ENV_PROD {
		gateway = APNS_GATEWAY_PROD
	}

	appStoreClient = apns.NewClient(gateway, config.Env.APNS.AppStoreCert, config.Env.APNS.AppStoreKey)
	inHouseClient = apns.NewClient(gateway, config.Env.APNS.InHouseCert, config.Env.APNS.InHouseKey)
	voipClient = apns.NewClient(gateway, config.Env.APNS.VoipCert, config.Env.APNS.VoipKey)
}

func send(pn *apns.PushNotification, deviceProfile string) error {
	var resp *apns.PushNotificationResponse

	if deviceProfile == models.DEVICE_PROFILE_APPSTORE {
		resp = appStoreClient.Send(pn)
	} else if deviceProfile == models.DEVICE_PROFILE_VOIP {
		resp = voipClient.Send(pn)
	} else {
		resp = inHouseClient.Send(pn)
	}

	info, _ := pn.PayloadString()
	seelog.Tracef("[APNS Push] Send: %s, (Token: %s|Profile: %s)",
		info, pn.DeviceToken, deviceProfile)

	if !resp.Success {
		seelog.Errorf("[APNS Push] Error: %s, (Token: %s|Profile: %s)",
			resp.Error.Error(), pn.DeviceToken, deviceProfile)
		return errors.New("推送失败")
	}

	return nil
}
