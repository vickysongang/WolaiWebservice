// POIFeedModel.go
package models

import (
	"encoding/json"
	"time"

	"github.com/astaxie/beego/orm"
	"github.com/cihub/seelog"

	"WolaiWebservice/config"
)

const (
	FEEDTYPE_MICROBLOG = iota
	FEEDTYPE_SHARE     = iota
	FEEDTYPE_REPOST    = iota
)

type POIFeed struct {
	Id              string            `json:"id" orm:"pk;column(feed_id)"`
	Creator         *User             `json:"creatorInfo" orm:"-"`
	CreateTimestamp float64           `json:"createTimestamp" orm:"-"`
	FeedType        int64             `json:"feedType"`
	Text            string            `json:"text"`
	ImageList       []string          `json:"imageList,omitempty" orm:"-"`
	OriginFeed      *POIFeed          `json:"originFeedInfo,omitempty" orm:"-"`
	Attribute       map[string]string `json:"attribute,omitempty" orm:"-"`
	LikeCount       int64             `json:"likeCount" orm:"-"`
	CommentCount    int64             `json:"commentCount" orm:"-"`
	RepostCount     int64             `json:"-" orm:"-"`
	HasLiked        bool              `json:"hasLiked" orm:"-"`
	HasFaved        bool              `json:"-" orm:"-"`
	Created         int64             `json:"-" orm:"column(creator)"`
	CreateTime      time.Time         `json:"-" orm:"auto_now_add;type(datetime)"`
	ImageInfo       string            `json:"-"`
	AttributeInfo   string            `json:"-"`
	OriginFeedId    string            `json:"-"`
	PlateType       string            `json:"plateType"`
	TopFlag         bool              `json:"topFlag" orm:"-"`
	DeleteFlag      string            `json:"-"`
	TopSeq          string            `json:"-"`
}

type POIFeedLike struct {
	Id         int64     `json:"-" orm:"pk"`
	FeedId     string    `json:"feedId"`
	UserId     int64     `json:"userId"`
	CreateTime time.Time `json:"-" orm:"auto_now_add;type(datetime)"`
}

type POIFeedDetail struct {
	Feed       *POIFeed        `json:"feedInfo"`
	LikedUsers []User          `json:"likedUsers"`
	Comments   POIFeedComments `json:"comments"`
}

type POIFeeds []POIFeed

func (f *POIFeed) TableName() string {
	return "feed"
}

func (fl *POIFeedLike) TableName() string {
	return "feed_like"
}

func init() {
	orm.RegisterModel(new(POIFeed), new(POIFeedLike))
}

func NewPOIFeed() POIFeed {
	return POIFeed{ImageList: make([]string, 9), Attribute: make(map[string]string)}
}

func (f *POIFeed) IncreaseLike() {
	f.LikeCount = f.LikeCount + 1
}

func (f *POIFeed) DecreaseLike() {
	f.LikeCount = f.LikeCount - 1
}

func (f *POIFeed) IncreaseComment() {
	f.CommentCount = f.CommentCount + 1
}

func (f *POIFeed) IncreaseRepost() {
	f.RepostCount = f.RepostCount + 1
}

func InsertPOIFeed(feed *POIFeed) (*POIFeed, error) {
	o := orm.NewOrm()
	_, err := o.Insert(feed)
	if err != nil {
		seelog.Error("feed:", feed, " ", err.Error())
		return nil, err
	}
	return feed, nil
}

func UpdateFeedInfo(feedId string, feedInfo map[string]interface{}) {
	o := orm.NewOrm()
	var params orm.Params = make(orm.Params)
	for k, v := range feedInfo {
		params[k] = v
	}
	_, err := o.QueryTable("feed").Filter("feed_id", feedId).Update(params)
	if err != nil {
		seelog.Error("feedId:", feedId, " feedInfo:", feedInfo, " ", err.Error())
	}
}

func UpdateFeedTopSeq() {
	o := orm.NewOrm()
	var params orm.Params = make(orm.Params)
	params["top_seq"] = ""
	o.QueryTable("feed").Filter("top_seq__isnull", false).Update(params)
}

func GetFeed(feedId string) (*POIFeed, error) {
	feed := POIFeed{}
	qb, _ := orm.NewQueryBuilder(config.Env.Database.Type)
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
	creator, _ := ReadUser(feed.Created)
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

func InsertPOIFeedLike(feedLike *POIFeedLike) *POIFeedLike {
	o := orm.NewOrm()
	_, err := o.Insert(feedLike)
	if err != nil {
		seelog.Error("feedLike:", feedLike, " ", err.Error())
		return nil
	}
	return feedLike
}

//物理删除Feed
func DeletePOIFeed(feedId string) {
	o := orm.NewOrm()
	o.QueryTable("feed").Filter("feed_id", feedId).Delete()
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

func GetFeedLikeList(feedId string) []User {
	o := orm.NewOrm()
	qb, _ := orm.NewQueryBuilder(config.Env.Database.Type)
	qb.Select("user_id").From("feed_like").Where("feed_id=?")
	sql := qb.String()
	var userIds []int64
	_, err := o.Raw(sql, feedId).QueryRows(&userIds)
	if err != nil {
		seelog.Error("feedId:", feedId, " ", err.Error())
	}
	users := make([]User, len(userIds))
	for i := range userIds {
		user, _ := ReadUser(userIds[i])
		if user == nil {
			continue
		}
		users[i] = *(user)
	}
	return users
}

func HasLikedFeed(feed *POIFeed, user *User) bool {
	o := orm.NewOrm()
	count, err := o.QueryTable("feed_like").Filter("feed_id", feed.Id).Filter("user_id", user.Id).Count()
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
	feedComment.Creator, _ = ReadUser(feedComment.Created)

	var imageList []string
	err = json.Unmarshal([]byte(feedComment.ImageInfo), &imageList)
	if err == nil {
		feedComment.ImageList = imageList
	} else {
		seelog.Error("unmarshal ImageInfo:", feedComment.ImageInfo, " ", err.Error())
		return nil, err
	}
	if feedComment.ReplyToId != 0 {
		feedComment.ReplyTo, _ = ReadUser(feedComment.ReplyToId)
	}
	return &feedComment, nil
}

func GetFeedComments(feedId string) POIFeedComments {
	o := orm.NewOrm()
	qb, _ := orm.NewQueryBuilder(config.Env.Database.Type)
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

func GetFeedFlowAtrium(start, pageCount int) (POIFeeds, error) {
	o := orm.NewOrm()
	var feedIds []string
	qb, _ := orm.NewQueryBuilder(config.Env.Database.Type)
	qb.Select("feed.feed_id").From("feed").InnerJoin("users").On("feed.creator = users.id").
		OrderBy("feed.create_time").Desc().Limit(pageCount).Offset(start)
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

//获取置顶动态
func GetTopFeedFlowAtrium(plateType string) (POIFeeds, error) {
	o := orm.NewOrm()
	var feedIds []string
	qb, _ := orm.NewQueryBuilder(config.Env.Database.Type)
	qb.Select("feed.feed_id").From("feed").InnerJoin("users").On("feed.creator = users.id").
		Where("feed.plate_type like ? and feed.top_seq is not null and feed.top_seq <> ''").OrderBy("feed.top_seq").Desc()
	sql := qb.String()
	_, err := o.Raw(sql, "%"+plateType+"%").QueryRows(&feedIds)
	feeds := make(POIFeeds, len(feedIds))
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
func GetFeedFlowAtriumByPlateType(start, pageCount int, plateType string) (POIFeeds, error) {
	o := orm.NewOrm()
	var feedIds []string
	qb, _ := orm.NewQueryBuilder(config.Env.Database.Type)
	qb.Select("feed.feed_id").From("feed").InnerJoin("users").On("feed.creator = users.id").Where("feed.plate_type like ? and (feed.top_seq is null or feed.top_seq='')").
		OrderBy("feed.create_time").Desc().Limit(pageCount).Offset(start)
	sql := qb.String()
	_, err := o.Raw(sql, "%"+plateType+"%").QueryRows(&feedIds)
	feeds := make(POIFeeds, len(feedIds))
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

func GetFeedFlowUserFeed(userId int64, start, pageCount int) POIFeeds {
	o := orm.NewOrm()
	var feedIds []string
	qb, _ := orm.NewQueryBuilder(config.Env.Database.Type)
	qb.Select("feed_id").From("feed").Where("creator = ?").OrderBy("create_time").Asc().Limit(pageCount).Offset(start)
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

func GetFeedFlowUserFeedLike(userId int64, start, pageCount int) POIFeeds {
	var feedIds []string
	o := orm.NewOrm()
	qb, _ := orm.NewQueryBuilder(config.Env.Database.Type)
	qb.Select("feed_id").From("feed_like").Where("user_id = ?").OrderBy("create_time").Asc().Limit(pageCount).Offset(start)
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
