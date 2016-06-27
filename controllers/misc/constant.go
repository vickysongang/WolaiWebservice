package misc

import (
	"WolaiWebservice/models"
	miscService "WolaiWebservice/service/misc"
	"errors"
)

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
