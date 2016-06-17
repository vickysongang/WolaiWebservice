// feed_comment
package feed

import (
	"WolaiWebservice/config"
	"WolaiWebservice/models"
	"encoding/json"

	"github.com/astaxie/beego/orm"
	"github.com/cihub/seelog"
)

func GetFeedCommentCount(feedId string) int64 {
	o := orm.NewOrm()
	count, err := o.QueryTable("feed_comment").Filter("feed_id", feedId).Count()
	if err != nil {
		return 0
	}
	return count
}

func GetFeedComment(feedCommentId string) (*models.FeedComment, error) {
	feedComment := models.FeedComment{}
	o := orm.NewOrm()
	qb, _ := orm.NewQueryBuilder(config.Env.Database.Type)
	qb.Select("comment_id,feed_id,creator,create_time,text,image_info,reply_to").From("feed_comment").
		Where("comment_id = ?").OrderBy("create_time").Asc()
	sql := qb.String()
	err := o.Raw(sql, feedCommentId).QueryRow(&feedComment)
	if err != nil {
		seelog.Error("feedCommentId:", feedCommentId, " ", err.Error())
		return nil, err
	}
	timestampNano := feedComment.CreateTime.UnixNano()
	timestampMillis := timestampNano / 1000
	timestamp := float64(timestampMillis) / 1000000.0
	feedComment.CreateTimestamp = timestamp
	feedComment.Creator, _ = models.ReadUser(feedComment.Created)

	var imageList []string
	err = json.Unmarshal([]byte(feedComment.ImageInfo), &imageList)
	if err == nil {
		feedComment.ImageList = imageList
	} else {
		seelog.Error("unmarshal ImageInfo:", feedComment.ImageInfo, " ", err.Error())
		return nil, err
	}
	if feedComment.ReplyToId != 0 {
		feedComment.ReplyTo, _ = models.ReadUser(feedComment.ReplyToId)
	}
	return &feedComment, nil
}

func GetFeedComments(feedId string) models.FeedComments {
	o := orm.NewOrm()
	qb, _ := orm.NewQueryBuilder(config.Env.Database.Type)
	qb.Select("comment_id").From("feed_comment").Where("feed_id=?")
	sql := qb.String()
	var commentIds []string
	_, err := o.Raw(sql, feedId).QueryRows(&commentIds)
	if err != nil {
		seelog.Error("feedId:", feedId, " ", err.Error())
	}
	feedComments := make(models.FeedComments, len(commentIds))
	for i := range commentIds {
		feedComment, err := GetFeedComment(commentIds[i])
		if err == nil {
			feedComments[i] = *feedComment
		}
	}
	return feedComments
}
