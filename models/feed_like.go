// feed_like
package models

import (
	"time"

	"github.com/astaxie/beego/orm"
	"github.com/cihub/seelog"
)

type FeedLike struct {
	Id         int64     `json:"-" orm:"pk"`
	FeedId     string    `json:"feedId"`
	UserId     int64     `json:"userId"`
	CreateTime time.Time `json:"-" orm:"auto_now_add;type(datetime)"`
}

func (fl *FeedLike) TableName() string {
	return "feed_like"
}
func init() {
	orm.RegisterModel(new(FeedLike))
}

func InsertFeedLike(feedLike *FeedLike) *FeedLike {
	o := orm.NewOrm()
	_, err := o.Insert(feedLike)
	if err != nil {
		seelog.Error("feedLike:", feedLike, " ", err.Error())
		return nil
	}
	return feedLike
}

func DeleteFeedLike(userId int64, feedId string) *FeedLike {
	feedLike := FeedLike{UserId: userId, FeedId: feedId}
	o := orm.NewOrm()
	_, err := o.QueryTable("feed_like").Filter("user_id", userId).Filter("feed_id", feedId).Delete()
	if err != nil {
		seelog.Error("feedLike:", feedLike, " ", err.Error())
		return nil
	}
	return &feedLike
}
