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

type QaPkgShowInfo struct {
	*models.QaPkg
	Name    string `json:"showName"`
	Content string `json:"showContent"`
	Price   string `json:"showPrice"`
	Comment string `json:"showComment"`
}

type QaPkgModuleInfo struct {
	ModuleId      int64            `json:"moduleId"`
	ModuleName    string           `json:"moduleName"`
	ModuleComment string           `json:"moduleComment"`
	QaPkgs        []*QaPkgShowInfo `json:"qaPkgs"`
}

type MonthlyQaPkg struct {
	RecordId      int64  `json:"recordId"`
	QaPkgId       int64  `json:"qaPkgId"`
	Title         string `json:"title"`
	CurrentMonth  int64  `json:"currentMonth"`
	TotalMonth    int64  `json:"totalMonth"`
	LeftTime      int64  `json:"leftTime"`
	TotalTime     int64  `json:"totalTime"`
	ExpireComment string `json:"expireComment"`
}

type PermanentQaPkg struct {
	RecordId      int64  `json:"recordId"`
	QaPkgId       int64  `json:"qaPkgId"`
	Title         string `json:"title"`
	LeftTime      int64  `json:"leftTime"`
	TotalTime     int64  `json:"totalTime"`
	ExpireComment string `json:"expireComment"`
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
		return nil, errors.New("家教时间包资料异常")
	}
	for _, module := range modules {
		if module.Active != "Y" {
			continue
		}
		qaPkgs, err := qapkgService.GetModuleQaPkgs(module.Id)
		if err != nil {
			continue
		}

		var moduleInfo QaPkgModuleInfo
		moduleInfo.ModuleId = module.Id
		moduleInfo.ModuleName = module.Name
		moduleInfo.ModuleComment = module.Comment
		for _, qaPkg := range qaPkgs {
			showInfo := QaPkgShowInfo{}
			showInfo.QaPkg = qaPkg

			if qaPkg.Type == models.QA_PKG_TYPE_PERMANENT {
				showInfo.Name = fmt.Sprintf("%d分钟%s", qaPkg.TimeLength, qaPkg.Title)
				showInfo.Content = fmt.Sprintf("%d分钟", qaPkg.TimeLength)
				showInfo.Price = fmt.Sprintf("%.2f元（原价%.2f元）", float64(qaPkg.DiscountPrice)/100, float64(qaPkg.OriginalPrice)/100)
				showInfo.Comment = "购买该优惠包后可以任意使用在线家教功能"

			} else if qaPkg.Type == models.QA_PKG_TYPE_MONTHLY {
				showInfo.Name = fmt.Sprintf("%s-%d%s", module.Name, qaPkg.Month, "个月")
				showInfo.Content = fmt.Sprintf("%d分钟/月", qaPkg.TimeLength)
				showInfo.Price = fmt.Sprintf("%.2f元（原价%.2f元）", float64(qaPkg.DiscountPrice)/100, float64(qaPkg.OriginalPrice)/100)
				showInfo.Comment = "购买该优惠包后可以任意使用在线家教功能"
			}
			moduleInfo.QaPkgs = append(moduleInfo.QaPkgs, &showInfo)

		}
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
		qapkg := MonthlyQaPkg{}
		qapkg.TotalMonth = int64(len(monthlyQaPkgRecords))
		var index int64
		for _, record := range monthlyQaPkgRecords {
			index++
			if now.After(record.TimeFrom) && record.TimeTo.After(now) {
				qapkg.RecordId = record.Id
				qapkg.QaPkgId = record.QaPkgId
				qaPkg, _ := models.ReadQaPkg(record.QaPkgId)
				qaModule, _ := models.ReadQaPkgModule(qaPkg.ModuleId)
				qapkg.Title = qaModule.Name
				qapkg.TotalTime = qaPkg.TimeLength
				qapkg.LeftTime = record.LeftTime
				expireTime := record.TimeTo.Format("2006-01-02")
				qapkg.ExpireComment = fmt.Sprintf("＊%v%s", expireTime, "前有效")
				totalQaTimeLength = totalQaTimeLength + record.LeftTime
				qapkg.CurrentMonth = index
				userMonthlyQaPkgs = append(userMonthlyQaPkgs, &qapkg)
				break
			}
		}
	}
	detail.UserMonthlyQaPkgs = userMonthlyQaPkgs

	userPermanentQaPkgs := make([]*PermanentQaPkg, 0)
	permanentQaPkgRecords, err := qapkgService.GetPermanentQaPkgRecords(userId)
	if err != nil {
		return nil, err
	}
	for _, record := range permanentQaPkgRecords {
		qapkg := PermanentQaPkg{}
		qaPkg, _ := models.ReadQaPkg(record.QaPkgId)
		if record.Type == models.QA_PKG_TYPE_GIVEN {
			if !(now.After(record.TimeFrom) && record.TimeTo.After(now)) {
				continue
			}
			qapkg.Title = fmt.Sprintf("赠送家教时间%d分钟", qaPkg.TimeLength)
		} else {
			qapkg.Title = fmt.Sprintf("%d分钟%s", qaPkg.TimeLength, qaPkg.Title)
		}

		qapkg.RecordId = record.Id
		qapkg.QaPkgId = record.QaPkgId

		qapkg.TotalTime = qaPkg.TimeLength
		qapkg.LeftTime = record.LeftTime
		userPermanentQaPkgs = append(userPermanentQaPkgs, &qapkg)

		totalQaTimeLength = totalQaTimeLength + record.LeftTime
	}
	detail.UserPermanentQaPkgs = userPermanentQaPkgs
	detail.TotalQaTimeLength = totalQaTimeLength
	detail.HasDiscount = qapkgService.HasQaPkgDiscount()
	return &detail, nil
}
