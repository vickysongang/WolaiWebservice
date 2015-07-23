package main

import (
	_ "encoding/json"
)

type POIFeedComment struct {
	Id              string   `json:"id"`
	FeedId          string   `json:"feedId"`
	Creator         *POIUser `json:"creatorInfo"`
	CreateTimestamp float64  `json:"createTimestamp"`
	Text            string   `json:"text"`
	ImageList       []string `json:"imageList,omitempty"`
	ReplyTo         *POIUser `json:"replyTo,omitempty"`
	LikeCount       int      `json:"likeCount"`
}

type POIFeedComments []POIFeedComment

func NewPOIFeedComment() POIFeedComment {
	return POIFeedComment{ImageList: make([]string, 9)}
}
