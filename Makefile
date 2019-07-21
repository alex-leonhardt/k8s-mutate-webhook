NAME = mutateme


app: deps
	go build -v -o $(NAME) cmd/main.go

deps:
	go get -v ./...

test: deps
	go test -v ./... -cover
	