// POICourseHandler
package handlers

import (
	"WolaiWebservice/models"
	"WolaiWebservice/utils"
	"time"

	"github.com/astaxie/beego/toolbox"
	seelog "github.com/cihub/seelog"
)

func POICourseExpiredHandler() {
	expiredTask := toolbox.NewTask("expiredTask", "0 1 0 * * *", func() error {
		processTime := time.Now().Format(utils.TIME_FORMAT)
		seelog.Debug("Process expired course:", processTime)
		userToCourses, _ := models.QueryExpiredCourses(processTime)
		for _, userToCourse := range userToCourses {
			updateInfo := map[string]interface{}{
				"Status": "expired",
			}
			models.UpdateUserCourseInfoById(userToCourse.Id, updateInfo)
		}
		return nil
	})
	toolbox.AddTask("expiredTask", expiredTask)
	toolbox.StartTask()
}
