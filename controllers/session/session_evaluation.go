package session

import (
	"math/rand"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"

	"WolaiWebservice/models"
	evaluationService "WolaiWebservice/service/evaluation"
)

type evaluationInfo struct {
	Type      string             `json:"type"`
	Evalution *models.Evaluation `json:"evaluationInfo"`
}

func CheckRandNumInSlice(slice []int64, randNum int64) bool {
	for _, v := range slice {
		if v == randNum {
			return true
		}
	}
	return false
}

func GetRandNumSlice(sliceSize int64, length int64) []int64 {
	var result []int64
	if sliceSize <= 0 || length <= 0 {
		return result
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for {
		randNum := r.Int63n(length)
		if !CheckRandNumInSlice(result, randNum) {
			result = append(result, randNum)
		}
		if sliceSize > length {
			if int64(len(result)) == length {
				break
			}
		} else {
			if int64(len(result)) == sliceSize {
				break
			}
		}
	}
	return result
}

func QuerySystemEvaluationLabels(userId, sessionId, count int64) ([]*models.EvaluationLabel, error) {
	labels := []*models.EvaluationLabel{}
	session, _ := models.ReadSession(sessionId)

	//如果当前用户是学生，则要返回老师的标签信息，如果当前用户是老师，则要返回学生的标签信息
	//学生
	if userId == session.Creator {
		teacher, _ := models.ReadUser(session.Tutor)
		//个人标签
		teacherPersonalLabels, err := evaluationService.QueryEvaluationLabels(teacher.Gender, models.PERSONAL_EVALUATION_LABEL, models.TEACHER_EVALUATION_LABEL)
		if err != nil {
			return nil, err
		}
		personalCount := (count - 2) / 2
		for _, v := range GetRandNumSlice(personalCount, int64(len(teacherPersonalLabels))) {
			labels = append(labels, teacherPersonalLabels[v])
		}
		//讲课风格
		teacherStyleLabels, err := evaluationService.QueryEvaluationLabels(teacher.Gender, models.STYLE_EVALUATION_LABEL, models.TEACHER_EVALUATION_LABEL)
		if err != nil {
			return nil, err
		}
		for _, v := range GetRandNumSlice(count-personalCount-2, int64(len(teacherStyleLabels))) {
			labels = append(labels, teacherStyleLabels[v])
		}
		//科目标签
		order, _ := models.ReadOrder(session.OrderId)
		teacherSubjectLabels, err := evaluationService.QueryEvaluationLabelsBySubject(order.SubjectId)
		if err != nil {
			return nil, err
		}
		for _, v := range GetRandNumSlice(2, int64(len(teacherSubjectLabels))) {
			labels = append(labels, teacherSubjectLabels[v])
		}
	} else if userId == session.Tutor { //老师
		student, _ := models.ReadUser(session.Creator)
		//个人标签
		studentPersonalLabels, err := evaluationService.QueryEvaluationLabels(student.Gender, models.PERSONAL_EVALUATION_LABEL, models.STUDENT_EVALUATION_LABEL)
		if err != nil {
			return nil, err
		}
		personalCount := count / 2
		for _, v := range GetRandNumSlice(personalCount, int64(len(studentPersonalLabels))) {
			labels = append(labels, studentPersonalLabels[v])
		}
		//能力程度
		studentAbilityLabels, err := evaluationService.QueryEvaluationLabels(student.Gender, models.ABILITY_EVALUATION_LABEL, models.STUDENT_EVALUATION_LABEL)
		if err != nil {
			return nil, err
		}
		for _, v := range GetRandNumSlice(count-personalCount, int64(len(studentAbilityLabels))) {
			labels = append(labels, studentAbilityLabels[v])
		}
	}
	return labels, nil
}

func CreateEvaluation(userId, targetId, sessionId, chapterId int64, evaluationContent string) (*models.Evaluation, error) {
	user, _ := models.ReadUser(userId)
	if user.AccessRight == models.USER_ACCESSRIGHT_TEACHER { //导师插入评价

		if chapterId == 0 { //兼容旧版，旧版未传该字段
			session, err := models.ReadSession(sessionId)
			if err == nil {
				order, err := models.ReadOrder(session.OrderId)
				if err == nil {
					chapterId = order.ChapterId
				}
			}
		}

		if chapterId != 0 { // 课程插入评价申请
			apply, _ := evaluationService.GetEvaluationApply(userId, chapterId, 0)
			chapter, _ := models.ReadCourseCustomChapter(chapterId)
			if apply.Id == 0 {
				evaluationApply := models.EvaluationApply{
					UserId:    userId,
					SessionId: sessionId,
					CourseId:  chapter.CourseId,
					ChapterId: chapterId,
					Status:    models.EVALUATION_APPLY_STATUS_CREATED,
					Content:   evaluationContent,
				}
				_, err := models.InsertEvaluationApply(&evaluationApply)
				if err != nil {
					return nil, err
				}
			} else {
				if apply.Status == models.EVALUATION_APPLY_STATUS_IDLE {
					eveluationInfo := map[string]interface{}{
						"SessionId": sessionId,
						"Status":    models.EVALUATION_APPLY_STATUS_CREATED,
						"Content":   evaluationContent,
					}
					models.UpdateEvaluationApply(apply.Id, eveluationInfo)
				}
			}
			return nil, nil
		} else { //答疑插入评价
			content, err := evaluateSession(sessionId, userId, targetId, chapterId, evaluationContent)
			return content, err
		}
	} else { //学生插入评价
		content, err := evaluateSession(sessionId, userId, targetId, chapterId, evaluationContent)
		return content, err
	}
	return nil, nil
}

func evaluateSession(sessionId, userId, targetId, chapterId int64, evaluationContent string) (*models.Evaluation, error) {
	session, _ := models.ReadSession(sessionId)
	if targetId == 1 {
		if userId == session.Creator {
			targetId = session.Tutor
		} else {
			targetId = session.Creator
		}
	}
	oldEvaluation, _ := evaluationService.QueryEvaluation(userId, sessionId)
	if oldEvaluation.Id == 0 {
		evaluation := models.Evaluation{
			UserId:    userId,
			TargetId:  targetId,
			SessionId: sessionId,
			ChapterId: chapterId,
			Content:   evaluationContent}
		content, err := models.InsertEvaluation(&evaluation)
		return content, err
	} else {
		evaluationInfo := map[string]interface{}{
			"UserId":    userId,
			"TargetId":  targetId,
			"SessionId": sessionId,
			"ChapterId": chapterId,
			"Content":   evaluationContent,
		}
		err := models.UpdateEvaluation(oldEvaluation.Id, evaluationInfo)
		if err != nil {
			return nil, err
		} else {
			content, err := models.ReadEvaluation(oldEvaluation.Id)
			return content, err
		}
	}
	return nil, nil
}

//这是版本兼容导致的一坨屎，请绕行
func QueryEvaluationInfo(userId, sessionId, targetId, chapterId int64) ([]*evaluationInfo, error) {
	evalutionInfos := make([]*evaluationInfo, 0)
	selfEvaluation := evaluationInfo{}
	otherEvaluation := evaluationInfo{}
	var isStudent bool
	var studentEvaluation, teacherEvaluation *models.Evaluation
	if sessionId != 0 {
		session, _ := models.ReadSession(sessionId)
		studentEvaluation, _ = evaluationService.QueryEvaluation(session.Creator, sessionId)
		teacherEvaluation, _ = evaluationService.QueryEvaluation(session.Tutor, sessionId)
		if userId == session.Creator {
			isStudent = true
		} else if userId == session.Tutor {
			isStudent = false
		}
	} else {
		chapter, _ := models.ReadCourseCustomChapter(chapterId)
		studentEvaluation, _ = evaluationService.QueryEvaluationByChapter(chapter.UserId, chapterId, 0)
		teacherEvaluation, _ = evaluationService.QueryEvaluationByChapter(chapter.TeacherId, chapterId, 0)
		if userId == chapter.UserId {
			isStudent = true
		} else if userId == chapter.TeacherId {
			isStudent = false
		}
	}

	//旧版评价表里targetId为0，新版不为0，故根据该字段来判断获取的是旧版评论还是新版评论
	if isStudent {
		//导师和学生都有评论
		if teacherEvaluation.Id != 0 && studentEvaluation.Id != 0 {
			//导师是新版学生是旧版或导师是旧版学生是新版，考虑到兼容性的问题，学生看到的评论内容忽略掉导师的评论
			if (teacherEvaluation.TargetId != 0 && studentEvaluation.TargetId == 0) || (teacherEvaluation.TargetId == 0 && studentEvaluation.TargetId != 0) {
				if targetId != 0 { //如果学生版本是新版，查看旧版内容考虑兼容性忽略旧版评论
					if !(strings.HasPrefix(studentEvaluation.Content, "[") && strings.HasSuffix(studentEvaluation.Content, "]")) {
						selfEvaluation.Type = "student"
						selfEvaluation.Evalution = studentEvaluation
						evalutionInfos = append(evalutionInfos, &selfEvaluation)
					}
				} else {
					if !(strings.HasPrefix(studentEvaluation.Content, "{") && strings.HasSuffix(studentEvaluation.Content, "}")) {
						selfEvaluation.Type = "student"
						selfEvaluation.Evalution = studentEvaluation
						evalutionInfos = append(evalutionInfos, &selfEvaluation)
					}
				}
			} else {
				if targetId == 0 {
					if !(strings.HasPrefix(studentEvaluation.Content, "{") && strings.HasSuffix(studentEvaluation.Content, "}")) {
						selfEvaluation.Type = "student"
						selfEvaluation.Evalution = studentEvaluation
						evalutionInfos = append(evalutionInfos, &selfEvaluation)
					}
					if !(strings.HasPrefix(teacherEvaluation.Content, "{") && strings.HasSuffix(teacherEvaluation.Content, "}")) {
						otherEvaluation.Type = "teacher"
						otherEvaluation.Evalution = teacherEvaluation
						evalutionInfos = append(evalutionInfos, &otherEvaluation)
					}
				} else {
					if !(strings.HasPrefix(studentEvaluation.Content, "[") && strings.HasSuffix(studentEvaluation.Content, "]")) {
						selfEvaluation.Type = "student"
						selfEvaluation.Evalution = studentEvaluation
						evalutionInfos = append(evalutionInfos, &selfEvaluation)
					}
					if !(strings.HasPrefix(teacherEvaluation.Content, "[") && strings.HasSuffix(teacherEvaluation.Content, "]")) {
						otherEvaluation.Type = "teacher"
						otherEvaluation.Evalution = teacherEvaluation
						evalutionInfos = append(evalutionInfos, &otherEvaluation)
					}
				}
			}
		} else if teacherEvaluation.Id == 0 && studentEvaluation.Id != 0 { //导师未评论，学生有评论
			if targetId != 0 { //如果学生版本是新版，查看旧版内容考虑兼容性忽略旧版评论
				if !(strings.HasPrefix(studentEvaluation.Content, "[") && strings.HasSuffix(studentEvaluation.Content, "]")) {
					selfEvaluation.Type = "student"
					selfEvaluation.Evalution = studentEvaluation
					evalutionInfos = append(evalutionInfos, &selfEvaluation)
				}
			} else {
				if !(strings.HasPrefix(studentEvaluation.Content, "{") && strings.HasSuffix(studentEvaluation.Content, "}")) {
					selfEvaluation.Type = "student"
					selfEvaluation.Evalution = studentEvaluation
					evalutionInfos = append(evalutionInfos, &selfEvaluation)
				}
			}
		} else if teacherEvaluation.Id != 0 && studentEvaluation.Id == 0 { //导师有评论，学生未评论
			if (teacherEvaluation.TargetId != 0 && targetId != 0) || (teacherEvaluation.TargetId == 0 && targetId == 0) {
				otherEvaluation.Type = "teacher"
				otherEvaluation.Evalution = teacherEvaluation
				evalutionInfos = append(evalutionInfos, &otherEvaluation)
			}
		}
	} else {
		if teacherEvaluation.Id != 0 { //导师只能看到自己的评论，不能看到学生的评论
			if targetId != 0 { //如果导师当前版本是新版，查看旧版评价内容时由于兼容性的问题忽略旧评论
				if !(strings.HasPrefix(teacherEvaluation.Content, "[") && strings.HasSuffix(teacherEvaluation.Content, "]")) {
					selfEvaluation.Type = "teacher"
					selfEvaluation.Evalution = teacherEvaluation
					evalutionInfos = append(evalutionInfos, &selfEvaluation)
				}
			} else {
				if !(strings.HasPrefix(teacherEvaluation.Content, "{") && strings.HasSuffix(teacherEvaluation.Content, "}")) {
					selfEvaluation.Type = "teacher"
					selfEvaluation.Evalution = teacherEvaluation
					evalutionInfos = append(evalutionInfos, &selfEvaluation)
				}
			}
		}
	}
	return evalutionInfos, nil
}

func HasStudentSessionRecordEvaluated(sessionId int64, studentId int64) bool {
	o := orm.NewOrm()
	count, err := o.QueryTable("evaluation").Filter("session_id", sessionId).Filter("user_id", studentId).Count()
	if err != nil {
		return false
	}
	if count > 0 {
		return true
	}
	return false
}

func HasTeacherSessionRecordEvaluated(sessionId int64, teacherId int64) bool {
	o := orm.NewOrm()
	session, err := models.ReadSession(sessionId)
	if err != nil {
		return false
	}
	order, err := models.ReadOrder(session.OrderId)
	if err != nil {
		return false
	}
	if order.ChapterId != 0 {
		count, err := o.QueryTable("evaluation").Filter("chapter_id", order.ChapterId).Filter("user_id", teacherId).Count()
		if err != nil {
			return false
		}
		if count > 0 {
			return true
		}
	} else {
		count, err := o.QueryTable("evaluation").Filter("session_id", sessionId).Filter("user_id", teacherId).Count()
		if err != nil {
			return false
		}
		if count > 0 {
			return true
		}
	}
	return false
}
