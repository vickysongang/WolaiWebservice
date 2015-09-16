// POIEvaluation.go
package models

import (
	"time"

	"github.com/astaxie/beego/orm"
	"github.com/tmhenry/POIWolaiWebService/utils"
)

type POIEvaluation struct {
	Id         int64     `json:"-" orm:"pk"`
	UserId     int64     `json:"userId"`
	SessionId  int64     `json:"sessionId"`
	Content    string    `json:"content"`
	CreateTime time.Time `orm:"auto_now_add;type(datetime)"`
}

type POIEvaluationLabel struct {
	Id            int64  `json:"-" orm:"pk"`
	Name          string `json:"name"`
	GenderType    int64  `json:"-"`
	AttributeType string `json:"-"`
	ObjectType    string `json:"-"`
}

type POIEvaluationLabels []*POIEvaluationLabel

type POIEvaluationInfo struct {
	Type      string         `json:"type"`
	Evalution *POIEvaluation `json:"evaluationInfo"`
}

type POIEvaluationInfos []POIEvaluationInfo

func (e *POIEvaluation) TableName() string {
	return "evaluation"
}

func (el *POIEvaluationLabel) TableName() string {
	return "evaluation_label"
}

func init() {
	orm.RegisterModel(new(POIEvaluation), new(POIEvaluationLabel))
}

func InsertEvaluation(evalution *POIEvaluation) (*POIEvaluation, error) {
	o := orm.NewOrm()
	id, err := o.Insert(evalution)
	if err != nil {
		return nil, err
	}
	evalution.Id = id
	return evalution, nil
}

func InsertEvaluationLabel(evalutionLabel *POIEvaluationLabel) (int64, error) {
	o := orm.NewOrm()
	id, err := o.Insert(evalutionLabel)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func QueryEvaluation4Self(userId, sessionId int64) (*POIEvaluation, error) {
	o := orm.NewOrm()
	qb, _ := orm.NewQueryBuilder(utils.DB_TYPE)
	qb.Select("id,user_id,session_id,content,create_time").From("evaluation").
		Where("user_id = ? and session_id = ?")
	sql := qb.String()
	evalution := POIEvaluation{}
	err := o.Raw(sql, userId, sessionId).QueryRow(&evalution)
	if err != nil {
		return nil, err
	}
	return &evalution, nil
}

func QueryEvaluation4Other(userId, sessionId int64) (*POIEvaluation, error) {
	o := orm.NewOrm()
	qb, _ := orm.NewQueryBuilder(utils.DB_TYPE)
	qb.Select("id,user_id,session_id,content,create_time").From("evaluation").
		Where("user_id <> ? and session_id = ?")
	sql := qb.String()
	evalution := POIEvaluation{}
	err := o.Raw(sql, userId, sessionId).QueryRow(&evalution)
	if err != nil {
		return nil, err
	}
	return &evalution, nil
}

/*
 * 根据条件查询系统推荐标签
 * genderType 性别类型，0代表女性，1代表男性，2代表中性
 * attributeType 属性分类，personal代表个人标签，style代表讲课风格，subject代表科目标签，ability代表能力程度
 * objectType 对象分类，student代表学生，teacher代表老师,both代表两者均可以
 */
func QueryEvaluationLabels(genderType int64, attributeType, objectType string) (POIEvaluationLabels, error) {
	labels := make(POIEvaluationLabels, 0)
	o := orm.NewOrm()
	qb, _ := orm.NewQueryBuilder(utils.DB_TYPE)
	qb.Select("name").From("evaluation_label").Where("gender_type in (?,2) and attribute_type = ? and object_type in (?,'both')")
	sql := qb.String()
	_, err := o.Raw(sql, genderType, attributeType, objectType).QueryRows(&labels)
	if err != nil {
		return nil, err
	}
	return labels, nil
}

func QueryEvaluationLabelsBySubject(subjectId int64) (POIEvaluationLabels, error) {
	labels := POIEvaluationLabels{}
	o := orm.NewOrm()
	qb, _ := orm.NewQueryBuilder(utils.DB_TYPE)
	qb.Select("evaluation_label.name").From("evaluation_label").
		InnerJoin("evaluation_to_subject").On("evaluation_label.id = evaluation_to_subject.label_id").
		Where("evaluation_label.attribute_type = 'subject' and evaluation_to_subject.subject_id = ?")
	sql := qb.String()
	_, err := o.Raw(sql, subjectId).QueryRows(&labels)
	if err != nil {
		return nil, err
	}
	return labels, nil
}

func QueryEvaluationInfo(userId, sessionId int64) (*POIEvaluationInfos, error) {
	session := QuerySessionById(sessionId)
	self, err1 := QueryEvaluation4Self(userId, sessionId)
	other, err2 := QueryEvaluation4Other(userId, sessionId)

	selfEvaluation := POIEvaluationInfo{}
	otherEvaluation := POIEvaluationInfo{}

	evalutionInfos := make(POIEvaluationInfos, 0)
	if userId == session.Tutor {
		if err1 == nil {
			selfEvaluation.Type = "teacher"
			selfEvaluation.Evalution = self

			evalutionInfos = append(evalutionInfos, selfEvaluation)
		}
		if err2 == nil {
			otherEvaluation.Type = "student"
			otherEvaluation.Evalution = other

			evalutionInfos = append(evalutionInfos, otherEvaluation)
		}
	} else if userId == session.Created {
		if err1 == nil {
			selfEvaluation.Type = "student"
			selfEvaluation.Evalution = self

			evalutionInfos = append(evalutionInfos, selfEvaluation)
		}
		if err2 == nil {
			otherEvaluation.Type = "teacher"
			otherEvaluation.Evalution = other

			evalutionInfos = append(evalutionInfos, otherEvaluation)
		}
	}
	return &evalutionInfos, nil
}

func HasOrderInSessionEvaluated(sessionId int64) bool {
	o := orm.NewOrm()
	count, err := o.QueryTable("evaluation").Filter("session_id", sessionId).Count()
	if err != nil {
		return false
	}
	if count > 1 {
		return true
	}
	return false
}
