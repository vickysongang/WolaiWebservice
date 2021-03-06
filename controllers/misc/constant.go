package misc

import (
	"WolaiWebservice/config"
	"WolaiWebservice/config/settings"
	"WolaiWebservice/models"
	miscService "WolaiWebservice/service/misc"
	"errors"
	"fmt"

	"qiniupkg.com/api.v7/conf"
	"qiniupkg.com/api.v7/kodo"
)

type UpgradeConfig struct {
	UpgradeType    string `json:"upgradeType"`
	DownloadUrl    string `json:"downloadUrl"`
	UpgradeContent string `json:"upgradeContent"`
}

func GetGradeList(pid int64) (int64, []*models.Grade) {
	var err error
	var grades []*models.Grade
	if pid != -1 {
		grades, err = miscService.QueryGradesByPid(pid)
		if err != nil {
			return 2, nil
		}
	} else {
		grades, err = miscService.QueryAllGrades()
		if err != nil {
			return 2, nil
		}
	}

	return 0, grades
}

func GetSubjectList(gradeId int64) (int64, []*models.Subject) {
	subjects := make([]*models.Subject, 0)
	var err error
	if gradeId != 0 {
		gradeSubjects, err := miscService.QueryGradeSubjects(gradeId)
		if err != nil {
			return 2, nil
		}

		for _, gradeSubject := range gradeSubjects {
			subject, err := models.ReadSubject(gradeSubject.SubjectId)
			if err != nil {
				continue
			}
			subjects = append(subjects, subject)
		}
	} else {
		subjects, err = miscService.QueryAllSubjects()
		if err != nil {
			return 2, nil
		}
	}
	return 0, subjects
}

func GetHelpItemList() (int64, []*models.HelpItem) {
	items, err := miscService.QueryAllHelpItems()
	if err != nil {
		return 2, nil
	}
	return 0, items
}

func GetAdvBanner(userId int64, version string) (int64, *models.AdvBanner, error) {
	user, err := models.ReadUser(userId)
	if err != nil {
		return 2, nil, errors.New("用户不存在")
	}
	advBanners, _ := miscService.QueryAllAdvBanners()
	for _, advBanner := range advBanners {
		if advBanner.AccessRight != 0 && advBanner.AccessRight != user.AccessRight {
			continue
		}
		if advBanner.Version == "all" || advBanner.Version == version {
			return 0, advBanner, nil
		} else if version < advBanner.Version[1:] {
			return 0, advBanner, nil
		}
	}
	return 2, nil, errors.New("未找到匹配的广告")
}

func VersionUpgrade(deviceType string, version int64) (int64, *UpgradeConfig, error) {
	config := UpgradeConfig{}
	upgradeInfo := settings.DeviceUpgradeInfo(deviceType)
	if upgradeInfo == nil {
		config.UpgradeType = "none"
		config.DownloadUrl = ""
		config.UpgradeContent = ""
	} else {
		if version <= upgradeInfo.ForceMinVersion {
			config.UpgradeType = "force"
			config.DownloadUrl = upgradeInfo.DownloadUrl
			config.UpgradeContent = upgradeInfo.UpgradeContent
		} else if version < upgradeInfo.MaxVersion && version > upgradeInfo.ForceMinVersion {
			config.UpgradeType = "common"
			config.DownloadUrl = upgradeInfo.DownloadUrl
			config.UpgradeContent = upgradeInfo.UpgradeContent
		} else {
			config.UpgradeType = "none"
			config.DownloadUrl = ""
			config.UpgradeContent = ""
		}
	}
	return 0, &config, nil
}

func GetQiniuDownloadUrl(mediaId string, width, height int64) (string, error) {
	domain := config.Env.Qiniu.Domain
	downloadUrl := domain + "/" + mediaId
	if width != 0 && height != 0 {
		downloadUrl += fmt.Sprintf("?imageView2/0/w/%d/h/%d", width, height)
	} else if width != 0 && height == 0 {
		downloadUrl += fmt.Sprintf("?imageView2/0/w/%d", width)
	} else if width == 0 && height != 0 {
		downloadUrl += fmt.Sprintf("?imageView2/0/h/%d", height)
	}
	return downloadUrl, nil
}

func GetQiniuUploadToken() (string, error) {
	bucket := config.Env.Qiniu.Bucket
	conf.ACCESS_KEY = config.Env.Qiniu.AccessKey
	conf.SECRET_KEY = config.Env.Qiniu.SecretKey
	c := kodo.New(0, nil)
	policy := &kodo.PutPolicy{
		Scope:   bucket,
		Expires: 3600,
	}
	token := c.MakeUptoken(policy)
	return token, nil
}
