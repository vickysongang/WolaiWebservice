package apnsprovider

import (
	"github.com/anachronistic/apns"

	"WolaiWebservice/config"
)

const (
	APNS_GATEWAY_DEV  = "gateway.sandbox.push.apple.com:2195"
	APNS_GATEWAY_PROD = "gateway.push.apple.com:2195"

	APNS_ENV_DEV  = "develop"
	APNS_ENV_PROD = "production"
)

var apnsClient *apns.Client

func init() {
	var gateway string
	if config.Env.APNS.Env == APNS_ENV_DEV {
		gateway = APNS_GATEWAY_DEV
	} else if config.Env.APNS.Env == APNS_ENV_PROD {
		gateway = APNS_GATEWAY_PROD
	}

	apnsClient = apns.NewClient(gateway, config.Env.APNS.Cert, config.Env.APNS.Key)
}
