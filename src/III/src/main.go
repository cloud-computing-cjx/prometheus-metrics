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
			ok_ports := []string{}
			err_ports := []string{}
			var ok_ports_str string
			var err_ports_str string
			var err_msg string
			log.Printf("ports: %s", ports)
			for _, port := range ports {
				address := net.JoinHostPort(ip, port)
				// 3 second timeout
				conn, err := net.DialTimeout("tcp", address, 3*time.Second)
				log.Printf("conn: %s, err: %s, port: %s", conn, err, port)
				if err != nil {
					results[port] = "failed"
					err_msg = err.Error()
					err_ports = append(err_ports, port)
					// todo log handler
				} else {
					if conn != nil {
						results[port] = "success"
						_ = conn.Close()
						ok_ports = append(ok_ports, port)
					} else {
						results[port] = "failed"
						err_ports = append(err_ports, port)
					}
				}
				// log.Printf("opsQueued: %s", opsQueued.Colle)
			}
			if len(ok_ports) == 0 {
				ok_ports_str = "nil"
			} else {
				ok_ports_str = strings.Join(ok_ports, ", ")
			}
			if len(err_ports) == 0 {
				err_ports_str = "nil"
			} else {
				err_ports_str = strings.Join(err_ports, ", ")
			}
			healthPortsCheck.With(prometheus.Labels{"app_name": os.Getenv("APP_NAME"), "ok": ok_ports_str, "err": err_ports_str, "err_msg": err_msg}).Set(1)
			time.Sleep(3 * time.Second)
		}
	}()
}

var (
	healthPortsCheck = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "klb",
			Subsystem: "ai_service",
			Name:      "sidecar",
			Help:      "检测所有服务端口是否可用，正常端口写入到 ok, 异常端口写入到 err, 错误信息写入到 err_msg.",
		},
		[]string{
			"app_name",
			"ok",
			"err",
			"err_msg",
		},
	)
)

func main() {
	prometheus.MustRegister(healthPortsCheck)
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
