// POILocationController
package controllers

import "POIWolaiWebService/models"

func InsertLocation(userId int64, objectId, address, ip, userAgent string) (*models.POILocation, error) {
	location := models.POILocation{
		UserId:    userId,
		ObjectId:  objectId,
		Address:   address,
		Ip:        ip,
		UserAgent: userAgent,
	}
	l, err := models.InsertLocation(&location)
	return l, err
}
