package main

import (
	"net/http"
	_ "net/http/pprof"
	"strconv"
	"testing"
	"time"

	"github.com/cihub/seelog"

	orderController "WolaiWebservice/controllers/order"
	"WolaiWebservice/websocket"
)

func BenchmarkOrder(b *testing.B) {
	go func() {
		seelog.Critical(http.ListenAndServe(":6060", nil))
	}()

	for i := 0; i < b.N; i++ {
		_, _, order := orderController.CreateOrder(10010, 0, 3, 1, 1, "Y")

		msg := websocket.NewPOIWSMessage("", 10010, websocket.WS_ORDER2_CREATE)
		msg.Attribute["orderId"] = strconv.FormatInt(order.Id, 10)

		go websocket.InitOrderDispatch(msg, time.Now().Unix())
	}
}
