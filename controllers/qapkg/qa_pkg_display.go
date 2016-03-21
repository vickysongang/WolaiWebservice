// qa_pkg
package qapkg

import (
	"WolaiWebservice/models"
	"errors"
	"fmt"
	"time"

	qapkgService "WolaiWebservice/service/qapkg"

	"github.com/astaxie/beego/orm"
)

type QaPkgModuleInfo struct {
	ModuleId      int64           `json:"moduleId"`
	ModuleName    string          `json:"moduleName"`
	ModuleComment string          `json:"moduleComment"`
	QaPkgs        []*models.QaPkg `json:"qaPkgs"`
}

type MonthlyQaPkg struct {
	QaPkgId       int64  `json:"qaPkgId"`
	Title         string `json:"title"`
	CurrentMonth  int64  `json:"currentMonth"`
	TotalMonth    int64  `json:"totalMonth"`
	LeftTime      int64  `json:"leftTime"`
	TotalTime     int64  `json:"totalTime"`
	ExpireComment string `json:"expireComment"`
}

type PermanentQaPkg struct {
	QaPkgId   int64  `json:"qaPkgId"`
	Title     string `json:"title"`
	LeftTime  int64  `json:"leftTime"`
	TotalTime int64  `json:"totalTime"`
}

type UserQaPkgDetail struct {
	TotalQaTimeLength   int64             `json:"totalQaTimeLength"`
	UserMonthlyQaPkgs   []*MonthlyQaPkg   `json:"userMonthlyQaPkgs"`
	UserPermanentQaPkgs []*PermanentQaPkg `json:"userPermanentQaPkgs"`
	HasDiscount         bool              `json:"hasDiscount"`
}

func GetQaPkgList() ([]QaPkgModuleInfo, error) {
	pkgModuleInfos := make([]QaPkgModuleInfo, 0)
	o := orm.NewOrm()
	var modules []models.QaPkgModule
	_, err := o.QueryTable(new(models.QaPkgModule).TableName()).OrderBy("rank").All(&modules)
	if err != nil {
		return nil, errors.New("答疑包资料异常")
	}
	for _, module := range modules {
		qaPkgs, err := qapkgService.GetModuleQaPkgs(module.Id)
		if err != nil {
			continue
		}
		var moduleInfo QaPkgModuleInfo
		moduleInfo.ModuleId = module.Id
		moduleInfo.ModuleName = module.Name
		moduleInfo.ModuleComment = module.Comment
		moduleInfo.QaPkgs = qaPkgs
		pkgModuleInfos = append(pkgModuleInfos, moduleInfo)
	}
	return pkgModuleInfos, nil
}

func GetQaPkgDetail(userId int64) (*UserQaPkgDetail, error) {
	detail := UserQaPkgDetail{}
	now := time.Now()

	userMonthlyQaPkgs := make([]*MonthlyQaPkg, 0)
	monthlyQaPkgRecords, err := qapkgService.GetMonthlyQaPkgRecords(userId)
	if err != nil {
		return nil, err
	}
	var totalQaTimeLength int64
	if len(monthlyQaPkgRecords) > 0 {
		userMonthlyQaPkg := MonthlyQaPkg{}
		userMonthlyQaPkg.TotalMonth = int64(len(monthlyQaPkgRecords))
		var index int64
		for _, monthlyQaPkgRecord := range monthlyQaPkgRecords {
			index++
			if now.After(monthlyQaPkgRecord.TimeFrom) && monthlyQaPkgRecord.TimeTo.After(now) {
				userMonthlyQaPkg.QaPkgId = monthlyQaPkgRecord.QaPkgId
				qaPkg, _ := models.ReadQaPkg(monthlyQaPkgRecord.QaPkgId)
				qaModule, _ := models.ReadQaPkgModule(qaPkg.ModuleId)
				userMonthlyQaPkg.Title = qaModule.Name
				userMonthlyQaPkg.TotalTime = qaPkg.TimeLength
				userMonthlyQaPkg.LeftTime = monthlyQaPkgRecord.LeftTime
				expireTime := monthlyQaPkgRecord.TimeTo.Format("2006-01-02")
				userMonthlyQaPkg.ExpireComment = fmt.Sprintf("＊%v%s", expireTime, "前有效")
				totalQaTimeLength = totalQaTimeLength + monthlyQaPkgRecord.LeftTime
				break
			}
		}
		userMonthlyQaPkg.CurrentMonth = index
		userMonthlyQaPkgs = append(userMonthlyQaPkgs, &userMonthlyQaPkg)
	}
	detail.UserMonthlyQaPkgs = userMonthlyQaPkgs

	userPermanentQaPkgs := make([]*PermanentQaPkg, 0)
	permanentQaPkgRecords, err := qapkgService.GetPermanentQaPkgRecords(userId)
	if err != nil {
		return nil, err
	}
	for _, permanentQaPkgRecord := range permanentQaPkgRecords {
		userPermanentQaPkg := PermanentQaPkg{}
		userPermanentQaPkg.QaPkgId = permanentQaPkgRecord.QaPkgId
		qaPkg, _ := models.ReadQaPkg(permanentQaPkgRecord.QaPkgId)
		qaModule, _ := models.ReadQaPkgModule(qaPkg.ModuleId)
		userPermanentQaPkg.Title = qaModule.Name
		userPermanentQaPkg.TotalTime = qaPkg.TimeLength
		userPermanentQaPkg.LeftTime = permanentQaPkgRecord.LeftTime
		userPermanentQaPkgs = append(userPermanentQaPkgs, &userPermanentQaPkg)

		totalQaTimeLength = totalQaTimeLength + permanentQaPkgRecord.LeftTime
	}
	detail.UserPermanentQaPkgs = userPermanentQaPkgs
	detail.TotalQaTimeLength = totalQaTimeLength
	detail.HasDiscount = qapkgService.HasQaPkgDiscount()
	return &detail, nil
}
