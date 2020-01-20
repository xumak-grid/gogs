.PHONY: build build-exposer build-initial build-docker build-docker-aws

build: build-exposer build-initial

build-exposer:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build  -o ./bin/exposer exposer/main.go
	chmod +x ./bin/exposer
build-initial:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build  -o ./bin/initial initial/*.go
	chmod +x ./bin/initial

VERSION = 0.11.19
build-docker:
	docker build --no-cache -t :5000/grid/gogs:latest .
	docker tag :5000/grid/gogs:latest :5000/grid/gogs:$(VERSION)
	docker push :5000/grid/gogs:latest
	docker push :5000/grid/gogs:$(VERSION)

build-docker-aws:
	docker build --no-cache -t /grid/gogs:latest .
	docker tag /grid/gogs:latest /grid/gogs:$(VERSION)
	docker push /grid/gogs:latest
	docker push /grid/gogs:$(VERSION)

clean:
	docker rmi :5000/grid/gogs:latest || true
	docker rmi :5000/grid/gogs:$(VERSION) || true
	docker rmi /grid/gogs:latest || true
	docker rmi /grid/gogs:$(VERSION) || true
