package main

import (
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func TestTcpGather(t *testing.T) {
	os.Setenv("SERVICE_PORTS", "9090, 2112")
	os.Setenv("APP_NAME", "TestTcpGather")
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
