package main

import (
	_ "encoding/json"
)

const (
	FEEDTYPE_MICROBLOG = iota
	FEEDTYPE_SHARE     = iota
	FEEDTYPE_REPOST    = iota
)

type POIFeed struct {
	Id              string            `json:"id"`
	Creator         *POIUser          `json:"creatorInfo"`
	CreateTimestamp float64           `json:"createTimestamp"`
	FeedType        int               `json:"feedType"`
	Text            string            `json:"text"`
	ImageList       []string          `json:"imageList,omitempty"`
	OriginFeed      *POIFeed          `json:"originFeedInfo,omitempty"`
	Attribute       map[string]string `json:"attribute,omitempty"`
	LikeCount       int               `json:"likeCount"`
	CommentCount    int               `json:"commentCount"`
	RepostCount     int               `json:"repostCount"`
}

type POIFeeds []POIFeed

func NewPOIFeed() POIFeed {
	return POIFeed{ImageList: make([]string, 9), Attribute: make(map[string]string)}
}
