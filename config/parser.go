package config

import (
	"encoding/json"
	"time"

	"github.com/cihub/seelog"
)

// Parse duration
type duration struct {
	time.Duration
}

type UpgradeInfo struct {
	MaxVersion      int64  `json:"max_version"`
	ForceMinVersion int64  `json:"force_min_version"`
	DownloadUrl     string `json:"download_url"`
}

func (d *duration) UnmarshalText(text []byte) error {
	var err error
	d.Duration, err = time.ParseDuration(string(text))
	if err != nil {
		seelog.Error(string(text), " ", err.Error())
		return err
	}
	return nil
}

func (info *UpgradeInfo) UnmarshalConfig(config []byte) error {
	var err error
	err = json.Unmarshal(config, &info)
	return err
}
