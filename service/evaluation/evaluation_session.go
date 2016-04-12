// evaluation_session
package evaluation

import (
	"WolaiWebservice/config"
	"WolaiWebservice/models"

	"github.com/astaxie/beego/orm"
)

type courseSessionInfo struct {
	SessionId int64
	CourseId  int64
	ChapterId int64
	StudentId int64
}

func CheckSessionEvaluated(sessionId int64) bool {
	o := orm.NewOrm()
	exist := o.QueryTable(new(models.Evaluation).TableName()).Filter("session_id", sessionId).Exist()
	return exist
}

func GetLatestNotEvaluatedCourseSession(teacherId int64) (sessionId int64, courseId int64, chapterId int64, studentId int64, err error) {
	o := orm.NewOrm()
	qb, _ := orm.NewQueryBuilder(config.Env.Database.Type)
	qb.Select("orders.course_id,orders.chapter_id,sessions.id as session_id,orders.creator as student_id").
		From("orders").
		InnerJoin("sessions").On("orders.id = sessions.order_id").
		Where("sessions.tutor = ? and sessions.status = 'complete' and orders.chapter_id is not null and orders.chapter_id <> 0 and sessions.create_time > ?").
		OrderBy("sessions.create_time").Desc()
	sql := qb.String()
	var infos []courseSessionInfo
	_, err = o.Raw(sql, teacherId, "2016-04-07").QueryRows(&infos)
	if err == nil {
		for _, info := range infos {
			if !CheckSessionEvaluated(info.SessionId) {
				evaluationApply, _ := GetEvaluationApply(teacherId, info.ChapterId)
				if evaluationApply.Id != 0 && evaluationApply.Status != models.EVALUATION_APPLY_STATUS_IDLE {
					continue
				}
				sessionId = info.SessionId
				courseId = info.CourseId
				chapterId = info.ChapterId
				studentId = info.StudentId
				return
			}
		}
	}
	return
}
