package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// 检查多个端口
func tcpGather(ip string, ports []string) {
	go func() {
		for {
			results := make(map[string]string)
			log.Printf("ports: %s", ports)
			for _, port := range ports {
				address := net.JoinHostPort(ip, port)
				// 3 second timeout
				conn, err := net.DialTimeout("tcp", address, 3*time.Second)
				log.Printf("conn: %s, err: %s, port: %s", conn, err, port)
				if err != nil {
					results[port] = "failed"
					// err_string, err := json.Marshal(err)
					if err != nil {
						log.Printf("err: %s", err)
					}
					opsQueued.With(prometheus.Labels{"ok": "nil", "err": port}).Set(1)
					// todo log handler
				} else {
					if conn != nil {
						results[port] = "success"
						_ = conn.Close()
						opsQueued.With(prometheus.Labels{"ok": port, "err": "nil"}).Set(0)
					} else {
						results[port] = "failed"
						opsQueued.With(prometheus.Labels{"ok": "nil", "err": port}).Set(1)
					}
				}
				// log.Printf("opsQueued: %s", opsQueued.Colle)
			}
			time.Sleep(3 * time.Second)
		}
	}()
}

var (
	opsQueued = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "klb",
			Subsystem: "ai_service",
			Name:      os.Getenv("APP_NAME"),
			Help:      "检测所有服务端口是否可用，正常值为 0, 异常值为: 1.",
		},
		[]string{
			"err",
			"ok",
		},
	)
)

func main() {
	prometheus.MustRegister(opsQueued)
	ip := "localhost"
	var ports []string
	service_ports, ok := os.LookupEnv("SERVICE_PORTS")
	if ok {
		ports = strings.Split(service_ports, ", ")
	}
	tcpGather(ip, ports)

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":65500", nil)
}
