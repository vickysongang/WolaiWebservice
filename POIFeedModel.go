// POIFeedModel.go
package main

import (
	"encoding/json"

	"github.com/astaxie/beego/orm"
	seelog "github.com/cihub/seelog"
)

func InsertPOIFeed(userId int64, feedId string, feedType int64, text string, imageStr string,
	originFeedId string, attributeStr string) (*POIFeed, error) {
	feed := POIFeed{Created: userId, FeedType: feedType, Text: text, ImageInfo: imageStr,
		AttributeInfo: attributeStr, Id: feedId, OriginFeedId: originFeedId}
	o := orm.NewOrm()
	_, err := o.Insert(&feed)
	if err != nil {
		seelog.Error("feed:", feed, " ", err.Error())
		return nil, err
	}
	return &feed, nil
}

func GetFeed(feedId string) (*POIFeed, error) {
	feed := POIFeed{}
	qb, _ := orm.NewQueryBuilder(DB_TYPE)
	qb.Select("feed_id,creator,create_time,feed_type,text,image_info,attribute_info,origin_feed_id").
		From("feed").Where("feed_id = ?")
	sql := qb.String()
	o := orm.NewOrm()
	err := o.Raw(sql, feedId).QueryRow(&feed)
	if err != nil {
		seelog.Error("feedId:", feedId, " ", err.Error())
		return nil, err
	}
	timestampNano := feed.CreateTime.UnixNano()
	timestampMillis := timestampNano / 1000
	timestamp := float64(timestampMillis) / 1000000.0
	feed.CreateTimestamp = timestamp
	creator := QueryUserById(feed.Created)
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
	feed.LikeCount = GetPOIFeedLikeCount(feed.Id)
	feed.CommentCount = GETPOIFeedCommentCount(feed.Id)
	return &feed, nil
}

func InsertPOIFeedLike(userId int64, feedId string) *POIFeedLike {
	o := orm.NewOrm()
	feedLike := POIFeedLike{UserId: userId, FeedId: feedId}
	_, err := o.Insert(&feedLike)
	if err != nil {
		seelog.Error("feedLike:", feedLike, " ", err.Error())
		return nil
	}
	return &feedLike
}

func DeletePOIFeedLike(userId int64, feedId string) *POIFeedLike {
	feedLike := POIFeedLike{UserId: userId, FeedId: feedId}
	o := orm.NewOrm()
	_, err := o.QueryTable("feed_like").Filter("user_id", userId).Filter("feed_id", feedId).Delete()
	if err != nil {
		seelog.Error("feedLike:", feedLike, " ", err.Error())
		return nil
	}
	return &feedLike
}

func GetPOIFeedLikeCount(feedId string) int64 {
	o := orm.NewOrm()
	count, err := o.QueryTable("feed_like").Filter("feed_id", feedId).Count()
	if err != nil {
		seelog.Error("feedId:", feedId, " ", err.Error())
		return 0
	}
	return count
}

func GETPOIFeedCommentCount(feedId string) int64 {
	o := orm.NewOrm()
	count, err := o.QueryTable("feed_comment").Filter("feed_id", feedId).Count()
	if err != nil {
		seelog.Error("feedId:", feedId, " ", err.Error())
		return 0
	}
	return count
}

func GetFeedLikeList(feedId string) POIUsers {
	o := orm.NewOrm()
	qb, _ := orm.NewQueryBuilder(DB_TYPE)
	qb.Select("user_id").From("feed_like").Where("feed_id=?")
	sql := qb.String()
	var userIds []int64
	_, err := o.Raw(sql, feedId).QueryRows(&userIds)
	if err != nil {
		seelog.Error("feedId:", feedId, " ", err.Error())
	}
	users := make(POIUsers, len(userIds))
	for i := range userIds {
		users[i] = *(QueryUserById(userIds[i]))
	}
	return users
}

func HasLikedFeed(feed *POIFeed, user *POIUser) bool {
	o := orm.NewOrm()
	count, err := o.QueryTable("feed_like").Filter("feed_id", feed.Id).Filter("user_id", user.UserId).Count()
	if err != nil {
		seelog.Error("feed:", feed, " user:", user, " ", err.Error())
		return false
	}
	if count > 0 {
		return true
	}
	return false
}

func GetFeedComment(feedCommentId string) (*POIFeedComment, error) {
	feedComment := POIFeedComment{}
	o := orm.NewOrm()
	qb, _ := orm.NewQueryBuilder(DB_TYPE)
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
	feedComment.Creator = QueryUserById(feedComment.Created)

	var imageList []string
	err = json.Unmarshal([]byte(feedComment.ImageInfo), &imageList)
	if err == nil {
		feedComment.ImageList = imageList
	} else {
		seelog.Error("unmarshal ImageInfo:", feedComment.ImageInfo, " ", err.Error())
		return nil, err
	}
	if feedComment.ReplyToId != 0 {
		feedComment.ReplyTo = QueryUserById(feedComment.ReplyToId)
	}
	return &feedComment, nil
}

func GetFeedComments(feedId string) POIFeedComments {
	o := orm.NewOrm()
	qb, _ := orm.NewQueryBuilder(DB_TYPE)
	qb.Select("comment_id").From("feed_comment").Where("feed_id=?")
	sql := qb.String()
	var commentIds []string
	_, err := o.Raw(sql, feedId).QueryRows(&commentIds)
	if err != nil {
		seelog.Error("feedId:", feedId, " ", err.Error())
	}
	feedComments := make(POIFeedComments, len(commentIds))
	for i := range commentIds {
		feedComment, err := GetFeedComment(commentIds[i])
		if err == nil {
			feedComments[i] = *feedComment
		}
	}
	return feedComments
}

func GetFeedFlowAtrium(start, pageNum int) (POIFeeds, error) {
	o := orm.NewOrm()
	var feedIds []string
	qb, _ := orm.NewQueryBuilder(DB_TYPE)
	qb.Select("feed_id").From("feed").OrderBy("create_time").Desc().Limit(pageNum).Offset(start)
	sql := qb.String()
	_, err := o.Raw(sql).QueryRows(&feedIds)
	feeds := make(POIFeeds, len(feedIds))
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

func GetFeedFlowUserFeed(userId int64, start, pageNum int) POIFeeds {
	o := orm.NewOrm()
	var feedIds []string
	qb, _ := orm.NewQueryBuilder(DB_TYPE)
	qb.Select("feed_id").From("feed").Where("creator = ?").OrderBy("create_time").Asc().Limit(pageNum).Offset(start)
	sql := qb.String()
	_, err := o.Raw(sql, userId).QueryRows(&feedIds)
	if err != nil {
		seelog.Error("userId:", userId, " ", err.Error())
	}
	feeds := make(POIFeeds, len(feedIds))
	for i := range feedIds {
		feedId := feedIds[i]
		feed, err := GetFeed(feedId)
		if err == nil {
			feeds[i] = *(feed)
		}
	}
	return feeds
}

func GetFeedFlowUserFeedLike(userId int64, start, pageNum int) POIFeeds {
	var feedIds []string
	o := orm.NewOrm()
	qb, _ := orm.NewQueryBuilder(DB_TYPE)
	qb.Select("feed_id").From("feed_like").Where("user_id = ?").OrderBy("create_time").Asc().Limit(pageNum).Offset(start)
	sql := qb.String()
	_, err := o.Raw(sql, userId).QueryRows(&feedIds)
	if err != nil {
		seelog.Error("userId:", userId, " ", err.Error())
	}
	feeds := make(POIFeeds, len(feedIds))
	for i := range feedIds {
		feedId := feedIds[i]
		feed, err := GetFeed(feedId)
		if err == nil {
			feeds[i] = *(feed)
		}
	}
	return feeds
}
