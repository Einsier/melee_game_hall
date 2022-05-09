package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	//定义游戏指标
	GaugeGameRoomCount = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "game_room_count",
		Help: "当前游戏房间数量",
	})
	// GaugeGameRoomCount.Inc() 游戏房间数量 +1
	// GaugeGameRoomCount.Dec() 游戏房间数量 -1
	// GaugeGameRoomCount.Add(n) 游戏房间数量 +n
	// GaugeGameRoomCount.Sub(n) 游戏房间数量 -n
	// GaugeGameRoomCount.Set(n) 游戏房间数量设置为 n
)

func Start() {
	// 将指标注册到 Prometheus 默认仓库
	prometheus.MustRegister(GaugeGameRoomCount)

	// Serve the default Prometheus metrics registry over HTTP on /metrics.
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":8080", nil)
}
