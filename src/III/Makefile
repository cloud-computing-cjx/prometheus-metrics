export tag=v1.0

# build:
# 	echo "building prometheus-metrics binary"
# 	mkdir -p bin/amd64
# 	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/amd64 ../src/main.go

release:
	echo "building prometheus-metrics container"
	docker build -t jinyumantang/prometheus-metrics:${tag} .

push: release
	echo "pushing jinyumantang/prometheus-metrics"
	docker push jinyumantang/prometheus-metrics:${tag}