// feed
package feed

import (
	"WolaiWebservice/config"
	"WolaiWebservice/models"
	"encoding/json"

	"github.com/astaxie/beego/orm"
	"github.com/cihub/seelog"
)

func GetFeed(feedId string) (*models.Feed, error) {
	feed := models.Feed{}
	qb, _ := orm.NewQueryBuilder(config.Env.Database.Type)
	qb.Select("feed_id,creator,create_time,feed_type,text,image_info,attribute_info,origin_feed_id").
		From("feed").Where("feed_id = ?")
	sql := qb.String()
	o := orm.NewOrm()
	err := o.Raw(sql, feedId).QueryRow(&feed)
	if err != nil {
		return nil, err
	}
	timestampNano := feed.CreateTime.UnixNano()
	timestampMillis := timestampNano / 1000
	timestamp := float64(timestampMillis) / 1000000.0
	feed.CreateTimestamp = timestamp
	creator, _ := models.ReadUser(feed.Created)
	feed.Creator = creator
	var imageList []string
	err = json.Unmarshal([]byte(feed.ImageInfo), &imageList)
	if err == nil {
		feed.ImageList = imageList
	} else {
		seelog.Error("unmarshal ImageInfo:", feed.ImageInfo, " ", err.Error())
		return nil, err
	}
	var attributeMap map[string]string
	err = json.Unmarshal([]byte(feed.AttributeInfo), &attributeMap)
	if err == nil {
		feed.Attribute = attributeMap
	} else {
		seelog.Error("unmarshal AttributeInfo:", feed.AttributeInfo, " ", err.Error())
		return nil, err
	}
	if feed.OriginFeedId != "" && len(feed.OriginFeedId) > 0 {
		originFeed, err := GetFeed(feed.OriginFeedId)
		if err != nil {
			return nil, err
		}
		feed.OriginFeed = originFeed
	}
	feed.LikeCount = GetFeedLikeCount(feed.Id)
	feed.CommentCount = GetFeedCommentCount(feed.Id)
	return &feed, nil
}

func GetFeedFlowAtrium(start, pageCount int) (models.Feeds, error) {
	o := orm.NewOrm()
	var feedIds []string
	qb, _ := orm.NewQueryBuilder(config.Env.Database.Type)
	qb.Select("feed.feed_id").From("feed").
		InnerJoin("users").On("feed.creator = users.id").
		OrderBy("feed.create_time").Desc().Limit(pageCount).Offset(start)
	sql := qb.String()
	_, err := o.Raw(sql).QueryRows(&feedIds)
	feeds := make(models.Feeds, len(feedIds))
	if err != nil {
		seelog.Error(err.Error())
		return feeds, err
	}
	for i := range feedIds {
		feed, err := GetFeed(feedIds[i])
		if err == nil {
			feeds[i] = *feed
		}
	}
	return feeds, nil
}

//获取置顶动态
func GetTopFeedFlowAtrium(plateType string) (models.Feeds, error) {
	o := orm.NewOrm()
	var feedIds []string
	qb, _ := orm.NewQueryBuilder(config.Env.Database.Type)
	qb.Select("feed.feed_id").From("feed").
		InnerJoin("users").On("feed.creator = users.id").
		Where("feed.plate_type like ? and feed.top_seq is not null and feed.top_seq <> ''").
		OrderBy("feed.top_seq").Desc()
	sql := qb.String()
	_, err := o.Raw(sql, "%"+plateType+"%").QueryRows(&feedIds)
	feeds := make(models.Feeds, len(feedIds))
	if err != nil {
		seelog.Error(err.Error())
		return feeds, err
	}
	for i := range feedIds {
		feed, err := GetFeed(feedIds[i])
		feed.TopFlag = true
		if err == nil {
			feeds[i] = *feed
		}
	}
	return feeds, nil
}

//获取板块动态
func GetFeedFlowAtriumByPlateType(start, pageCount int, plateType string) (models.Feeds, error) {
	o := orm.NewOrm()
	var feedIds []string
	qb, _ := orm.NewQueryBuilder(config.Env.Database.Type)
	qb.Select("feed.feed_id").From("feed").
		InnerJoin("users").On("feed.creator = users.id").
		Where("feed.plate_type like ? and (feed.top_seq is null or feed.top_seq='')").
		OrderBy("feed.create_time").Desc().Limit(pageCount).Offset(start)
	sql := qb.String()
	_, err := o.Raw(sql, "%"+plateType+"%").QueryRows(&feedIds)
	feeds := make(models.Feeds, len(feedIds))
	if err != nil {
		seelog.Error(err.Error())
		return feeds, err
	}
	for i := range feedIds {
		feed, err := GetFeed(feedIds[i])
		feed.TopFlag = false
		if err == nil {
			feeds[i] = *feed
		}
	}
	return feeds, nil
}

func GetFeedFlowUserFeed(userId int64, start, pageCount int) models.Feeds {
	o := orm.NewOrm()
	var feedIds []string
	qb, _ := orm.NewQueryBuilder(config.Env.Database.Type)
	qb.Select("feed_id").From("feed").
		Where("creator = ?").
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
