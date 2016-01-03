package main

import (
	"strconv"
	"testing"
	"time"

	orderController "WolaiWebservice/controllers/order"
	"WolaiWebservice/websocket"
)

func BenchmarkOrder(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _, order := orderController.CreateOrder(10010, 0, 3, 1, 1, "Y")

		msg := websocket.NewPOIWSMessage("", 10010, websocket.WS_ORDER2_CREATE)
		msg.Attribute["orderId"] = strconv.FormatInt(order.Id, 10)

		go websocket.InitOrderDispatch(msg, time.Now().Unix())
	}
}
