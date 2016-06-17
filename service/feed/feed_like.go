// feed_like
package feed

import (
	"WolaiWebservice/config"
	"WolaiWebservice/models"

	"github.com/astaxie/beego/orm"
	"github.com/cihub/seelog"
)

func GetFeedLikeCount(feedId string) int64 {
	o := orm.NewOrm()
	count, err := o.QueryTable("feed_like").Filter("feed_id", feedId).Count()
	if err != nil {
		return 0
	}
	return count
}

func GetFeedFlowUserFeedLike(userId int64, start, pageCount int) models.Feeds {
	var feedIds []string
	o := orm.NewOrm()
	qb, _ := orm.NewQueryBuilder(config.Env.Database.Type)
	qb.Select("feed_id").From("feed_like").Where("user_id = ?").
		OrderBy("create_time").Asc().Limit(pageCount).Offset(start)
	sql := qb.String()
	_, err := o.Raw(sql, userId).QueryRows(&feedIds)
	if err != nil {
		seelog.Error("userId:", userId, " ", err.Error())
	}
	feeds := make(models.Feeds, len(feedIds))
	for i := range feedIds {
		feedId := feedIds[i]
		feed, err := GetFeed(feedId)
		if err == nil {
			feeds[i] = *(feed)
		}
	}
	return feeds
}

func GetFeedLikeList(feedId string) []models.User {
	o := orm.NewOrm()
	qb, _ := orm.NewQueryBuilder(config.Env.Database.Type)
	qb.Select("user_id").From("feed_like").Where("feed_id=?")
	sql := qb.String()
	var userIds []int64
	_, err := o.Raw(sql, feedId).QueryRows(&userIds)
	if err != nil {
		seelog.Error("feedId:", feedId, " ", err.Error())
	}
	users := make([]models.User, len(userIds))
	for i := range userIds {
		user, _ := models.ReadUser(userIds[i])
		if user == nil {
			continue
		}
		users[i] = *(user)
	}
	return users
}

func HasLikedFeed(feed *models.Feed, user *models.User) bool {
	o := orm.NewOrm()
	count, err := o.QueryTable("feed_like").
		Filter("feed_id", feed.Id).
		Filter("user_id", user.Id).Count()
	if err != nil {
		seelog.Error("feed:", feed, " user:", user, " ", err.Error())
		return false
	}
	if count > 0 {
		return true
	}
	return false
}
