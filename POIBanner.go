package main

import (
	"fmt"
)

type POIBanner struct {
	Id      int64  `json:"-"`
	MediaId string `json:"mediaId"`
	URL     string `json:"url"`
	Order   int64  `json:"order"`
}
type POIBanners []POIBanner

func (dbm *POIDBManager) QueryBannerList() POIBanners {
	stmtQuery, err := dbm.dbClient.Prepare(
		`SELECT id, media_id, url, rank FROM banners WHERE active = 1 ORDER BY rank ASC`)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	defer stmtQuery.Close()

	var id int64
	var mediaId string
	var url string
	var order int64
	banners := make(POIBanners, 0)

	rows, err := stmtQuery.Query()
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	for rows.Next() {
		_ = rows.Scan(&id, &mediaId, &url, &order)
		banners = append(banners, POIBanner{Id: id, MediaId: mediaId, URL: url, Order: order})
	}

	return banners
}
