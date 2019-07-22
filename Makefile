NAME = mutateme
IMAGE_PREFIX = alexleonhardt
IMAGE_NAME = $$(basename `pwd`)
IMAGE_VERSION = $$(git log --abbrev-commit --format=%h -s | head -n 1)

app: deps
	go build -v -o $(NAME) cmd/main.go

deps:
	go get -v ./...

test: deps
	go test -v ./... -cover
	
docker:
	docker build -t $(IMAGE_PREFIX)/$(IMAGE_NAME):$(IMAGE_VERSION) .
	docker tag $(IMAGE_PREFIX)/$(IMAGE_NAME):$(IMAGE_VERSION) $(IMAGE_PREFIX)/$(IMAGE_NAME):latest

push:
	docker push $(IMAGE_PREFIX)/$(IMAGE_NAME):$(IMAGE_VERSION) 
	docker push $(IMAGE_PREFIX)/$(IMAGE_NAME):latest

