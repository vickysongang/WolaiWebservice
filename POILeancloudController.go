package main

import ()

func SendCommentNotification(feedCommentId string) {
	feedComment := RedisManager.GetFeedComment(feedCommentId)
	feed := RedisManager.GetFeed(feedComment.FeedId)
	if feedComment == nil || feed == nil {
		return
	}

	lcTMsg := NewLCCommentNotification(feedCommentId)
	if lcTMsg == nil {
		return
	}

	// if someone comments himself...
	if feedComment.Creator.UserId != feed.Creator.UserId {
		LCSendTypedMessage(1000, feed.Creator.UserId, lcTMsg)
	}

	if feedComment.ReplyTo != nil {
		// if someone replies the author... the poor man should not be notified twice
		if feedComment.ReplyTo.UserId != feed.Creator.UserId {
			LCSendTypedMessage(1000, feedComment.ReplyTo.UserId, lcTMsg)
		}
	}

	return
}

func SendLikeNotification(userId int64, timestamp float64, feedId string) {
	user := QueryUserById(userId)
	feed := RedisManager.GetFeed(feedId)
	if user == nil || feed == nil {
		return
	}

	if user.UserId == feed.Creator.UserId {
		return
	}

	lcTMsg := NewLCLikeNotification(userId, timestamp, feedId)
	if lcTMsg == nil {
		return
	}

	LCSendTypedMessage(1000, feed.Creator.UserId, lcTMsg)

	return
}

func SendSessionNotification(sessionId int64, oprCode int64) {
	session := DbManager.QuerySessionById(sessionId)
	if session == nil {
		return
	}

	lcTMsg := NewSessionNotification(sessionId, oprCode)
	if lcTMsg == nil {
		return
	}

	switch oprCode {
	case -1:
		LCSendTypedMessage(session.Teacher.UserId, session.Creator.UserId, lcTMsg)
	case 1:
		LCSendTypedMessage(session.Creator.UserId, session.Teacher.UserId, lcTMsg)
	case 2:
		LCSendTypedMessage(session.Teacher.UserId, session.Creator.UserId, lcTMsg)
	case 3:
		LCSendTypedMessage(session.Creator.UserId, session.Teacher.UserId, lcTMsg)
		LCSendTypedMessage(session.Teacher.UserId, session.Creator.UserId, lcTMsg)
	}
}
