export tag=v1.0

# build:
# 	echo "building prometheus-metrics binary"
# 	mkdir -p bin/amd64
# 	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/amd64 ../src/main.go

release: build
	echo "building prometheus-metrics container"
	docker build -t 10.0.1.191/prometheus-metrics:${tag} .

push: release
	echo "pushing 10.0.1.191/prometheus-metrics"
	docker push 10.0.1.191/prometheus-metrics:${tag}