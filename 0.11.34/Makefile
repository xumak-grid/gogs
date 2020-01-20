.PHONY: build-docker-aws

VERSION = 0.11.34

build-docker-aws:
	docker build --no-cache -t 281327226678.dkr.ecr.us-east-1.amazonaws.com/grid/gogs:$(VERSION) .
	docker push 281327226678.dkr.ecr.us-east-1.amazonaws.com/grid/gogs:$(VERSION)

clean:
	docker rmi 281327226678.dkr.ecr.us-east-1.amazonaws.com/grid/gogs:$(VERSION) || true
